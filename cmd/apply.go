package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cvhariharan/flowctl/pkg/client"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply flow definitions from YAML files to a running flowctl server",
	Long: `Reconcile flow definitions declared in YAML against a flowctl server.

Reads one or more YAML files (or a directory of them, with --recursive) and
creates or updates each flow on the server. Existing flows are looked up by
the slug derived from metadata.name; missing flows are created, existing ones
are updated. The command is idempotent.

YAML files may contain multiple documents separated by ---.`,
	Example: `  flowctl apply -f flow.yaml
  flowctl apply -f flows/ -R
  cat flow.yaml | flowctl apply -f -
  flowctl apply -f bundle.yaml --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, _ := cmd.Flags().GetStringSlice("filename")
		recursive, _ := cmd.Flags().GetBool("recursive")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		server, _ := cmd.Flags().GetString("server")
		token, _ := cmd.Flags().GetString("token")
		insecure, _ := cmd.Flags().GetBool("insecure")

		if server == "" {
			server = os.Getenv("FLOWCTL_SERVER")
		}
		if token == "" {
			token = os.Getenv("FLOWCTL_TOKEN")
		}

		if len(paths) == 0 {
			return fmt.Errorf("at least one -f/--filename is required")
		}
		if server == "" {
			return fmt.Errorf("server URL is required: pass --server or set FLOWCTL_SERVER")
		}
		if token == "" {
			return fmt.Errorf("api token is required: pass --token or set FLOWCTL_TOKEN")
		}

		return runApply(applyOptions{
			paths:     paths,
			recursive: recursive,
			dryRun:    dryRun,
			server:    server,
			token:     token,
			insecure:  insecure,
		})
	},
	SilenceUsage: true,
}

func init() {
	applyCmd.Flags().StringSliceP("filename", "f", nil, "YAML file, directory, or - for stdin (repeatable)")
	applyCmd.Flags().BoolP("recursive", "R", false, "Recurse into directories passed to -f")
	applyCmd.Flags().Bool("dry-run", false, "Resolve and print the create/update plan; make no writes")
	applyCmd.Flags().String("server", "", "Server base URL (or FLOWCTL_SERVER)")
	applyCmd.Flags().String("token", "", "Bearer API token (or FLOWCTL_TOKEN)")
	applyCmd.Flags().Bool("insecure", false, "Skip TLS certificate verification")
	rootCmd.AddCommand(applyCmd)
}

type applyOptions struct {
	paths     []string
	recursive bool
	dryRun    bool
	server    string
	token     string
	insecure  bool
}

type namedBlob struct {
	name string // source path or "<stdin>"
	data []byte
}

type yamlDoc struct {
	id   string         // slug(metadata.name) — matches the server's GenerateSlug
	ns   string         // metadata.namespace
	body map[string]any // full raw map for JSON marshaling
}

func (d yamlDoc) label() string { return "flow/" + d.id }

func runApply(opts applyOptions) error {
	blobs, err := gatherInputs(opts.paths, opts.recursive)
	if err != nil {
		return err
	}

	docs, parseErrs := decodeDocs(blobs)

	cl, err := newAPIClient(opts.server, opts.token, opts.insecure)
	if err != nil {
		return err
	}

	ctx := context.Background()
	var created, configured, errCount int
	for _, e := range parseErrs {
		fmt.Fprintf(os.Stderr, "%s: %v\n", e.source, e.err)
		errCount++
	}
	for _, d := range docs {
		action, err := reconcile(ctx, cl, d, opts.dryRun)
		if err != nil {
			fmt.Printf("%s\terror: %v\n", d.label(), err)
			errCount++
			continue
		}
		fmt.Printf("%s\t%s\n", d.label(), action)
		switch action {
		case "created", "would create":
			created++
		case "configured", "would update":
			configured++
		}
	}

	if opts.dryRun {
		fmt.Printf("\nPlan: %d to create, %d to update, %d errors\n", created, configured, errCount)
	} else {
		fmt.Printf("\n%d created, %d configured, %d errors\n", created, configured, errCount)
	}

	if errCount > 0 {
		return fmt.Errorf("apply finished with %d error(s)", errCount)
	}
	return nil
}

// gatherInputs resolves -f values into a list of named byte blobs.
// `-` reads stdin; a directory is walked (respecting recursive) keeping only
// .yaml/.yml files; a regular file is read as-is.
func gatherInputs(paths []string, recursive bool) ([]namedBlob, error) {
	var out []namedBlob
	stdinUsed := false
	for _, p := range paths {
		if p == "-" {
			if stdinUsed {
				continue
			}
			stdinUsed = true
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return nil, fmt.Errorf("read stdin: %w", err)
			}
			out = append(out, namedBlob{name: "<stdin>", data: data})
			continue
		}
		info, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", p, err)
		}
		if info.IsDir() {
			err := filepath.WalkDir(p, func(path string, d fs.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if d.IsDir() {
					if !recursive && path != p {
						return fs.SkipDir
					}
					return nil
				}
				if !isYAMLFile(path) {
					return nil
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("read %s: %w", path, err)
				}
				out = append(out, namedBlob{name: path, data: data})
				return nil
			})
			if err != nil {
				return nil, err
			}
			continue
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", p, err)
		}
		out = append(out, namedBlob{name: p, data: data})
	}
	return out, nil
}

func isYAMLFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}

type docErr struct {
	source string
	err    error
}

// decodeDocs splits each blob into YAML documents, probes the kind, and
// extracts the natural-key fields. Per-doc errors are reported but do not
// abort the batch.
func decodeDocs(blobs []namedBlob) ([]yamlDoc, []docErr) {
	var docs []yamlDoc
	var errs []docErr
	for _, b := range blobs {
		dec := yaml.NewDecoder(bytes.NewReader(b.data))
		idx := 0
		for {
			var node yaml.Node
			if err := dec.Decode(&node); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				errs = append(errs, docErr{source: docSource(b.name, idx), err: fmt.Errorf("yaml decode: %w", err)})
				break
			}
			idx++
			if node.Kind == 0 || isEmptyNode(&node) {
				continue
			}
			body := map[string]any{}
			if err := node.Decode(&body); err != nil {
				errs = append(errs, docErr{source: docSource(b.name, idx-1), err: fmt.Errorf("yaml decode: %w", err)})
				continue
			}
			d, err := buildDoc(body)
			if err != nil {
				errs = append(errs, docErr{source: docSource(b.name, idx-1), err: err})
				continue
			}
			docs = append(docs, d)
		}
	}
	return docs, errs
}

func docSource(blob string, idx int) string {
	if idx == 0 {
		return blob
	}
	return fmt.Sprintf("%s#%d", blob, idx)
}

func isEmptyNode(n *yaml.Node) bool {
	if n == nil {
		return true
	}
	if n.Kind == yaml.DocumentNode && len(n.Content) == 0 {
		return true
	}
	if n.Kind == yaml.ScalarNode && n.Value == "" && n.Tag == "!!null" {
		return true
	}
	return false
}

func buildDoc(body map[string]any) (yamlDoc, error) {
	if kind := stringField(body, "kind"); kind != "" && !strings.EqualFold(kind, "Flow") {
		return yamlDoc{}, fmt.Errorf("unsupported kind %q (this version of apply only supports Flow)", kind)
	}
	delete(body, "kind") // server doesn't expect kind in the request body

	meta, _ := body["metadata"].(map[string]any)
	if meta == nil {
		return yamlDoc{}, fmt.Errorf("metadata is required")
	}
	name := stringField(meta, "name")
	if name == "" {
		return yamlDoc{}, fmt.Errorf("metadata.name is required")
	}
	ns := stringField(meta, "namespace")
	if ns == "" {
		return yamlDoc{}, fmt.Errorf("metadata.namespace is required")
	}

	// id is always derived from name to match the server's HandleCreateFlow,
	// which calls GenerateSlug(name) and ignores any client-supplied id.
	return yamlDoc{id: slugify(name), ns: ns, body: body}, nil
}

func stringField(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return strings.TrimSpace(v)
}

// slugify mirrors handlers.GenerateSlug: slug.Make(name) with "-" → "_".
func slugify(s string) string {
	return strings.ReplaceAll(slug.Make(s), "-", "_")
}

func reconcile(ctx context.Context, cl *client.Client, d yamlDoc, dryRun bool) (string, error) {
	found, err := lookupFlow(ctx, cl, d.ns, d.id)
	if err != nil {
		return "", err
	}
	if dryRun {
		if found {
			return "would update", nil
		}
		return "would create", nil
	}
	if found {
		var body client.FlowUpdateReq
		if err := mapToTyped(buildUpdateBody(d.body), &body); err != nil {
			return "", fmt.Errorf("convert update body: %w", err)
		}
		resp, err := cl.UpdateFlow(ctx, d.ns, d.id, body)
		if err != nil {
			return "", err
		}
		if err := checkStatus(resp); err != nil {
			return "", err
		}
		return "configured", nil
	}
	var body client.FlowCreateReq
	if err := mapToTyped(d.body, &body); err != nil {
		return "", fmt.Errorf("convert create body: %w", err)
	}
	resp, err := cl.CreateFlow(ctx, d.ns, body)
	if err != nil {
		return "", err
	}
	if err := checkStatus(resp); err != nil {
		return "", err
	}
	return "created", nil
}

// lookupFlow returns true if a flow with the given slug exists in namespace.
// A 404 is mapped to (false, nil); other non-2xx responses are returned as errors.
func lookupFlow(ctx context.Context, cl *client.Client, ns, id string) (bool, error) {
	resp, err := cl.GetFlow(ctx, ns, id)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}
	body, _ := io.ReadAll(resp.Body)
	return false, fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
}

// checkStatus reads and discards the response body, returning an error if the
// status code is not 2xx.
func checkStatus(resp *http.Response) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
}

// mapToTyped JSON-marshals a generic map and unmarshals it into a typed
// struct (one of client.FlowCreateReq / client.FlowUpdateReq). This avoids
// duplicating the schema while still letting the generated client do the
// wire encoding.
func mapToTyped(m map[string]any, out any) error {
	buf, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, out)
}

// buildUpdateBody flattens a YAML "metadata" block into the shape expected by
// FlowUpdateReq (top-level prefix/description/allow_overlap/...).
// metadata.{id,name,namespace} live in the URL and are dropped.
func buildUpdateBody(create map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range create {
		if k == "metadata" {
			continue
		}
		out[k] = v
	}
	if meta, ok := create["metadata"].(map[string]any); ok {
		for _, k := range []string{"prefix", "description", "allow_overlap", "user_schedulable", "max_retries", "schedules"} {
			if v, ok := meta[k]; ok {
				out[k] = v
			}
		}
	}
	if _, ok := out["inputs"]; !ok {
		out["inputs"] = []any{}
	}
	if _, ok := out["actions"]; !ok {
		out["actions"] = []any{}
	}
	return out
}

// newAPIClient builds a generated *client.Client wired with bearer-token
// auth, a 30s timeout, and (optionally) TLS-verification skipped.
func newAPIClient(server, token string, insecure bool) (*client.Client, error) {
	base := strings.TrimRight(server, "/") + "/"
	tr := &http.Transport{}
	if insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	httpClient := &http.Client{Timeout: 30 * time.Second, Transport: tr}
	addAuth := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		return nil
	}
	return client.NewClient(base,
		client.WithHTTPClient(httpClient),
		client.WithRequestEditorFn(addAuth),
	)
}


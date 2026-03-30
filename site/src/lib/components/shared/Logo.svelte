<script lang="ts">
  import logoSvg from '$lib/assets/logo.svg';
  import fullLogoSvg from '$lib/assets/full-logo.svg';
  import fullLogoLightSvg from '$lib/assets/full-logo-light.svg';
  import { resolvedTheme } from '$lib/stores/theme';

  let {
    height = '7rem',
    iconOnly = false
  }: {
    height?: string;
    iconOnly?: boolean;
  } = $props();

  let logoSrc = $derived(
    iconOnly
      ? logoSvg
      : $resolvedTheme === 'dark'
        ? fullLogoLightSvg
        : fullLogoSvg
  );
</script>

<!-- Logo -->
<img
  src={logoSrc}
  alt="Flowctl"
  style="height: {height}; width: auto; image-rendering: -webkit-optimize-contrast; image-rendering: crisp-edges;"
/>

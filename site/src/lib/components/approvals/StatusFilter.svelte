<script lang="ts">
  let {
    value = $bindable(''),
    onChange
  }: {
    value: string;
    onChange: (status: string) => void;
  } = $props();

  function handleChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    value = target.value;
    onChange(target.value);
  }
</script>

<select
  {value}
  onchange={handleChange}
>
  <option value="">All Statuses</option>
  <option value="pending">Pending</option>
  <option value="approved">Approved</option>
  <option value="rejected">Rejected</option>
</select>

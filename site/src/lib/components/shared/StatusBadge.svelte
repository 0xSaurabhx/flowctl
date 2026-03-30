<script lang="ts">
  let { value }: { value: string } = $props();

  const capitalizeStatus = (status: string) => {
    if (!status) return '';
    return status.charAt(0).toUpperCase() + status.slice(1);
  };

  const getBadgeClass = (status: string) => {
    switch (status) {
      case 'completed':
      case 'approved':
        return 'success';
      case 'pending':
      case 'cancelled':
        return 'warning';
      case 'running':
        return '';
      case 'errored':
      case 'rejected':
        return 'danger';
      default:
        return '';
    }
  };
</script>

<span class="badge {getBadgeClass(value)}">
  {capitalizeStatus(value)}
</span>

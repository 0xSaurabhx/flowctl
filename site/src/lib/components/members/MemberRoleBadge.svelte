<script lang="ts">
  let { value }: { value: string } = $props();

  function getRoleClass(role: string) {
    switch (role) {
      case 'user': return 'success';
      case 'operator': return 'info';
      case 'reviewer': return 'warning';
      case 'admin': return 'danger';
      default: return '';
    }
  }

  function capitalizeRole(role: string) {
    return role.charAt(0).toUpperCase() + role.slice(1);
  }
</script>

<span class="badge {getRoleClass(value)}">
  {capitalizeRole(value)}
</span>

<style>
  .badge {
    display: inline-flex;
    align-items: center;
    padding: 0.125rem 0.625rem;
    border-radius: 9999px;
    font-size: 0.75rem;
    font-weight: 500;
    background: var(--faint);
    color: var(--foreground);
  }
  .badge.success {
    background: color-mix(in srgb, var(--success) 15%, transparent);
    color: var(--success);
  }
  .badge.info {
    background: color-mix(in srgb, var(--primary) 15%, transparent);
    color: var(--primary);
  }
  .badge.warning {
    background: color-mix(in srgb, var(--warning) 15%, transparent);
    color: var(--warning);
  }
  .badge.danger {
    background: color-mix(in srgb, var(--danger) 15%, transparent);
    color: var(--danger);
  }
</style>

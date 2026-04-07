<script lang="ts">
  import { notifications, type Notification } from '$lib/stores/notifications';
  import { fly, scale } from 'svelte/transition';
  import { flip } from 'svelte/animate';
  import {
    IconCircleCheck,
    IconAlertCircle,
    IconAlertTriangle,
    IconInfoCircle,
    IconX
  } from '@tabler/icons-svelte';

  const handleDismiss = (id: string) => {
    notifications.remove(id);
  };

  const getIcon = (type: Notification['type']) => {
    switch (type) {
      case 'success': return IconCircleCheck;
      case 'error': return IconAlertCircle;
      case 'warning': return IconAlertTriangle;
      case 'info':
      default: return IconInfoCircle;
    }
  };

  const getToastClass = (type: Notification['type']) => {
    switch (type) {
      case 'success': return 'success';
      case 'error': return 'danger';
      case 'warning': return 'warning';
      case 'info':
      default: return '';
    }
  };
</script>

<div class="toast-container" data-placement="top-right" role="region" aria-label="Notifications">
  {#each $notifications as notification (notification.id)}
    {@const IconComponent = getIcon(notification.type)}
    <div
      class="toast {getToastClass(notification.type)}"
      role="alert"
      aria-live={notification.type === 'error' ? 'assertive' : 'polite'}
      in:fly={{ x: 300, duration: 300 }}
      out:scale={{ duration: 200 }}
      animate:flip={{ duration: 200 }}
    >
      <div class="hstack gap-3">
        <IconComponent size={18} aria-hidden="true" />
        <div class="flex-1">
          <h3 class="text-sm">{notification.title}</h3>
          <p class="text-sm text-light">{notification.message}</p>
        </div>
        {#if notification.dismissible}
          <button
            onclick={() => handleDismiss(notification.id)}
            aria-label="Dismiss notification"
            class="ghost small"
          >
            <IconX size={16} aria-hidden="true" />
          </button>
        {/if}
      </div>
    </div>
  {/each}
</div>

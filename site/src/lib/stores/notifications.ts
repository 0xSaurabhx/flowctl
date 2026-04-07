export interface Notification {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  title: string;
  message: string;
  duration?: number;
  dismissible?: boolean;
}

type NotificationType = Notification['type'];

const variantMap: Record<NotificationType, string | undefined> = {
  success: 'success',
  error: 'danger',
  warning: 'warning',
  info: undefined
};

function toastNotification(type: NotificationType, title: string, message: string, options?: Partial<Notification>) {
  const duration = options?.duration ?? (type === 'error' ? 0 : 5000);

  if (typeof window !== 'undefined' && typeof (window as any).ot?.toast === 'function') {
    (window as any).ot.toast(message, title, {
      variant: variantMap[type],
      duration
    });
  }
}

function createNotificationStore() {
  return {
    // Keep the same public API so callers don't need to change
    add: (notification: Omit<Notification, 'id'>) => {
      toastNotification(notification.type, notification.title, notification.message, notification);
      return `notification-${Date.now()}`;
    },
    remove: (_id: string) => { /* ot.toast handles auto-dismiss */ },
    clear: () => {
      if (typeof window !== 'undefined' && typeof (window as any).ot?.toast?.clear === 'function') {
        (window as any).ot.toast.clear();
      }
    },
    success: (title: string, message: string, options?: Partial<Notification>) => {
      toastNotification('success', title, message, options);
    },
    error: (title: string, message: string, options?: Partial<Notification>) => {
      toastNotification('error', title, message, { duration: 0, ...options });
    },
    warning: (title: string, message: string, options?: Partial<Notification>) => {
      toastNotification('warning', title, message, options);
    },
    info: (title: string, message: string, options?: Partial<Notification>) => {
      toastNotification('info', title, message, options);
    },
    // Svelte store compatibility - provide subscribe that returns empty array
    subscribe: (fn: (value: Notification[]) => void) => {
      fn([]);
      return () => {};
    }
  };
}

export const notifications = createNotificationStore();

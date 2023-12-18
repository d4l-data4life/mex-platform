let isLeaving = false;
window.addEventListener('unload', () => (isLeaving = true));

export const catchRetryableAction = async (
  action: Function,
  canClose: boolean = true,
  closeHandler: (fromRetry: boolean) => void = null
) => {
  try {
    await action();
  } catch (error) {
    if (isLeaving) {
      return;
    }

    console.error(error);
    window.dispatchEvent(
      new CustomEvent('showError', {
        bubbles: false,
        cancelable: false,
        detail: {
          retryHandler: () => catchRetryableAction(action, canClose, closeHandler),
          closeHandler,
          error,
          canClose,
        },
      })
    );
  }
};

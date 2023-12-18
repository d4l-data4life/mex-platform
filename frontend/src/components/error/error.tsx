import { Component, h, Listen, State } from '@stencil/core';
import { ResponseError } from 'models/response-error';

@Component({
  tag: 'mex-error',
})
export class ErrorComponent {
  modalContentButtonRefs: HTMLButtonElement[] = [];

  @State() error?: ResponseError | Error;
  @State() retryHandlers: Function[] = [];
  @State() closeHandlers: Function[] = [];
  @State() canClose: boolean;
  @State() wantsToClose: boolean = false;

  @Listen('showError', {
    target: 'window',
  })
  showErrorHandler(event: CustomEvent) {
    const { error, retryHandler, closeHandler, canClose } = event.detail;

    this.retryHandlers.push(retryHandler);
    this.error = error;
    this.canClose = canClose;
    this.wantsToClose = false;
    canClose && closeHandler && this.closeHandlers.push(closeHandler);
  }

  @Listen('closeError', {
    target: 'window',
  })
  closeErrorHandler() {
    this.canClose && this.close();
  }

  get isShown() {
    return !!this.error;
  }

  get canRetry() {
    return !!this.retryHandlers.length;
  }

  get errorText() {
    const { status = null } = (this.error as ResponseError) ?? {};
    if (status === null) {
      return 'error.client';
    }

    return status === 0 ? 'error.network' : 'error.server';
  }

  close(fromRetry: boolean = false) {
    if (!fromRetry && !this.wantsToClose) {
      this.wantsToClose = true;
      this.canClose = false;
      return;
    }

    const { closeHandlers } = this;
    this.error = null;
    this.retryHandlers = [];
    this.closeHandlers = [];
    closeHandlers.map((handler) => handler(fromRetry));
  }

  setFocus(closeRef: HTMLButtonElement) {
    this.isShown && (this.modalContentButtonRefs[this.modalContentButtonRefs.length - 1] ?? closeRef)?.focus();
  }

  async retry() {
    const { retryHandlers } = this;
    this.close(true);
    await Promise.all(retryHandlers.map((handler) => handler()));
  }

  render() {
    const { canRetry, canClose, wantsToClose } = this;
    const closeButtonConfig = {
      clickHandler: () => this.close(),
      testAttr: 'error:modal:close',
    };
    const retryButtonConfig = {
      label: 'error.retry',
      clickHandler: () => this.retry(),
      testAttr: 'error:modal:retry',
    };

    return (
      this.isShown && (
        <mex-modal handleClose={canClose ? this.close.bind(this) : null} handleSetFocus={this.setFocus.bind(this)}>
          {!wantsToClose && (
            <mex-modal-contents
              context={this}
              illustration="error"
              caption="error.title"
              text={this.errorText}
              buttons={
                canRetry
                  ? [
                      ...(canClose
                        ? [
                            {
                              ...closeButtonConfig,
                              label: 'error.ignoreIssue',
                              modifier: 'secondary' as 'secondary',
                            },
                          ]
                        : []),
                      retryButtonConfig,
                    ]
                  : []
              }
              data-test="error:modal"
            />
          )}
          {wantsToClose && (
            <mex-modal-contents
              context={this}
              illustration="confirm"
              caption="error.closeConfirm.title"
              text="error.closeConfirm.text"
              buttons={
                canRetry
                  ? [
                      { ...retryButtonConfig, modifier: 'secondary' },
                      {
                        ...closeButtonConfig,
                        label: 'error.closeConfirm.continue',
                      },
                    ]
                  : []
              }
              data-test="error:closeConfirm:modal"
            />
          )}
        </mex-modal>
      )
    );
  }
}

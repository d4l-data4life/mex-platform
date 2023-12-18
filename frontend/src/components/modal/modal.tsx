import { Component, Element, h, Host, Listen, Prop } from '@stencil/core';

@Component({
  tag: 'mex-modal',
  styleUrl: 'modal.css',
})
export class ModalComponent {
  @Element() hostEl: HTMLElement;

  #closeRef?: HTMLButtonElement;

  @Prop() handleClose?: () => void;
  @Prop() handleSetFocus?: (closeRef?: HTMLButtonElement) => void;

  @Listen('focusin', { target: 'window' })
  focusinHandler(event: Event) {
    const target = event.target as HTMLElement;
    const captureFocus = !target.closest('mex-modal') && !target.closest('mex-header');
    captureFocus && this.setFocus();
  }

  @Listen('scrollModalToTop')
  scrollModalToTopHandler() {
    this.hostEl.scrollTo(0, 0);
  }

  close() {
    this.handleClose?.();
  }

  setFocus() {
    requestAnimationFrame(() => this.handleSetFocus?.(this.#closeRef));
  }

  onKeyUp({ key }: KeyboardEvent) {
    key === 'Escape' && this.close();
  }

  componentDidRender() {
    this.setFocus();
  }

  render() {
    return (
      <Host class="modal__backdrop" onClick={() => this.close()}>
        <div
          class="modal"
          role="dialog"
          aria-modal="true"
          onKeyUp={(event) => this.onKeyUp(event)}
          onClick={(event) => event.stopPropagation()}
        >
          <slot />
          {!!this.handleClose && (
            <button class="modal__close" ref={(el) => (this.#closeRef = el)} onClick={() => this.close()}>
              <mex-icon-close />
            </button>
          )}
        </div>
      </Host>
    );
  }
}

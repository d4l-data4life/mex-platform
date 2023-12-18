import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-close-inline',
})
export class IconCloseInlineComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path d="M4.2 4.2a11 11 0 1 1 15.6 15.6A11 11 0 0 1 4.2 4.2Zm11.3 3.2L12 11 8.5 7.4l-1 1.2 3.4 3.4-3.5 3.5 1.2 1 3.4-3.4 3.5 3.5 1-1.2-3.4-3.4 3.5-3.5-1.2-1z" />
        </svg>
      </Host>
    );
  }
}

import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-chevron',
})
export class IconChevronComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path d="m20 6 2 2.1L11.8 18 2 8.1 4 6l8 7.8Z" />
        </svg>
      </Host>
    );
  }
}

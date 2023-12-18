import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-close',
})
export class IconCloseComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path d="M21.6 2.5a1.4 1.4 0 0 1 .1 1.9l-.1.1-7.5 7.6 7.5 7.5a1.4 1.4 0 0 1-1.9 2.1l-.1-.1L12 14l-7.5 7.5a1.4 1.4 0 0 1-2.2-1.9l.1-.1 7.6-7.5-7.6-7.6a1.4 1.4 0 0 1 1.8-2.2h.1l.2.2 7.5 7.6 7.5-7.6a1.4 1.4 0 0 1 2 0z" />
        </svg>
      </Host>
    );
  }
}

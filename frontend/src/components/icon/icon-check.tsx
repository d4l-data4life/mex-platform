import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-check',
})
export class IconCheckComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path d="m19.5 4.2-11 11.2-4-4.2a.6.6 0 0 0-.8 0l-1.5 1.6a.6.6 0 0 0 0 .8l6 6.2c.2.3.6.3.8 0L21.8 6.6a.6.6 0 0 0 0-.9l-1.5-1.5a.6.6 0 0 0-.8 0z" />
        </svg>
      </Host>
    );
  }
}

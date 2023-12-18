import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-warning',
})
export class IconWarningComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path d="m21.75 18.35-8.32-14.5a1.66 1.66 0 0 0-2.9 0l-8.3 14.5C1.58 19.52 2.39 21 3.69 21H20.3c1.31 0 2.13-1.48 1.44-2.65zM11.12 9.4c0-.56.44-1.04 1-1.04.54 0 1 .46 1 1.04v4.37c0 .57-.44 1.05-1 1.05-.54 0-1-.46-1-1.05zm.98 8.66c-.56 0-1-.48-1-1.05 0-.59.46-1.04 1-1.04.56 0 1 .48 1 1.04 0 .57-.44 1.05-1 1.05z" />
        </svg>
      </Host>
    );
  }
}

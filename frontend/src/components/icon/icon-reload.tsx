import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-reload',
})
export class IconReloadComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
            <path d="M3 5v5h5" />
            <path d="M7 17A8 8 0 1 0 7 6l-4 4" />
          </g>
        </svg>
      </Host>
    );
  }
}

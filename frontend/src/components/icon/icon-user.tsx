import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-user',
})
export class IconUserComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5">
            <circle cx="12" cy="12" r="9.5" />
            <path d="M17.5 19.3v-1.7c0-1.8-2.2-3.3-3.8-3.3h-3.4c-1.6 0-3.8 1.5-3.8 3.3v1.6" />
            <circle cx="12" cy="9" r="2.8" />
          </g>
        </svg>
      </Host>
    );
  }
}

import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-search',
})
export class IconSearchComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <g fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2">
            <path d="m21 21-6-6" />
            <circle cx="10" cy="10" r="7" />
          </g>
        </svg>
      </Host>
    );
  }
}

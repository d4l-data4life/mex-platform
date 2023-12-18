import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-resource',
})
export class IconResourceComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path
            fill="var(--color-icon-indicator-background)"
            d="M12.75 4h-5C6.78 4 6 4.8 6 5.8v12.4c0 1 .78 1.8 1.75 1.8h8.5c.97 0 1.75-.8 1.75-1.8V8.4Z"
          />
          <g
            fill="none"
            stroke="var(--color-icon-indicator)"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width=".83"
          >
            <path d="m13.93 2.58 5.49 5.65V19.2c0 .61-.24 1.16-.63 1.56a2.13 2.13 0 0 1-1.54.66v0H6.75c-.6 0-1.14-.25-1.54-.66-.39-.4-.63-.95-.63-1.56v0V4.8c0-.61.24-1.16.63-1.56a2.13 2.13 0 0 1 1.54-.66Z" />
            <path d="M13.5 3v5.5h6m-4.5 4H8.33m6.67 3H8.33M10 8.42H8.33" />
          </g>
        </svg>
      </Host>
    );
  }
}

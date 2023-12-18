import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-source',
})
export class IconSourceComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path
            fill="var(--color-icon-indicator-background)"
            d="M20 17.24A1.8 1.8 0 0 1 18.17 19H5.83A1.8 1.8 0 0 1 4 17.24V6.76C4 5.8 4.47 5 5.83 5H8.1l1.83 2.64h8.23A1.8 1.8 0 0 1 20 9.41Z"
          />
          <path
            fill="none"
            stroke="var(--color-icon-indicator)"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width=".87"
            d="m9.55 3.44 1.82 2.72h7.9c.64 0 1.23.25 1.65.67.4.4.64.94.64 1.54v9.98c0 .6-.25 1.15-.64 1.54-.42.42-1 .67-1.65.67v0H4.73c-.64 0-1.23-.25-1.65-.67a2.82 2.39 0 0 1-.64-1.54v0-12.7c0-.64.18-1.21.58-1.61.36-.36.9-.6 1.7-.6h4.83Z"
          />
          <path
            fill="var(--color-icon-indicator)"
            fill-rule="nonzero"
            d="M12.58 10a.42.5 0 0 1 .08 1H6.42a.42.5 0 0 1-.08-1h.08Z"
          />
        </svg>
      </Host>
    );
  }
}

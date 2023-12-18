import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-email',
})
export class IconEmailComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path
            fill="var(--color-icon-indicator-background)"
            d="M7 6h10a3 3 0 0 1 3 3v6a3 3 0 0 1-3 3H7a3 3 0 0 1-3-3V9a3 3 0 0 1 3-3Z"
          />
          <path
            fill="var(--color-icon-indicator)"
            d="M19.27 4A2.7 2.7 0 0 1 22 6.67v10.66A2.7 2.7 0 0 1 19.27 20H4.73A2.7 2.7 0 0 1 2 17.33V6.67A2.7 2.7 0 0 1 4.73 4Zm1.88 3.93-8.63 5.9a.93.93 0 0 1-.93.07l-.11-.06-8.62-5.9v9.4c0 .46 1.24 1.77 1.79 1.9h14.62c.5 0 1.88-1.42 1.88-1.9v-9.4Zm-1.88-3.17H4.73c-.35 0-1.13.7-1.56 1.28a.7.7 0 0 1 .16.06l.1 1.04L12 13.11l8.57-5.97c.17-.12.38-.35.57-.54-.1-.53-1.39-1.84-1.87-1.84Z"
          />
        </svg>
      </Host>
    );
  }
}

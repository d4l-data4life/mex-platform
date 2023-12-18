import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-lock',
})
export class IconLockComponent {
  @Prop() classes = '';
  @Prop() open = false;

  render() {
    const { open } = this;

    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path
            fill={open ? 'var(--color-green-lighter)' : 'var(--color-red-lighter)'}
            d="M6.67 11h10.66c.92 0 1.67.75 1.67 1.67v5.66c0 .92-.75 1.67-1.67 1.67H6.67C5.75 20 5 19.25 5 18.33v-5.66c0-.92.75-1.67 1.67-1.67z"
          />

          <g fill="none" stroke="var(--color-icon-indicator)">
            <rect width="17" height="12" x="3.5" y="9.5" rx="1.7" />
            {!open && <path d="M7.5 10V6.7c0-2.3 2-4.2 4.5-4.2s4.5 1.9 4.5 4.2V10" />}
            {open && <path d="M7.5 10V7c0-2.3 1.71-4.19 3.93-4.49a4.74 5.69 0 0 1 3.23 1.2" />}
          </g>
        </svg>
      </Host>
    );
  }
}

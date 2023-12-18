import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-platform',
})
export class IconPlatformComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path
            fill="var(--color-icon-indicator-background)"
            d="M7.1 6.97 4 12v4.25c0 .97.8 1.75 1.8 1.75h12.4c1 0 1.8-.78 1.8-1.75V12l-3.1-5.03A1.8 1.8 0 0 0 15.27 6H8.72c-.69 0-1.31.38-1.62.97z"
          />
          <path
            fill="none"
            stroke="var(--color-icon-indicator)"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width=".83"
            d="M20.67 11.5H4m12.28-6.92a2.25 2.25 0 0 1 1.98 1.2v0l3.16 6.12v5.35c0 .6-.25 1.14-.66 1.54-.4.39-.95.63-1.56.63v0H4.8c-.61 0-1.16-.24-1.56-.63a2.13 2.13 0 0 1-.66-1.54v0-5.35l3.15-6.12c.2-.37.48-.67.82-.87a2.25 2.25 0 0 1 1.17-.33v0zM7.5 15.5h1m2 0h1"
          />
        </svg>
      </Host>
    );
  }
}

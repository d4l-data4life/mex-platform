import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-info',
})
export class IconInfoComponent {
  @Prop() classes = '';
  @Prop() hasBackground = false;

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          {this.hasBackground && <circle fill="var(--color-icon-indicator-background)" cx="12" cy="12" r="7" />}
          <path d="M12 2a10 10 0 1 1 0 20 10 10 0 0 1 0-20zm0 1.5a8.5 8.5 0 1 0 0 17 8.5 8.5 0 0 0 0-17zm0 7.5c.5 0 1 .5 1 1v4.8c0 .7-.4 1.2-1 1.2a1 1 0 0 1-1-1v-4.8c0-.7.5-1.2 1-1.2Zm0-4c.6 0 1 .4 1 .9s-.4.8-.9.9H12c-.6 0-1-.4-1-.9s.4-.8.9-.9Z" />
        </svg>
      </Host>
    );
  }
}

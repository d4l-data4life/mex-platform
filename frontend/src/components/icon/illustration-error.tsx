import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-illustration-error',
})
export class IllustrationErrorComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host
        aria-hidden="true"
        role="presentation"
        class={`illustration ${this.classes}`}
        style={{ '--illustration-width': '64px', '--illustration-height': '80px' }}
      >
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 80">
          <rect width="40" height="8" x="12" y="72" fill="var(--color-gray-4)" rx="4" />
          <path
            fill="var(--color-primary-1)"
            d="M24.7 0c13.6 0 24.7 11 24.7 24.7a24.6 24.6 0 0 1-5.9 16l19.9 19.8a2 2 0 0 1-2.7 3l-.2-.1-19.8-19.9a24.6 24.6 0 0 1-16 5.9 24.7 24.7 0 0 1 0-49.4Zm0 4a20.6 20.6 0 1 0 0 41.3 20.6 20.6 0 0 0 0-41.2z"
          />
          <path
            fill="var(--color-primary-1)"
            d="M16.8 13.6a1 1 0 0 1 1 .9v.1l.4 5.8 7.9 1.2a1 1 0 0 1 .8.7l.8 3.9 7.8 2a1 1 0 0 1-.4 2H35l-6.8-1.8 1.8 8.3a1 1 0 0 1-2 .6V37l-2-9.3a1 1 0 0 1-.4-1.8l-.5-2.5h-.5a1 1 0 0 1 0 .4l-3 7.8a1 1 0 0 1-2-.5l.1-.2 3-7.7V23l-5.6-.8a1 1 0 0 1-.9-.8v-.2l-.4-6.6a1 1 0 0 1 1-1z"
          />
        </svg>
      </Host>
    );
  }
}

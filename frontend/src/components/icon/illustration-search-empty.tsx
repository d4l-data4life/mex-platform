import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-illustration-search-empty',
})
export class IllustrationSearchEmptyComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host
        aria-hidden="true"
        role="presentation"
        class={`illustration ${this.classes}`}
        style={{ '--illustration-width': '44px', '--illustration-height': '56px' }}
      >
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 44 56">
          <rect width="28" height="6" x="8" y="50" fill="var(--color-gray-4)" rx="2.8" />
          <path
            fill="var(--color-primary-1)"
            d="M16.99.01a16.99 16.987 0 0 1 12.966 27.978l13.676 13.664a1.402 1.401 0 0 1-1.883 2.062l-.1-.09L27.983 29.95a16.92 16.917 0 0 1-10.993 4.034A16.995 16.992 0 0 1 16.99 0zm0 2.803a14.197 14.194 0 1 0 0 28.388 14.202 14.2 0 0 0 0-28.398zm4.375 8.378a1.201 1.201 0 0 1 1.602-.08l.1.08c.421.43.461 1.101.08 1.602l-.09.1-4.245 4.234 4.245 4.244c.441.44.471 1.141.08 1.602l-.08.1a1.201 1.201 0 0 1-1.601.08l-.09-.09-4.245-4.244-4.245 4.244a1.201 1.201 0 0 1-1.602.08l-.1-.08a1.201 1.201 0 0 1-.08-1.602l.09-.09 4.245-4.244-4.245-4.244a1.201 1.201 0 0 1-.08-1.602l.08-.09a1.201 1.201 0 0 1 1.461-.18l.14.1.1.08 4.236 4.244z"
          />
        </svg>
      </Host>
    );
  }
}

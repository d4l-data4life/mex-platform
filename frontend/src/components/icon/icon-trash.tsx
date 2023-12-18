import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-trash',
})
export class IconTrashComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path d="M14.4 2c1.48 0 2.74 1.02 2.74 2.37v.68h3.01a.85.85 0 1 1 0 1.7h-.17v12.57A2.71 2.71 0 0 1 17.24 22H7.78c-1.5 0-2.74-1.2-2.74-2.68V6.75h-1.2a.85.85 0 1 1 0-1.7h4.04v-.68C7.88 3.02 9.14 2 10.62 2Zm3.88 4.75H6.74v12.57c0 .54.46.99 1.04.99h9.46c.58 0 1.04-.45 1.04-.99zm-7.8 3.39c.46 0 .84.38.84.84v5.09a.85.85 0 1 1-1.7 0v-5.09c0-.46.38-.84.85-.84zm4.07 0c.47 0 .85.38.85.84v5.09a.85.85 0 1 1-1.7 0v-5.09c0-.46.38-.84.85-.84zm-.15-6.45h-3.78c-.61 0-1.05.35-1.05.68v.68h5.87v-.68c0-.33-.43-.68-1.04-.68z" />
        </svg>
      </Host>
    );
  }
}

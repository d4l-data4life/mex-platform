import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-datum',
})
export class IconDatumComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path fill="var(--color-icon-indicator-background)" d="M4 10h16v5H4Zm3-5h10v4H7Zm0 11h10v4H7Z" />
          <path
            fill="var(--color-icon-indicator)"
            d="M17.9 3c.86 0 1.1.7 1.1 1.56V9.2h1.44c.86 0 1.56.7 1.56 1.57v3.51c0 .86-.7 1.56-1.56 1.56H19v4.6c0 .87-.24 1.57-1.1 1.57H6.1c-.86 0-1.1-.7-1.1-1.56v-4.6H3.56c-.86 0-1.56-.7-1.56-1.57v-3.51c0-.87.7-1.56 1.56-1.56H5V4.56C5 3.7 5.24 3 6.1 3Zm.3 13H5.8v4.6c0 .22.3.57.52.57h11.36c.23 0 .52-.35.52-.57zm-2.77 2c.31 0 .57.22.57.5s-.26.5-.57.5H8.57c-.31 0-.57-.22-.57-.5s.26-.5.57-.5Zm5.17-8H3.4c-.22 0-.4.2-.4.45v4.1c0 .25.18.45.4.45h17.2a.37.37 0 0 0 .28-.13.49.49 0 0 0 .12-.32v-4.1c0-.25-.18-.45-.4-.45ZM8.5 12a.5.5 0 1 1 0 1h-3a.5.5 0 1 1 0-1zm3 0a.5.5 0 1 1 0 1h-1a.5.5 0 1 1 0-1zm7 0a.5.5 0 1 1 0 1h-5a.5.5 0 1 1 0-1zm-.82-8.2H6.32c-.23 0-.52.35-.52.56V9h12.4V4.36c0-.21-.3-.56-.52-.56ZM15.43 6c.31 0 .57.22.57.5s-.26.5-.57.5H8.57C8.26 7 8 6.78 8 6.5s.26-.5.57-.5Z"
          />
        </svg>
      </Host>
    );
  }
}

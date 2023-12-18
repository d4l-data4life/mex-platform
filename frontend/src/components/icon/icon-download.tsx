import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-download',
})
export class IconDownloadComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path d="M2.69 12a.7.7 0 0 1 .69.72v4.63a3.16 3.16 0 0 0 3.08 3.21h11.09c1.7 0 3.07-1.43 3.08-3.21v-4.53a.73.73 0 0 1 .2-.5.67.67 0 0 1 .48-.21.7.7 0 0 1 .69.71v4.53c0 2.57-2 4.65-4.46 4.65H6.46A4.56 4.56 0 0 1 2 17.35v-4.63c0-.4.3-.72.69-.72ZM12 2c.38 0 .68.29.68.65v12.13l3.16-3a.7.7 0 0 1 .96 0 .62.62 0 0 1 0 .91l-4.32 4.12A.67.67 0 0 1 12 17a.7.7 0 0 1-.48-.2L7.2 12.7a.63.63 0 0 1-.2-.46.63.63 0 0 1 .2-.46.67.67 0 0 1 .48-.19.7.7 0 0 1 .48.2l3.16 3V2.65c0-.36.3-.65.68-.65Z" />
        </svg>
      </Host>
    );
  }
}

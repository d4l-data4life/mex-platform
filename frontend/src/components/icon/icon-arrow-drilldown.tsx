import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-arrow-drilldown',
})
export class IconArrowDrilldownComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path d="m17.5 13.3-7 7 1.7 1.7L22 12 12.3 2l-1.7 1.8 6.8 7H7.9c-.6 0-3.4-1.9-3.4-2.5V2H2v6.5c0 2 3 4.8 4.8 4.8z" />
        </svg>
      </Host>
    );
  }
}

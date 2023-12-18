import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-arrow',
})
export class IconArrowComponent {
  @Prop() classes = '';

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <path d="m12 2-1.8 1.8 7 7H2v2.5h15.3l-7 7L12 22l10-10Z" />
        </svg>
      </Host>
    );
  }
}

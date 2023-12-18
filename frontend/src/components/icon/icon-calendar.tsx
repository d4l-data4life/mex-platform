import { Component, Prop, Host, h } from '@stencil/core';

@Component({
  tag: 'mex-icon-calendar',
})
export class IconCalendarComponent {
  @Prop() classes = '';
  @Prop() arrows = false;

  render() {
    return (
      <Host aria-hidden="true" role="presentation" class={`icon ${this.classes}`}>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class={`icon ${this.classes}`}>
          <rect width="14" height="15" x="5" y="5" fill="var(--color-icon-indicator-background)" rx="1.67" />
          <g fill="none" stroke="var(--color-icon-indicator)">
            <rect width="17" height="18" x="3.5" y="3.5" rx="2" />
            <path d="M8.26 2v3.48M15.74 2v3.48M4 8.53h16" />
          </g>

          {this.arrows && (
            <path
              fill="var(--color-icon-indicator)"
              d="m14.5 11-.61.61 2.45 2.45H12v.88h4.34l-2.45 2.45.61.61 3.5-3.5Zm-5 0 .61.61-2.45 2.45H12v.88H7.66l2.45 2.45-.61.61L6 14.5Z"
            />
          )}
        </svg>
      </Host>
    );
  }
}

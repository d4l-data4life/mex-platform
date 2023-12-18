import { Component, h, Host, Prop } from '@stencil/core';

@Component({
  tag: 'mex-placeholder',
  styleUrl: 'placeholder.css',
  shadow: true,
})
export class PlaceholderComponent {
  @Prop() lines = 1;
  @Prop() inline = false;
  @Prop() text = '';
  @Prop() width?: string;
  @Prop() height?: string;

  render() {
    return (
      <Host class={{ placeholder: true, 'placeholder--inline': this.inline }}>
        {new Array(this.lines).fill(null).map(() => {
          const width = this.width ?? this.text ? 'auto' : `${50 + Math.floor(Math.random() * 40)}%`;
          const animationDuration = 1500 + Math.floor(Math.random() * 1500);

          return (
            <div
              aria-hidden="true"
              role="presentation"
              class="placeholder__line"
              style={{ width, height: this.height ?? 'auto', '--animation-duration': `${animationDuration}ms` }}
              innerHTML={(this.text || '-') + '<br />'}
            />
          );
        })}
      </Host>
    );
  }
}

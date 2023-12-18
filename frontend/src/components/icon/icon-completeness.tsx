import { Component, Prop, Host, h, Fragment } from '@stencil/core';

@Component({
  tag: 'mex-icon-completeness',
  styleUrl: 'icon-completeness.css',
})
export class IconCloseComponent {
  @Prop() value: number;
  @Prop() classes: string = '';
  @Prop() barsCount = 5;
  @Prop() monochrome = false;

  get activeColor() {
    const { value, monochrome } = this;
    if (monochrome) {
      return '--color-highlight-1';
    }

    if (value < 20) {
      return '--color-red';
    }

    if (value < 60) {
      return '--color-yellow';
    }

    return '--color-green';
  }

  get activePercentage() {
    const barPercentage = Math.ceil(100 / this.barsCount);
    const barCount = Math.ceil(this.value / barPercentage);
    return 100 - barCount * barPercentage;
  }

  getBars() {
    return (
      <Fragment>
        {new Array(this.barsCount).fill(null).map((_, index) => (
          <span key={index} class="icon-completeness__bar" />
        ))}
      </Fragment>
    );
  }

  render() {
    return (
      <Host class={`icon-completeness ${this.classes}`}>
        <div class="icon-completeness__bars">{this.getBars()}</div>
        <div
          class="icon-completeness__bars icon-completeness__bars--active"
          style={{
            color: `var(${this.activeColor})`,
            '--icon-completeness-active-height': `${this.activePercentage}%`,
          }}
        >
          {this.getBars()}
        </div>
      </Host>
    );
  }
}

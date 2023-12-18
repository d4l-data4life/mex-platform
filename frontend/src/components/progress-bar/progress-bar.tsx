import { Component, h, Host, Prop } from '@stencil/core';

@Component({
  tag: 'mex-progress-bar',
  styleUrl: 'progress-bar.css',
})
export class ProgressBarComponent {
  @Prop() progress: number = 0;
  @Prop() spaced: boolean = false;

  render() {
    return (
      <Host class={{ 'progress-bar': true, 'progress-bar--spaced': this.spaced }}>
        <div class="progress-bar__progress" style={{ width: `${this.progress}%` }} />
      </Host>
    );
  }
}

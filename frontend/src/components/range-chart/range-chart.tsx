import { Component, h, Prop, State } from '@stencil/core';
import stores from 'stores';

export interface RangeChartPoint {
  count: number;
  value: number;
}

@Component({
  tag: 'mex-range-chart',
  styleUrl: 'range-chart.css',
})
export class RangeChartComponent {
  @Prop() points: RangeChartPoint[];
  @Prop() value: number[];
  @Prop() handleClick?: (point: RangeChartPoint) => void;

  @State() hoveredIndex?: number;

  get min() {
    return this.points?.[0]?.value;
  }

  get max() {
    return this.points?.[this.points.length - 1]?.value;
  }

  get offsetLeft() {
    const { value, min, points } = this;
    return ((((value?.[0] ?? min) - min) / points.length) * 100).toFixed(4);
  }

  get offsetRight() {
    const { value, min, max, points } = this;
    return (100 - (((value?.[1] ?? max) + 1 - min) / points.length) * 100).toFixed(4);
  }

  get percentages() {
    const counts = (this.points ?? []).map(({ count }) => count);
    const max = Math.max(...counts);
    return counts.map((count) => Math.round((count / max) * 100));
  }

  renderContent(percentages: number[], active: boolean) {
    return (
      <div
        class={{
          'range-chart__bars': true,
          'range-chart__bars--active': active,
        }}
        style={{
          '--range-chart-offset-left': active ? `${this.offsetLeft}%` : '0',
          '--range-chart-offset-right': active ? `${this.offsetRight}%` : '0',
        }}
      >
        {percentages.map((perc, index) => (
          <div key={index} class="range-chart__bar" style={{ height: `${perc}%` }} />
        ))}
      </div>
    );
  }

  render() {
    const { points, percentages } = this;

    return (
      <div class="range-chart" aria-hidden="true" role="presentation">
        <div class="range-chart__bars range-chart__bars--help">
          {percentages.map((_, index) => (
            <mex-tooltip
              key={index}
              class="range-chart__tooltip"
              text={stores.i18n.t('filters.chartTooltip', {
                value: points[index]?.value,
                count: points[index]?.count,
              })}
              method="hover"
              handleClick={this.handleClick ? () => this.handleClick(points[index]) : null}
            >
              <span slot="toggle" />
            </mex-tooltip>
          ))}
        </div>
        {this.renderContent(percentages, true)}
        {this.renderContent(percentages, false)}
      </div>
    );
  }
}

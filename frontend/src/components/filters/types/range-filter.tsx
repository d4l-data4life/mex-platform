import { Component, Event, EventEmitter, Fragment, h, Prop, State, Watch } from '@stencil/core';
import { ANALYTICS_CUSTOM_EVENTS } from 'config';
import { SearchFacet } from 'config/search';
import services from 'services';
import stores from 'stores';
import { SearchResultsFacetBucket } from 'stores/search';
import { buildAscNumSequence } from 'utils/search';

@Component({
  tag: 'mex-range-filter',
  styleUrl: 'range-filter.css',
})
export class RangeFilterComponent {
  @Prop() busy: boolean;
  @Prop() expanded: boolean;
  @Prop() facet: SearchFacet;
  @Prop() buckets: SearchResultsFacetBucket[];

  @State() value: number[];
  @State() inputValue: number[];

  @Event() updateSearch: EventEmitter;

  @Watch('buckets')
  onBucketsChange(newBuckets: SearchResultsFacetBucket[], oldBuckets: SearchResultsFacetBucket[]) {
    if (
      newBuckets?.[0]?.value !== oldBuckets?.[0]?.value ||
      newBuckets?.[newBuckets.length - 1]?.value !== oldBuckets?.[oldBuckets.length - 1]?.value
    ) {
      this.adjustValue();
    }
  }

  constructor() {
    this.onFilterChange = this.onFilterChange.bind(this);
  }

  get filter() {
    return stores.filters.get(this.facet.axis.name);
  }

  get filterValue(): number[] {
    const value = this.filter?.[0]
      ?.split('-')
      .map((year) => parseInt(year, 10))
      .filter((year) => !isNaN(year));

    return value?.length !== 2 ? [this.min, this.max] : value;
  }

  get canChange() {
    return this.expanded && !this.busy;
  }

  get hasBuckets() {
    return !!this.buckets?.length;
  }

  get isEmpty() {
    return !this.busy && !this.hasBuckets;
  }

  get min() {
    return parseInt(this.buckets?.[0]?.value, 10);
  }

  get max() {
    return parseInt(this.buckets?.[this.buckets.length - 1]?.value, 10);
  }

  get chartPoints() {
    return this.buckets?.map(({ value, count }) => ({
      value: parseInt(value, 10) as number,
      count,
    }));
  }

  onFilterChange() {
    this.adjustValue();
  }

  adjustValue(rawValue = this.filterValue, index = -1): number[] {
    const { min, max } = this;
    const value = buildAscNumSequence(rawValue, index, min, max);
    this.value = value;
    this.inputValue = value;
    return value;
  }

  handleChange(rawValue: number[], index = -1) {
    const value = this.adjustValue(rawValue, index);
    const { name } = this.facet.axis;
    stores.filters.set(name, [`${value[0]}-${value[1]}`]);
    this.updateSearch.emit();

    services.analytics.trackEvent(...ANALYTICS_CUSTOM_EVENTS.SEARCH_FILTER, `Set range: ${name}`);
  }

  handleInputChange(event: InputEvent, index: number) {
    const { inputValue } = this;
    const newValue = parseInt((event.target as HTMLInputElement).value, 10);
    inputValue[index] = isNaN(newValue) ? 0 : newValue;
    this.handleChange(inputValue, index);
  }

  renderInput(index: number) {
    return (
      <input
        class="range-filter__input"
        type="number"
        aria-label={stores.i18n.t(`filters.${index ? 'to' : 'from'}`)}
        readonly={!this.canChange}
        disabled={!this.expanded}
        min={this.min}
        max={this.max}
        placeholder={stores.i18n.t(`filters.${index ? 'to' : 'from'}`)}
        step={1}
        value={this.inputValue[index]}
        onChange={(event: InputEvent) => this.handleInputChange(event, index)}
        onKeyDown={(event: KeyboardEvent) => !this.canChange && event.preventDefault()}
        data-test="filter:input"
        data-test-context={this.facet.axis}
        data-test-key={index}
      />
    );
  }

  componentWillLoad() {
    this.adjustValue();
    stores.filters.addListener(this.facet.axis.name, this.onFilterChange);
  }

  disconnectedCallback() {
    stores.filters.removeListener(this.onFilterChange);
  }

  render() {
    const { min, max, canChange, hasBuckets, isEmpty } = this;

    return (
      <div class={{ 'range-filter': true, 'range-filter--empty': isEmpty }}>
        {isEmpty && stores.i18n.t('filters.empty')}

        {hasBuckets && (
          <Fragment>
            <div class="range-filter__inputs">
              {this.renderInput(0)}-{this.renderInput(1)}
            </div>

            <mex-range-chart
              class="range-filter__chart"
              points={this.chartPoints}
              value={this.inputValue}
              handleClick={canChange ? (point) => this.handleChange([point.value, point.value]) : null}
              data-test="filter:chart"
            />

            <mex-range-slider
              disabled={!canChange}
              min={min}
              max={max}
              value={this.value}
              handleChange={(value) => this.handleChange(value)}
              handleDrag={(value) => requestAnimationFrame(() => (this.inputValue = value))}
              mode="frame"
              highlightActiveRange
              testAttr="filter:slider"
            />
          </Fragment>
        )}
      </div>
    );
  }
}

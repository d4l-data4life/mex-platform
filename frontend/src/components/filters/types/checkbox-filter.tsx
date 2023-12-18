import { Component, Event, EventEmitter, h, Prop, State } from '@stencil/core';
import { ANALYTICS_CUSTOM_EVENTS } from 'config';
import { SearchFacet } from 'config/search';
import services from 'services';
import stores from 'stores';
import { SearchResultsFacetBucket } from 'stores/search';
import { formatValue } from 'utils/field';
import { CheckStatusEnum } from 'utils/filters';

@Component({
  tag: 'mex-checkbox-filter',
  styleUrl: 'checkbox-filter.css',
})
export class CheckboxFilterComponent {
  @Prop() busy: boolean;
  @Prop() expanded: boolean;
  @Prop() facet: SearchFacet;
  @Prop() buckets: SearchResultsFacetBucket[];

  @State() isLoadingMore = false;

  @Event() updateSearch: EventEmitter;

  get hasBuckets() {
    return !!this.buckets?.length;
  }

  get canToggle() {
    return this.expanded && !this.busy;
  }

  get axis() {
    return this.facet.axis;
  }

  get isEmpty() {
    return !this.busy && !this.hasBuckets;
  }

  get state() {
    return stores.filters.get(this.axis.name);
  }

  toggle(value: string, previouslyChecked: CheckStatusEnum) {
    const { axis } = this;
    const isReallyChecked = previouslyChecked === CheckStatusEnum.UNCHECKED;
    isReallyChecked ? stores.filters.add(axis.name, value) : stores.filters.remove(axis.name, value);
    this.updateSearch.emit();

    services.analytics.trackEvent(
      ...ANALYTICS_CUSTOM_EVENTS.SEARCH_FILTER,
      `Switched ${isReallyChecked ? 'on' : 'off'}: ${axis}`
    );
  }

  render() {
    const { axis, busy, canToggle, facet, isEmpty, hasBuckets, state } = this;

    return (
      <div class={{ 'checkbox-filter': true, 'checkbox-filter--empty': isEmpty }}>
        {isEmpty && stores.i18n.t('filters.empty')}

        {hasBuckets &&
          this.buckets.map(({ value, count, hierarchyInfo }, index) => (
            <mex-checkbox
              key={`${axis}--${value}`}
              classes="checkbox-filter__item"
              label={formatValue([hierarchyInfo?.display || value], axis.uiField)}
              secondaryText={`(${count})`}
              handleChange={(previouslyChecked) => this.toggle(value, previouslyChecked)}
              checked={state.includes(value) ? CheckStatusEnum.CHECKED : CheckStatusEnum.UNCHECKED}
              disabled={!this.canToggle || (!count && !state.includes(value))}
              data-test="filter:item"
              data-test-context={axis}
              data-test-key={index}
              data-test-active={String(state.includes(value))}
              testAttr="filter:item:checkbox"
            />
          ))}

        <mex-filter-footer busy={busy} facet={facet} hasItems={hasBuckets} isLoadMoreDisabled={!canToggle} />
      </div>
    );
  }
}

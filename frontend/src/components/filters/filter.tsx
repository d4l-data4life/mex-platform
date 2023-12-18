import { Component, h, Host, Prop, State } from '@stencil/core';
import { ANALYTICS_CUSTOM_EVENTS, FIELDS } from 'config';
import { SearchFacet } from 'config/search';
import services from 'services';
import stores from 'stores';
import { translateFieldName } from 'utils/field';

@Component({
  tag: 'mex-filter',
  styleUrl: 'filter.css',
})
export class FilterComponent {
  @Prop() isParentExpanded = true;
  @Prop() busy = true;
  @Prop() facet: SearchFacet;

  @State() isExpanded = true;

  get facetResponse() {
    return stores.search.facets?.find((item) => item.axis === this.facet.axis.name);
  }

  get buckets() {
    return this.facetResponse?.buckets;
  }

  get expanded() {
    return this.isParentExpanded && this.isExpanded;
  }

  get type() {
    return this.facet.type;
  }

  render() {
    const { type, busy, expanded, facet, facetResponse, buckets, isParentExpanded } = this;

    return (
      <Host class="filter" data-test="filter" data-test-context={this.facet.axis}>
        <button
          class="filter__headline"
          onClick={() => {
            this.isExpanded = !this.isExpanded;

            services.analytics.trackEvent(
              ...ANALYTICS_CUSTOM_EVENTS.SEARCH_NAVIGATION,
              this.isExpanded ? 'Filter expanded' : 'Filter collapsed'
            );
          }}
          disabled={!isParentExpanded}
          aria-expanded={String(this.expanded)}
          data-test="filter:toggle"
          data-test-active={String(this.expanded)}
        >
          <span>{translateFieldName(facet.axis.uiField ?? FIELDS[facet.axis.name])}</span>
          <mex-icon-chevron
            classes={`icon--inline icon--mirrorable ${this.isExpanded ? 'icon--mirrored-vertical' : ''}`}
          />
        </button>
        <mex-accordion expanded={this.isExpanded}>
          {type === 'exact' && <mex-checkbox-filter busy={busy} expanded={expanded} facet={facet} buckets={buckets} />}
          {type === 'yearRange' && <mex-range-filter busy={busy} expanded={expanded} facet={facet} buckets={buckets} />}
          {type === 'hierarchy' && (
            <mex-hierarchy-filter busy={busy} expanded={expanded} facet={facet} facetResponse={facetResponse} />
          )}
        </mex-accordion>
      </Host>
    );
  }
}

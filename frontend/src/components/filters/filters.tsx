import { Component, Event, EventEmitter, h, Host, State } from '@stencil/core';
import { ANALYTICS_CUSTOM_EVENTS, SEARCH_CONFIG } from 'config';
import services from 'services';
import stores from 'stores';
import { IS_DESKTOP, mobileViewportChanges } from 'utils/device';

@Component({
  tag: 'mex-filters',
  styleUrl: 'filters.css',
})
export class FiltersComponent {
  @State() isExpanded = IS_DESKTOP;
  @State() canToggle: boolean;

  @Event() updateSearch: EventEmitter;

  constructor() {
    this.expandIfNotMobile = this.expandIfNotMobile.bind(this);
  }

  get canReset() {
    return !stores.filters.isEmpty;
  }

  reset() {
    stores.filters.reset();
    this.updateSearch.emit();

    services.analytics.trackEvent(...ANALYTICS_CUSTOM_EVENTS.SEARCH_NAVIGATION, 'Reset');
  }

  expandIfNotMobile(isMobileViewport: boolean) {
    this.canToggle = isMobileViewport;
    if (!isMobileViewport && !this.isExpanded) {
      this.isExpanded = true;
    }
  }

  componentWillLoad() {
    mobileViewportChanges.addListener(this.expandIfNotMobile);
  }

  disconnectedCallback() {
    mobileViewportChanges.removeListener(this.expandIfNotMobile);
  }

  render() {
    const { isExpanded, canToggle, canReset } = this;
    const { isBusy } = stores.search;

    return (
      <Host class="filters" data-test="filters">
        <div class="filters__headline">
          <button
            class="filters__toggle"
            disabled={!canToggle}
            onClick={() => (this.isExpanded = !isExpanded)}
            aria-expanded={String(isExpanded)}
            data-test="filters:toggle"
            data-test-active={String(isExpanded)}
          >
            <span class="u-underline-2">{stores.i18n.t('filters.title')}</span>
            {canToggle && (
              <mex-icon-chevron
                classes={`icon--inline icon--mirrorable ${isExpanded ? 'icon--mirrored-vertical' : ''}`}
              />
            )}
          </button>
          <button class="filters__reset" disabled={!canReset} onClick={() => this.reset()} data-test="filters:reset">
            <mex-icon-reload classes="icon--inline icon--large" />
            <span>{stores.i18n.t('filters.reset')}</span>
          </button>
        </div>
        <mex-accordion expanded={isExpanded}>
          <div class="filters__items">
            {SEARCH_CONFIG.FACETS.map((facet) => (
              <mex-filter isParentExpanded={isExpanded} busy={isBusy} facet={facet} />
            ))}
          </div>
        </mex-accordion>
      </Host>
    );
  }
}

import { Component, h, Host } from '@stencil/core';
import stores from 'stores';
import { AnalyticsConsent } from 'stores/analytics';

@Component({
  tag: 'mex-tracking-banner',
  styleUrl: 'tracking-banner.css',
})
export class TrackingBannerComponent {
  acceptTracking() {
    stores.analytics.consents = [AnalyticsConsent.TRACKING];
  }

  rejectTracking() {
    stores.analytics.consents = [];
  }

  render() {
    const { t, convertLinks } = stores.i18n;

    return (
      <Host class="tracking-banner" data-test="trackingBanner">
        <div class="tracking-banner__inner">
          <p
            class="tracking-banner__text"
            data-test="trackingBanner:text"
            innerHTML={convertLinks(t('consent.analytics.text'))}
          />
          <button
            class="button button--primary tracking-banner__action"
            onClick={() => this.rejectTracking()}
            data-test="trackingBanner:reject"
          >
            {t('consent.analytics.reject')}
          </button>
          <button
            class="button button--primary tracking-banner__action"
            onClick={() => this.acceptTracking()}
            data-test="trackingBanner:accept"
          >
            {t('consent.analytics.accept')}
          </button>
        </div>
      </Host>
    );
  }
}

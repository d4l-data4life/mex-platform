import { ANALYTICS_HEARTBEAT_TRACKING_INTERVAL, ANALYTICS_CONFIG, ANALYTICS_IS_ENABLED } from 'config';
import {
  AnalyticsServiceProvider,
  MatomoFlag,
  MatomoParams,
  MATOMO_CUSTOM_DIMENSION_PRODUCT,
  MATOMO_CUSTOM_DIMENSION_PROJECT,
  MATOMO_URL,
} from 'config/analytics';
import stores from 'stores';
import { AnalyticsConsent } from 'stores/analytics';
import { redactUrl } from 'utils/analytics';
import { generateState } from 'utils/auth';
import { get } from 'utils/fetch-client';

export interface AnalyticsEvent {
  category: string;
  action: string;
  name?: string;
  value?: number;
}

export interface AnalyticsImpression {
  title: string;
  url: string;
  referrer: string;
}

export class AnalyticsService {
  #heartbeatTrackingInterval: number;

  constructor() {
    stores.analytics.onConsentsChange(() => this.trackImpression());
  }

  canTrack() {
    return (
      ANALYTICS_IS_ENABLED() &&
      stores.analytics.consents.includes(AnalyticsConsent.TRACKING) &&
      stores.auth.isAuthenticated
    );
  }

  get impression() {
    return { title: document.title, url: redactUrl(document.location.href), referrer: redactUrl(document.referrer) };
  }

  async trackImpression() {
    await this.track(this.impression, null, false);
    this.scheduleHeartbeatTracking();
  }

  async trackEvent(category: string, action: string, name?: string, value?: number) {
    await this.track(this.impression, { category, action, name, value }, false);
  }

  async trackHeartbeat() {
    await this.track(this.impression, null, true);
  }

  scheduleHeartbeatTracking() {
    window.clearInterval(this.#heartbeatTrackingInterval);

    if (!this.canTrack()) {
      return;
    }

    this.#heartbeatTrackingInterval = window.setInterval(
      this.trackHeartbeat.bind(this),
      ANALYTICS_HEARTBEAT_TRACKING_INTERVAL
    );
  }

  private async track(impression: AnalyticsImpression, event: AnalyticsEvent, ping: boolean): Promise<void> {
    if (!this.canTrack()) {
      return; // Tracking disabled, rejected or not allowed (yet)
    }

    switch (ANALYTICS_CONFIG.PROVIDER) {
      case AnalyticsServiceProvider.MATOMO:
        try {
          await this.trackMatomo(impression, event, ping);
        } catch (_) {}
        break;
      default:
        throw new Error('Unknown or unset analytics service provider');
    }
  }

  private async trackMatomo(impression: AnalyticsImpression, event: AnalyticsEvent, ping: boolean): Promise<void> {
    const [hour, minute, second] = new Date().toLocaleTimeString().split(':');
    const { visitorId } = stores.analytics;

    const params = {
      ...(ping
        ? {
            [MatomoParams.PING]: String(MatomoFlag.YES),
          }
        : {}),
      [MatomoParams.RECORD]: String(MatomoFlag.YES),
      [MatomoParams.API_VERSION]: String(1),
      [MatomoParams.NONCE]: generateState(),
      [MatomoParams.SITE_ID]: String(ANALYTICS_CONFIG.SITE_ID),
      [MatomoParams.SEND_IMAGE]: String(MatomoFlag.NO),
      [MatomoParams.USER_ID]: visitorId,
      [MatomoParams.VISITOR_ID]: visitorId,
      [MatomoParams.LOCAL_TIME_HOUR]: String(parseInt(hour, 10)),
      [MatomoParams.LOCAL_TIME_MINUTE]: String(parseInt(minute, 10)),
      [MatomoParams.LOCAL_TIME_SECOND]: String(parseInt(second, 10)),
      [MatomoParams.SCREEN_RESOLUTION]: `${screen.width}x${screen.height}`,
      [MatomoParams.URL]: impression.url,
      [MatomoParams.REFERRER_URL]: impression.referrer,
      [`${MatomoParams.CUSTOM_DIMENSION}2`]: String(MATOMO_CUSTOM_DIMENSION_PRODUCT),
      [`${MatomoParams.CUSTOM_DIMENSION}3`]: String(MATOMO_CUSTOM_DIMENSION_PROJECT),
      ...(event
        ? {
            [MatomoParams.EVENT_CATEGORY]: event.category,
            [MatomoParams.EVENT_ACTION]: event.action,
            ...(event.name ? { [MatomoParams.EVENT_NAME]: event.name } : {}),
            ...(event.value ? { [MatomoParams.EVENT_VALUE]: String(event.value) } : {}),
          }
        : {
            [MatomoParams.PAGE_TITLE]: impression.title,
          }),
    };

    await get({
      url: `${MATOMO_URL}?${new URLSearchParams(params).toString()}`,
      sendLanguage: false,
      authorized: false,
      contentType: null,
    });
  }
}

export default new AnalyticsService();

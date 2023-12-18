jest.mock('stencil-router-v2');

import * as config from 'config';
import analytics, { AnalyticsConsent } from './analytics';
import * as analyticsUtils from 'utils/analytics';

describe('analytics store', () => {
  beforeEach(() => analytics.reset());

  it('sets and gets (default: []) the consents', () => {
    expect(analytics.consents).toEqual([]);
    analytics.consents = [AnalyticsConsent.TRACKING];
    expect(analytics.consents).toEqual([AnalyticsConsent.TRACKING]);
  });

  it('returns if user has made a choice', () => {
    expect(analytics.hasChosen).toBe(false);

    analytics.consents = [];
    expect(analytics.hasChosen).toBe(true);

    analytics.consents = [AnalyticsConsent.TRACKING];
    expect(analytics.hasChosen).toBe(true);
  });

  it('includes DNT flag when determining if user has made a choice', () => {
    jest.spyOn(config, 'ANALYTICS_IS_DNT').mockReturnValueOnce(true);
    expect(analytics.hasChosen).toBe(true);
  });

  it('provides a mechanism to listen to consent changes', () => {
    let consents;
    analytics.onConsentsChange((newConsents) => (consents = newConsents));
    expect(consents).toBe(undefined);

    analytics.consents = [];
    expect(consents).toEqual([]);

    analytics.consents = [AnalyticsConsent.TRACKING];
    expect(consents).toEqual([AnalyticsConsent.TRACKING]);
  });

  it('gets the visitor ID and automatically generates and persists it if unset', () => {
    const MOCKED_VISITOR_ID = 'c15a8afe1b3446b89c831613e7acce4f';
    const generateVisitorIdSpy = jest.spyOn(analyticsUtils, 'generateVisitorId');
    generateVisitorIdSpy.mockReturnValueOnce(MOCKED_VISITOR_ID);

    expect(analytics.visitorId).toBe(MOCKED_VISITOR_ID);

    generateVisitorIdSpy.mockReturnValueOnce('aaaaa123456789bbbbbb123456789ccc');
    expect(analytics.visitorId).toBe(MOCKED_VISITOR_ID);
  });

  it('resets the store', () => {
    analytics.consents = [AnalyticsConsent.TRACKING];
    expect(analytics.consents).toEqual([AnalyticsConsent.TRACKING]);

    analytics.reset();
    expect(analytics.consents).toEqual([]);
  });
});

jest.mock('stencil-router-v2');

import * as config from 'config';
import stores from 'stores';
import { AnalyticsConsent } from 'stores/analytics';
import analyticsService from './analytics';
import * as fetchClient from 'utils/fetch-client';
import { AnalyticsServiceProvider } from 'config/analytics';

const requestSpy = jest.spyOn(fetchClient, 'request');
requestSpy.mockImplementation(async () => [null, new Headers()] as [any, Headers]);

const expectInRequestUrl = (fragment: string) =>
  expect(requestSpy).toHaveBeenCalledWith(expect.objectContaining({ url: expect.stringContaining(fragment) }));

describe('analytics service', () => {
  beforeAll(() => {
    global.screen = { ...global.screen, width: 720, height: 400 };
    jest.spyOn(stores.auth, 'isAuthenticated', 'get').mockReturnValue(true);
    config.ANALYTICS_CONFIG.PROVIDER = AnalyticsServiceProvider.MATOMO;
  });

  describe('canTrack()', () => {
    it('returns a boolean', () => expect(typeof analyticsService.canTrack()).toBe('boolean'));

    it('returns false if tracking is disabled for current environment (no site ID or DNT flag is set)', async () => {
      jest.spyOn(config, 'ANALYTICS_IS_ENABLED').mockReturnValueOnce(false);
      stores.analytics.consents = [AnalyticsConsent.TRACKING];
      expect(await analyticsService.canTrack()).toBe(false);
    });

    it('returns false if tracking consent is not granted', () => {
      jest.spyOn(config, 'ANALYTICS_IS_ENABLED').mockReturnValueOnce(true);
      stores.analytics.consents = [];
      expect(analyticsService.canTrack()).toBe(false);
    });

    it('returns true if tracking is enabled for current environment and tracking consent is granted', () => {
      jest.spyOn(config, 'ANALYTICS_IS_ENABLED').mockReturnValue(true);
      stores.analytics.consents = [AnalyticsConsent.TRACKING];
      expect(analyticsService.canTrack()).toBe(true);
    });
  });

  describe('impression getter', () => {
    it('returns an AnalyticsImpression object with current page title, url and referrer', () => {
      const titleSpy = jest.spyOn(document, 'title', 'get');
      const hrefSpy = jest.spyOn(document.location, 'href', 'get');

      titleSpy.mockReturnValue('Foo Bar');
      hrefSpy.mockReturnValue('https://fo.o/bar');
      Object.defineProperty(document, 'referrer', { value: 'https://fo.o/baz' });
      expect(analyticsService.impression).toEqual({
        title: 'Foo Bar',
        url: 'https://fo.o/bar',
        referrer: 'https://fo.o/baz',
      });

      titleSpy.mockReturnValue('Foo Baz');
      hrefSpy.mockReturnValue('https://fo.o/baz');
      Object.defineProperty(document, 'referrer', { value: 'https://fo.o/bar' });
      expect(analyticsService.impression).toEqual({
        title: 'Foo Baz',
        url: 'https://fo.o/baz',
        referrer: 'https://fo.o/bar',
      });
    });
  });

  describe('trackImpression', () => {
    beforeEach(() => {
      jest.spyOn(config, 'ANALYTICS_IS_ENABLED').mockReturnValueOnce(true);
      stores.analytics.consents = [AnalyticsConsent.TRACKING];
    });

    it('performs a request against the API of the tracking provider including the AnalyticsImpression', async () => {
      jest.spyOn(document, 'title', 'get').mockReturnValue('Test page');
      jest.spyOn(document.location, 'href', 'get').mockReturnValue('https://foo.bar/baz');
      Object.defineProperty(document, 'referrer', { value: 'https://bar.baz/foo' });

      jest.clearAllMocks();
      await analyticsService.trackImpression();
      expectInRequestUrl('Test+page');
      expectInRequestUrl('https%3A%2F%2Ffoo.bar%2Fbaz');
      expectInRequestUrl('https%3A%2F%2Fbar.baz%2Ffoo');
    });

    it('schedules heart beat tracking', async () => {
      const scheduleHeartbeatTrackingSpy = jest.spyOn(analyticsService, 'scheduleHeartbeatTracking');

      jest.clearAllMocks();
      await analyticsService.trackImpression();
      expect(scheduleHeartbeatTrackingSpy).toHaveBeenCalled();
    });
  });

  describe('trackEvent', () => {
    beforeEach(() => {
      jest.spyOn(config, 'ANALYTICS_IS_ENABLED').mockReturnValueOnce(true);
      stores.analytics.consents = [AnalyticsConsent.TRACKING];
    });

    it('performs a request against the API of the tracking provider including the AnalyticsEvent and url', async () => {
      jest.spyOn(document.location, 'href', 'get').mockReturnValue('https://foo.bar/baz');

      jest.clearAllMocks();
      await analyticsService.trackEvent('Foo category', 'Bar action');
      expectInRequestUrl('Foo+category');
      expectInRequestUrl('Bar+action');
      expectInRequestUrl('https%3A%2F%2Ffoo.bar%2Fbaz');

      jest.clearAllMocks();
      await analyticsService.trackEvent('Baz category', 'Foo action', 'Bar name', 9876543);
      expectInRequestUrl('Baz+category');
      expectInRequestUrl('Foo+action');
      expectInRequestUrl('Bar+name');
      expectInRequestUrl('9876543');
    });
  });

  describe('trackHeartbeat', () => {
    beforeEach(() => {
      jest.spyOn(config, 'ANALYTICS_IS_ENABLED').mockReturnValueOnce(true);
      stores.analytics.consents = [AnalyticsConsent.TRACKING];
    });

    it('performs a request against the API of the tracking provider including a ping parameter and the url', async () => {
      jest.spyOn(document.location, 'href', 'get').mockReturnValue('https://baz.foo/bar');

      jest.clearAllMocks();
      await analyticsService.trackHeartbeat();
      expectInRequestUrl('ping=');
      expectInRequestUrl('https%3A%2F%2Fbaz.foo%2Fbar');
    });
  });

  describe('scheduleHeartbeatTracking', () => {
    it('schedules any previous heartbeat tracking interval', async () => {
      const clearIntervalSpy = jest.spyOn(window, 'clearInterval');

      await analyticsService.scheduleHeartbeatTracking();
      expect(clearIntervalSpy).toHaveBeenCalled();
    });

    it('sets a new heartbeat tracking interval according to const (if tracking is enabled)', async () => {
      const setIntervalSpy = jest.spyOn(window, 'setInterval');
      const heartbeatTrackingIntervalConstSpy = jest
        .spyOn(config, 'ANALYTICS_HEARTBEAT_TRACKING_INTERVAL', 'get')
        .mockReturnValue(321);

      await analyticsService.scheduleHeartbeatTracking();
      expect(setIntervalSpy).toHaveBeenCalledWith(expect.anything(), 321);
      expect(heartbeatTrackingIntervalConstSpy).toHaveBeenCalled();
    });
  });
});

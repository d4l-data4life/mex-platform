jest.mock('stencil-router-v2');

import navigationService from './navigation';
import * as fetchClient from 'utils/fetch-client';
import * as config from 'config';

const { CONTENT_URL, NAVIGATION_CONFIG } = config;

const createLink = (label: string, url: string, testAttr?: string) => {
  return { label, url, testAttr };
};

const EXAMPLE_CONTENT = {
  [config.LANGUAGE_CODES[0]]: {
    header: [createLink('HomePage', 'https://example.com/', 'test:home')],
    footer: [createLink('Impressum', 'https://example.com/', 'test:home')],
    learningEnvironment: [createLink('Get Started', 'https://example.com/', 'test:home')],
    services: [createLink('More', 'https://example.com/', 'test:home')],
  },
};

const requestSpy = jest.spyOn(fetchClient, 'request');

describe('navigation service', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('fetch()', () => {
    it('fetches navigation links', async () => {
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_CONTENT, new Headers()]);

      expect(await navigationService.fetch()).toBe(EXAMPLE_CONTENT);
      expect(requestSpy).toHaveBeenCalledWith(
        expect.objectContaining({
          method: 'GET',
          url: CONTENT_URL,
        })
      );
    });

    it('throws an error when the request fails', async () => {
      const consoleSpy = jest.spyOn(console, 'error');
      consoleSpy.mockImplementationOnce(() => {});
      requestSpy.mockImplementationOnce(() => {
        throw new Error('Network error');
      });

      await expect(navigationService.fetch()).rejects.toThrow();
    });
  });

  describe('load()', () => {
    it('loads navigation links', async () => {
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_CONTENT, new Headers()]);

      await navigationService.load();
      const language = config.LANGUAGE_CODES[0];
      expect(NAVIGATION_CONFIG[language].HEADER).toBe(EXAMPLE_CONTENT[language].header);
      expect(NAVIGATION_CONFIG[language].FOOTER).toBe(EXAMPLE_CONTENT[language].footer);
      expect(NAVIGATION_CONFIG[language].LEARNING_ENVIRONMENT).toBe(EXAMPLE_CONTENT[language].learningEnvironment);
      expect(NAVIGATION_CONFIG[language].SERVICES).toBe(EXAMPLE_CONTENT[language].services);
    });
  });
});

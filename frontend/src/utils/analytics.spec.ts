jest.mock('stencil-router-v2');

import { generateNonce, generateVisitorId, redactUrl } from './analytics';

describe('analytics util', () => {
  beforeEach(() => {
    window.btoa = jest.fn((str) => Buffer.from(str.toString(), 'binary').toString('base64'));
  });

  describe('generateNonce()', () => {
    it('returns a 36-character string', () => {
      expect(typeof generateNonce()).toBe('string');
      expect(generateNonce().length).toBe(36);
    });
  });

  describe('generateVisitorId()', () => {
    it('returns a 32-character hexadecimal string', () => {
      expect(typeof generateVisitorId()).toBe('string');
      expect(/^[0-9a-f]{32}$/.test(generateVisitorId())).toBe(true);
    });
  });

  describe('redactUrl', () => {
    it('redacts the query from search urls', () => {
      expect(redactUrl('/foo/bar')).toBe('/foo/bar');
      expect('/search?foo=bar').toBe('/search?foo=bar');
      expect(redactUrl('https://mex.app/search/foo+bar')).toBe('https://mex.app/search/redacted');
      expect(redactUrl('/search/highly+sensitive+stuff?foo=bar')).toBe('/search/redacted?foo=bar');
    });
  });
});

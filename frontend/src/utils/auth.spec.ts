jest.mock('stencil-router-v2');

import { webcrypto } from 'crypto';
import * as config from 'config';
import { generateChallenge, generateState, parseJwtPayload } from './auth';

const overwriteCrypto = (newCrypto) => {
  if (newCrypto) {
    (window.crypto as any) = global.crypto = newCrypto as any;
  } else {
    delete (window as any).crypto;
  }
};

describe('auth util', () => {
  describe('generateChallenge()', () => {
    beforeEach(() => {
      window.btoa = jest.fn((str) => Buffer.from(str.toString(), 'binary').toString('base64'));
      overwriteCrypto(webcrypto);
    });

    it('requires web crypto API', async () => {
      overwriteCrypto(null);
      expect(await generateChallenge()).toEqual([null, null]);

      overwriteCrypto({});
      expect(await generateChallenge()).toEqual([null, null]);

      overwriteCrypto(webcrypto);
      expect(await generateChallenge()).not.toEqual([null, null]);
    });

    it('corresponds to the AUTH_CHALLENGE_VERIFIER_LENGTH and AUTH_CHALLENGE_VERIFIER_MASK consts', async () => {
      const verifierLengthSpy = jest.spyOn(config, 'AUTH_CHALLENGE_VERIFIER_LENGTH', 'get');
      const verifierMaskSpy = jest.spyOn(config, 'AUTH_CHALLENGE_VERIFIER_MASK', 'get');

      verifierLengthSpy.mockReturnValue(17);
      verifierMaskSpy.mockReturnValue('abc');
      expect((await generateChallenge())[0].match(/[abc]/g).length).toBe(17);

      verifierLengthSpy.mockReturnValue(32);
      verifierMaskSpy.mockReturnValue('defg');
      expect((await generateChallenge())[0].match(/[defg]/g).length).toBe(32);
    });

    it('generates a verifier and a challenge', async () => {
      const verifierLength = 32;
      jest.spyOn(config, 'AUTH_CHALLENGE_VERIFIER_LENGTH', 'get').mockReturnValue(verifierLength);
      jest.spyOn(config, 'AUTH_CHALLENGE_VERIFIER_MASK', 'get').mockReturnValue('abcdefghij');

      const { subtle } = crypto;
      overwriteCrypto({
        getRandomValues: jest.fn((arr) => arr.map((_, index) => (index / verifierLength) * 256)),
        subtle,
      });

      expect(await generateChallenge()).toEqual([
        'aaaabbbcccdddeeeffffggghhhiiijjj',
        '38JsTWCO0cR4E8Wd5jg7ygn69lg20PzQuIYP4DJOtTY',
      ]);
    });
  });

  describe('generateState()', () => {
    it("makes use of web crypto API's randomUUID() if present", () => {
      const randomUUIDSpy = jest.spyOn(webcrypto as any, 'randomUUID');

      overwriteCrypto(null);
      generateState();
      expect(randomUUIDSpy).not.toHaveBeenCalled();

      overwriteCrypto({});
      generateState();
      expect(randomUUIDSpy).not.toHaveBeenCalled();

      overwriteCrypto(webcrypto);
      generateState();
      expect(randomUUIDSpy).toHaveBeenCalled();
    });

    it('returns a 36-character string', () => {
      overwriteCrypto(null);
      expect(typeof generateState()).toBe('string');
      expect(generateState().length).toBe(36);

      overwriteCrypto(webcrypto);
      expect(typeof generateState()).toBe('string');
      expect(generateState().length).toBe(36);
    });
  });

  describe('parseJwtPayload()', () => {
    it('parses the JWT payload and returns an object', () => {
      window.atob = (str: string) => Buffer.from(str, 'base64').toString('binary');
      expect(parseJwtPayload('foo.eyJlbWFpbCI6ImpvaG5AZG8uZSJ9.bar')).toEqual({ email: 'john@do.e' });
    });
  });
});

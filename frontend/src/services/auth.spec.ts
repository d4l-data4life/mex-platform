jest.mock('stencil-router-v2');

import { webcrypto } from 'crypto';
import * as config from 'config';
import * as authUtils from 'utils/auth';
import { AuthServiceProvider } from 'config/auth';
import authService from './auth';
import stores from 'stores';
import * as fetchClient from 'utils/fetch-client';

const requestSpy = jest.spyOn(fetchClient, 'request');
requestSpy.mockImplementation(async () => [EXAMPLE_AUTH_RESPONSE, new Headers()] as [any, Headers]);

const EXAMPLE_AUTH_RESPONSE = {
  access_token: 'access-token-foo',
  refresh_token: 'refresh-token-bar',
} as any;

describe('auth service', () => {
  let authProviderSpy, authAuthorizeUriSpy;
  beforeAll(() => {
    authProviderSpy = jest.spyOn(config, 'AUTH_PROVIDER').mockImplementation(async () => AuthServiceProvider.MICROSOFT);
    authAuthorizeUriSpy = jest.spyOn(config, 'AUTH_AUTHORIZE_URI').mockImplementation(async () => '/auth/authorize');
    jest.spyOn(config, 'AUTH_TOKEN_URI').mockImplementation(async () => '/auth/token');
    jest.spyOn(config, 'AUTH_LOGOUT_URI').mockImplementation(async () => '/auth/logout');
  });

  beforeEach(() => {
    jest.clearAllMocks();
    window.btoa = jest.fn((str) => Buffer.from(str.toString(), 'binary').toString('base64'));
    (window.crypto as any) = global.crypto = webcrypto as any;
  });

  describe('generateAuthorizeUrl()', () => {
    it('returns a url containing the string from the AUTH_AUTHORIZE_URI const', async () => {
      expect(authAuthorizeUriSpy).not.toHaveBeenCalled();
      expect(await authService.generateAuthorizeUrl()).toContain('/auth/authorize');
      expect(authAuthorizeUriSpy).toHaveBeenCalled();
    });

    it('throws an error when AUTH_AUTHORIZE_URI is not present', async () => {
      authAuthorizeUriSpy.mockReturnValueOnce(undefined);
      await expect(authService.generateAuthorizeUrl()).rejects.toThrow();
    });

    it('generates a state and challenge', async () => {
      const generateChallengeSpy = jest.spyOn(authUtils, 'generateChallenge');
      const generateStateSpy = jest.spyOn(authUtils, 'generateState');

      await authService.generateAuthorizeUrl();
      expect(generateChallengeSpy).toHaveBeenCalled();
      expect(generateStateSpy).toHaveBeenCalled();
    });

    it('stores the verifier and state', async () => {
      const authStoreVerifierSetterSpy = jest.spyOn(stores.auth, 'verifier', 'set');
      const authStoreStateSetterSpy = jest.spyOn(stores.auth, 'state', 'set');

      await authService.generateAuthorizeUrl();
      expect(authStoreVerifierSetterSpy).toHaveBeenCalledWith(expect.stringContaining(''));
      expect(authStoreStateSetterSpy).toHaveBeenCalledWith(expect.stringContaining('-'));
    });

    it('generates a url string containing AUTH_CLIENT_ID', async () => {
      const authClientIdSpy = jest.spyOn(config, 'AUTH_CLIENT_ID');

      authClientIdSpy.mockImplementation(async () => 'client-id-bar');

      const url = await authService.generateAuthorizeUrl();
      expect(typeof url).toBe('string');
      expect(url).toContain('client-id-bar');
    });
  });

  describe('redeemCode()', () => {
    it('switches redeem method according to AUTH_PROVIDER const', async () => {
      expect(authProviderSpy).not.toHaveBeenCalled();
      expect(await authService.redeemCode('foo')).toBe(undefined);

      authProviderSpy.mockReturnValueOnce('exoticFooProvider');
      await expect(authService.redeemCode('foo')).rejects.toThrow();
    });

    it('requests the token endpoint handing over the code and verifier', async () => {
      const authStoreVerifierGetterSpy = jest.spyOn(stores.auth, 'verifier', 'get');
      authStoreVerifierGetterSpy.mockReturnValueOnce('my-foo-verifier');

      await authService.redeemCode('my-foo-code');

      expect(requestSpy).toHaveBeenCalledWith(
        expect.objectContaining({ body: expect.stringContaining('my-foo-code') })
      );
      expect(requestSpy).toHaveBeenCalledWith(
        expect.objectContaining({ body: expect.stringContaining('my-foo-verifier') })
      );
    });

    it('resets the auth store and then stores access token and refresh token', async () => {
      const authStoreResetSpy = jest.spyOn(stores.auth, 'reset');
      const authStoreAccessTokenSetterSpy = jest.spyOn(stores.auth, 'accessToken', 'set');
      const authStoreRefreshTokenSetterSpy = jest.spyOn(stores.auth, 'refreshToken', 'set');

      expect(authStoreResetSpy).not.toHaveBeenCalled();
      expect(authStoreAccessTokenSetterSpy).not.toHaveBeenCalled();
      expect(authStoreRefreshTokenSetterSpy).not.toHaveBeenCalled();

      await authService.redeemCode('my-foo-code');

      expect(authStoreResetSpy).toHaveBeenCalled();
      expect(authStoreAccessTokenSetterSpy).toHaveBeenCalledWith('access-token-foo');
      expect(authStoreRefreshTokenSetterSpy).toHaveBeenCalledWith('refresh-token-bar');
    });
  });

  describe('refreshSession()', () => {
    it('switches refresh method according to AUTH_PROVIDER const', async () => {
      expect(authProviderSpy).not.toHaveBeenCalled();
      await expect(authService.refreshSession()).resolves.not.toThrow();

      authProviderSpy.mockReturnValueOnce('exoticFooProvider');
      await expect(authService.refreshSession()).rejects.toThrow();
    });

    it('requires a refresh token to be present in the auth store', async () => {
      const authStoreRefreshTokenGetterSpy = jest.spyOn(stores.auth, 'refreshToken', 'get');
      authStoreRefreshTokenGetterSpy.mockReturnValue(null);

      await expect(authService.refreshSession()).rejects.toThrow();

      authStoreRefreshTokenGetterSpy.mockReturnValue('refresh-foo');
      await expect(authService.refreshSession()).resolves.not.toThrow();
    });

    it('requests the token endpoint handing over the refresh token', async () => {
      const authStoreRefreshTokenGetterSpy = jest.spyOn(stores.auth, 'refreshToken', 'get');
      authStoreRefreshTokenGetterSpy.mockReturnValue('refresh-token-bar');

      await authService.refreshSession();

      expect(requestSpy).toHaveBeenCalledWith(
        expect.objectContaining({ body: expect.stringContaining('refresh-token-bar') })
      );
    });
  });

  describe('expireSession()', () => {
    it('resets the auth store', () => {
      const authStoreResetSpy = jest.spyOn(stores.auth, 'reset');
      expect(authStoreResetSpy).not.toHaveBeenCalled();
      authService.expireSession();
      expect(authStoreResetSpy).toHaveBeenCalled();
    });

    it('displays a notification about the session having expired', () => {
      const notifcationsStoreAddSpy = jest.spyOn(stores.notifications, 'add');
      expect(notifcationsStoreAddSpy).not.toHaveBeenCalled();
      authService.expireSession();
      expect(notifcationsStoreAddSpy).toHaveBeenCalledWith('sessionExpired');
    });
  });

  describe('generateLogoutUrl()', () => {
    it('switches generation method according to AUTH_PROVIDER const', async () => {
      expect(authProviderSpy).not.toHaveBeenCalled();
      expect(typeof (await authService.generateLogoutUrl())).toBe('string');

      authProviderSpy.mockReturnValueOnce('exoticFooProvider');
      await expect(authService.generateLogoutUrl()).rejects.toThrow();
    });

    it('generates a url string containing AUTH_CLIENT_ID', async () => {
      const authClientIdSpy = jest.spyOn(config, 'AUTH_CLIENT_ID');

      authClientIdSpy.mockImplementation(async () => 'client-id-bar');

      const url = await authService.generateLogoutUrl();
      expect(url).toContain('/auth/logout');
      expect(url).toContain('client-id-bar');
    });
  });
});

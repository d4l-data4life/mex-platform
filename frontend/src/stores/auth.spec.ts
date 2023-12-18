jest.mock('stencil-router-v2');

import * as config from 'config';
import auth from './auth';
import * as authUtil from 'utils/auth';

describe('auth store', () => {
  let sessionStorageSetItemSpy, localStorageSetItemSpy;

  beforeAll(() => {
    auth.reset();

    sessionStorageSetItemSpy = jest.spyOn(window.sessionStorage, 'setItem');
    sessionStorageSetItemSpy.mockImplementation(() => {});
    localStorageSetItemSpy = jest.spyOn(window.localStorage, 'setItem');
    localStorageSetItemSpy.mockImplementation(() => {});
  });

  beforeEach(() => {
    jest.clearAllMocks();
    auth.reset();
  });

  it('gets and sets the verifier', () => {
    expect(auth.verifier).toBe(undefined);
    auth.verifier = 'verifier-foo';
    expect(auth.verifier).toBe('verifier-foo');
  });

  it('gets and sets the state', () => {
    expect(auth.state).toBe(undefined);
    auth.state = 'state-foo';
    expect(auth.state).toBe('state-foo');
  });

  it('gets and sets the access token', () => {
    expect(auth.accessToken).toBe(undefined);
    auth.accessToken = 'accessToken-foo';
    expect(auth.accessToken).toBe('accessToken-foo');
  });

  it('gets and sets the refresh token', () => {
    expect(auth.refreshToken).toBe(undefined);
    auth.refreshToken = 'refreshToken-foo';
    expect(auth.refreshToken).toBe('refreshToken-foo');
  });

  it('gets and sets the requested route', () => {
    expect(auth.requestedRoute).toBe(undefined);

    auth.requestedRoute = 'route-foo';
    expect(auth.requestedRoute).toBe('route-foo');

    auth.requestedRoute = 'https://evil.corp';
    expect(auth.requestedRoute).toBe('https/evil.corp');
  });

  it('resets the session store', () => {
    auth.accessToken = 'accessToken-foo';
    expect(auth.accessToken).toBe('accessToken-foo');

    auth.resetSession();

    expect(auth.accessToken).toBe(undefined);
  });

  it('resets the store(s)', () => {
    auth.verifier = 'verifier-foo';
    auth.state = 'state-foo';
    auth.accessToken = 'accessToken-foo';
    auth.isReturning = true;

    expect(auth.verifier).toBe('verifier-foo');
    expect(auth.state).toBe('state-foo');
    expect(auth.accessToken).toBe('accessToken-foo');
    expect(auth.isReturning).toBe(true);

    auth.reset();

    expect(auth.verifier).toBe(undefined);
    expect(auth.state).toBe(undefined);
    expect(auth.accessToken).toBe(undefined);
    expect(auth.isReturning).toBe(false);
  });

  it('has an isAuthenticated getter (bool) to return if there is an access token', () => {
    expect(auth.isAuthenticated).toBe(false);
    auth.accessToken = 'foo-bar-baz';
    expect(auth.isAuthenticated).toBe(true);
  });

  it('persists verifier and state', () => {
    expect(sessionStorageSetItemSpy).not.toHaveBeenCalled();

    auth.verifier = 'verifier-foo';
    auth.state = 'state-bar';

    expect(sessionStorageSetItemSpy).toHaveBeenCalledTimes(2);
    expect(sessionStorageSetItemSpy).toHaveBeenCalledWith(expect.stringContaining('verifier'), '"verifier-foo"');
    expect(sessionStorageSetItemSpy).toHaveBeenCalledWith(expect.stringContaining('state'), '"state-bar"');
  });

  it('persists access token and refresh token if feature flag is set', () => {
    expect(sessionStorageSetItemSpy).not.toHaveBeenCalled();
    const persistTokensFeatureFlagSpy = jest.spyOn(config, 'AUTH_PERSIST_TOKENS', 'get');

    persistTokensFeatureFlagSpy.mockReturnValue(false);
    jest.resetModules();

    auth.accessToken = 'accessToken-foo';
    auth.refreshToken = 'refreshToken-bar';

    expect(sessionStorageSetItemSpy).not.toHaveBeenCalled();

    persistTokensFeatureFlagSpy.mockReturnValue(true);
    jest.resetModules();

    auth.accessToken = 'accessToken-foo';
    auth.refreshToken = 'refreshToken-bar';

    expect(sessionStorageSetItemSpy).toHaveBeenCalledTimes(2);
    expect(sessionStorageSetItemSpy).toHaveBeenCalledWith(expect.stringContaining('accessToken'), '"accessToken-foo"');
    expect(sessionStorageSetItemSpy).toHaveBeenCalledWith(
      expect.stringContaining('refreshToken'),
      '"refreshToken-bar"'
    );
  });

  it('persists flag about returning visitor', () => {
    expect(localStorageSetItemSpy).not.toHaveBeenCalled();
    auth.isReturning = true;
    expect(localStorageSetItemSpy).toHaveBeenCalledWith(expect.stringContaining('isReturning'), 'true');
  });

  it('parses the access token JWT payload and returns the email adress', () => {
    jest.spyOn(authUtil, 'parseJwtPayload').mockImplementation(() => ({ email: 'jane@do.e' }));
    expect(auth.userEmail).toBe('jane@do.e');
  });
});

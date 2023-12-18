jest.mock('stencil-router-v2');

import { ResponseError } from 'models/response-error';
import * as config from 'config';
import stores from 'stores';
import services from 'services';
import { del, get, post, put, request } from './fetch-client';
import { AuthServiceProvider } from 'config/auth';

const DUMMY_URL = 'https://foo.bar.dev';

const FETCH_SUCCESS = {
  ok: true,
  status: 200,
  headers: new Headers({ 'Content-Type': 'application/json' }),
};

const FETCH_ARGUMENTS = {
  credentials: 'omit',
  headers: { 'Content-Type': 'application/json', 'X-User-Language': expect.anything() },
  method: 'GET',
};

const expectValue = (key: string, value: any) => expect.objectContaining({ [key]: value });

const expectNestedValue = (key: string, name: string, value: any) =>
  expect.objectContaining({ [key]: expect.objectContaining({ [name]: value }) });

describe('fetch client', () => {
  let fetchSpy;
  let jsonMock;
  let textMock;

  beforeAll(() => {
    fetchSpy = jest.spyOn(global, 'fetch');
  });

  beforeEach(() => {
    jest.resetAllMocks();
    jsonMock = jest.fn(() => ({ success: true }));
    textMock = jest.fn(() => 'success');
    fetchSpy.mockImplementation(() => ({ ...FETCH_SUCCESS, json: jsonMock, text: textMock }));
  });

  describe('request()', () => {
    it('uses the fetch api', async () => {
      await request({ url: DUMMY_URL });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, FETCH_ARGUMENTS);
    });

    it('sends Authorization header if argument is true', async () => {
      stores.auth.accessToken = 'foo-auth';
      await request({ url: DUMMY_URL, authorized: true });
      expect(fetchSpy).toHaveBeenCalledWith(
        DUMMY_URL,
        expectNestedValue('headers', 'Authorization', 'Bearer foo-auth')
      );

      stores.auth.reset();
      await request({ url: DUMMY_URL, authorized: true });
      expect(fetchSpy).toHaveBeenCalledWith(
        DUMMY_URL,
        expectNestedValue('headers', 'Authorization', 'Bearer undefined')
      );
    });

    it('allows providing custom headers', async () => {
      await request({ url: DUMMY_URL, headers: { 'X-Unit-Test': 'foo' } });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expectNestedValue('headers', 'X-Unit-Test', 'foo'));
      expect(fetchSpy).toHaveBeenCalledWith(
        DUMMY_URL,
        expectNestedValue('headers', 'Content-Type', 'application/json')
      );
    });

    it('allows specifying if user language shall be sent (default: true)', async () => {
      const languageSpy = jest.spyOn(stores.i18n, 'language', 'get');

      languageSpy.mockReturnValueOnce('es');
      await request({ url: DUMMY_URL });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expectNestedValue('headers', 'X-User-Language', 'es'));

      jest.clearAllMocks();
      languageSpy.mockReturnValueOnce('fr');
      await request({ url: DUMMY_URL, sendLanguage: true });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expectNestedValue('headers', 'X-User-Language', 'fr'));

      jest.clearAllMocks();
      languageSpy.mockReturnValueOnce('no');
      await request({ url: DUMMY_URL, sendLanguage: false });
      expect(fetchSpy).not.toHaveBeenCalledWith(DUMMY_URL, expectNestedValue('headers', 'X-User-Language', 'no'));
    });

    it('allows specifying the request method', async () => {
      await request({ url: DUMMY_URL, method: 'POST' });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expectValue('method', 'POST'));
    });

    it('allows providing a request body', async () => {
      const body = JSON.stringify({ foo: 'bar' });
      await request({ url: DUMMY_URL, body });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expectValue('body', body));
    });

    it('throws a response error with status code 0 when the request itself repeatedly fails (e.g. network offline, timeout, cors issues...)', async () => {
      const consoleSpy = jest.spyOn(console, 'error');
      const errorMessage = 'NetworkError when attempting to fetch resource.';

      new Array(2).fill(null).forEach(() =>
        fetchSpy.mockImplementationOnce(() => {
          throw new Error(errorMessage);
        })
      );
      consoleSpy.mockImplementationOnce(() => {});

      expect.assertions(4);

      try {
        await request({ url: DUMMY_URL });
      } catch (error) {
        expect(error).toBeInstanceOf(ResponseError);
        expect(error.message).toBe(errorMessage);
        expect(error.status).toBe(0);
        expect(consoleSpy).toHaveBeenCalledWith('network error', error);
      }
    });

    it('does not throw a response error when only the first request fails (e.g. timeout) but the repeated request works', async () => {
      fetchSpy.mockImplementationOnce(() => {
        throw new Error('First request did not work');
      });

      expect.assertions(0);
      try {
        await request({ url: DUMMY_URL });
      } catch (_) {
        expect(true).toBe(true);
      }
    });

    it('allows to specify a timeout after which the first request attempt is cancelled', async () => {
      await request({ url: DUMMY_URL });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expect.not.objectContaining({ signal: expect.anything() }));

      await request({ url: DUMMY_URL, timeout: 35 });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expect.objectContaining({ signal: expect.anything() }));
    });

    it('silently refreshes the session when status code is 401 and the call was authorized', async () => {
      jest.spyOn(config, 'AUTH_PROVIDER').mockImplementation(async () => AuthServiceProvider.NOAUTH);
      jest.spyOn(config, 'AUTH_TOKEN_URI').mockImplementation(async () => '/auth/token');

      const refreshSessionSpy = jest.spyOn(services.auth, 'refreshSession');
      fetchSpy.mockImplementationOnce(() => ({
        ok: false,
        status: 401,
        headers: new Headers({ 'Content-Type': 'application/json' }),
        text: textMock,
      }));

      textMock.mockReturnValueOnce('Unauthorized');
      jsonMock.mockReturnValueOnce({
        access_token: 'access-token-refreshed',
        refresh_token: 'refresh-token-refreshed',
      });

      stores.auth.refreshToken = 'refresh-token-old';

      expect(refreshSessionSpy).not.toHaveBeenCalled();
      expect(await request({ url: DUMMY_URL, authorized: true })).toEqual([{ success: true }, FETCH_SUCCESS.headers]);

      expect(refreshSessionSpy).toHaveBeenCalled();
      expect(fetchSpy).toHaveBeenCalledWith(
        expect.anything(),
        expectValue('body', expect.stringContaining('refresh-token-old'))
      );
      expect(fetchSpy).toHaveBeenCalledWith(
        DUMMY_URL,
        expectNestedValue('headers', 'Authorization', 'Bearer access-token-refreshed')
      );
      expect(stores.auth.accessToken).toBe('access-token-refreshed');
      expect(stores.auth.refreshToken).toBe('refresh-token-refreshed');
    });

    it('expires the session when silent refresh of the session fails', async () => {
      const expireSessionSpy = jest.spyOn(services.auth, 'expireSession');
      fetchSpy.mockImplementation(() => ({
        ok: false,
        status: 401,
        headers: new Headers({ 'Content-Type': 'application/json' }),
        text: textMock,
      }));

      stores.auth.refreshToken = 'refresh-token-old';
      textMock.mockReturnValue('Unauthorized');

      expect(expireSessionSpy).not.toHaveBeenCalled();
      await expect(request({ url: DUMMY_URL, authorized: true })).rejects.toThrow();
      expect(expireSessionSpy).toHaveBeenCalled();
    });

    it('returns null when status code is 204', async () => {
      fetchSpy.mockImplementationOnce(() => ({ ...FETCH_SUCCESS, status: 204 }));
      expect(await request({ url: DUMMY_URL })).toEqual([null, FETCH_SUCCESS.headers]);
    });

    it('returns json when status code is in the 2xx range and json header is returned', async () => {
      expect(await request({ url: DUMMY_URL })).toEqual([{ success: true }, FETCH_SUCCESS.headers]);
      expect(jsonMock).toHaveBeenCalled();
      expect(textMock).not.toHaveBeenCalled();
    });

    it('returns string when status code is in the 2xx range and non-json header is returned', async () => {
      const headers = new Headers({ 'Content-Type': 'text' });
      fetchSpy.mockImplementation(() => ({ ...FETCH_SUCCESS, headers, json: jsonMock, text: textMock }));

      expect(await request({ url: DUMMY_URL })).toEqual(['success', headers]);
      expect(textMock).toHaveBeenCalled();
      expect(jsonMock).not.toHaveBeenCalled();
    });

    it('throws a response error when status code is not in the 2xx range', async () => {
      const errorMessage = 'Unauthorized';

      fetchSpy.mockImplementationOnce(() => ({
        ok: false,
        status: 401,
        text: jest.fn(() => errorMessage),
      }));

      expect.assertions(3);

      try {
        await request({ url: DUMMY_URL });
      } catch (error) {
        expect(error).toBeInstanceOf(ResponseError);
        expect(error.message).toBe(errorMessage);
        expect(error.status).toBe(401);
      }
    });
  });

  describe('get()', () => {
    it('is an alias for request() with method=GET', async () => {
      await get({ url: DUMMY_URL });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expectValue('method', 'GET'));
    });
  });

  describe('post()', () => {
    it('is an alias for request() with method=POST', async () => {
      await post({ url: DUMMY_URL });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expectValue('method', 'POST'));
    });
  });

  describe('put()', () => {
    it('is an alias for request() with method=PUT', async () => {
      await put({ url: DUMMY_URL });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expectValue('method', 'PUT'));
    });
  });

  describe('del()', () => {
    it('is an alias for request() with method=DELETE', async () => {
      await del({ url: DUMMY_URL });
      expect(fetchSpy).toHaveBeenCalledWith(DUMMY_URL, expectValue('method', 'DELETE'));
    });
  });
});

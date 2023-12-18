import { IS_PREVIEW_MODE } from 'config';
import { ResponseError } from 'models/response-error';
import authService from 'services/auth';
import stores from 'stores';

interface FetchClientHeaders {
  [key: string]: string;
}

interface FetchClientArguments {
  url: string;
  method?: string;
  contentType?: string;
  headers?: FetchClientHeaders;
  body?: any;
  credentials?: 'omit' | 'same-origin' | 'include';
  authorized?: boolean;
  sendLanguage?: boolean;
  isRetry?: boolean;
  timeout?: number;
}

export const request = async <T>(args: FetchClientArguments): Promise<[T, Headers]> => {
  const {
    url,
    method = 'GET',
    headers = {},
    body = null,
    contentType = 'application/json',
    credentials = 'omit',
    authorized = false,
    sendLanguage = true,
    isRetry = false,
    timeout = null,
  } = args;

  let response;
  let timeoutId;
  let timeoutController;

  if (timeout) {
    timeoutController = new AbortController();
    timeoutId = setTimeout(() => timeoutController.abort(), timeout);
  }

  // cache-busting for preview mode
  const previewAppendix = IS_PREVIEW_MODE ? `${url.includes('?') ? '&' : '?'}_t=${Date.now()}` : '';

  try {
    response = await fetch(url + previewAppendix, {
      method: method,
      headers: {
        ...(contentType ? { 'Content-Type': contentType } : {}),
        ...(authorized
          ? {
              Authorization: `Bearer ${stores.auth.accessToken}`,
            }
          : {}),
        ...(sendLanguage
          ? {
              'X-User-Language': stores.i18n.language,
            }
          : {}),
        ...headers,
      },
      ...(body ? { body } : {}),
      credentials,
      signal: timeoutController?.signal,
    });
  } catch (e) {
    if (!isRetry) {
      return request({ ...args, isRetry: true, timeout: null });
    }

    // most likely network or cors error
    console.error('network error', e);
    throw new ResponseError(e.message, 0);
  } finally {
    clearTimeout(timeoutId);
  }

  if (response.status === 401 && authorized && !isRetry) {
    try {
      await authService.refreshSession();
      return request({ ...args, isRetry: true });
    } catch (_) {
      authService.expireSession();
    }
  }

  if (response.status === 401 && isRetry) {
    authService.expireSession();
  }

  if (response.status === 204) {
    return [null, response.headers];
  }

  if (response.status >= 200 && response.status < 300) {
    return /application\/json/.test(response.headers.get('Content-Type') || '')
      ? [await response.json(), response.headers]
      : [await response.text(), response.headers];
  }

  throw new ResponseError(await response.text(), response.status);
};

export const get = async <T>(params: Omit<FetchClientArguments, 'method'>): Promise<[T, Headers]> =>
  request({ ...params, method: 'GET' });

export const post = async <T>(params: Omit<FetchClientArguments, 'method'>): Promise<[T, Headers]> =>
  request({ ...params, method: 'POST' });

export const put = async <T>(params: Omit<FetchClientArguments, 'method'>): Promise<[T, Headers]> =>
  request({ ...params, method: 'PUT' });

export const del = async <T>(params: Omit<FetchClientArguments, 'method'>): Promise<[T, Headers]> =>
  request({ ...params, method: 'DELETE' });

export default { request, get, post, put, del };

import { Env } from '@stencil/core';

const { DYNAMIC_ENV_URL } = Env;

export interface DynamicEnvVars {
  AUTH_PROVIDER?: string;
  AUTH_CLIENT_ID?: string;
  AUTH_AUTHORIZE_URI?: string;
  AUTH_TOKEN_URI?: string;
  AUTH_LOGOUT_URI?: string;
}

export default new Promise<DynamicEnvVars>(async (resolve) => {
  try {
    resolve(
      !DYNAMIC_ENV_URL ? {} : await (await fetch(DYNAMIC_ENV_URL, { method: 'GET', credentials: 'omit' })).json()
    );
  } catch {
    resolve({});
  }
});

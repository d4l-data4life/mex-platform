import {
  AUTH_PROVIDER,
  AUTH_CLIENT_ID,
  AUTH_SCOPE,
  AUTH_AUTHORIZE_URI,
  AUTH_TOKEN_URI,
  AUTH_LOGOUT_URI,
  ROUTES,
} from 'config';
import { AuthServiceProvider } from 'config/auth';
import stores from 'stores';
import { generateChallenge, generateState } from 'utils/auth';
import { post } from 'utils/fetch-client';

export interface AuthTokenResponse {
  access_token: string;
  expires_in: number;
  ext_expires_in: number;
  refresh_token: string;
  scope: string;
  token_type: 'Bearer';
}

type BodyStringifier = (data: Record<string, string>) => [string, string];

const stringifierFormUrlenconded: BodyStringifier = (d) => [
  new URLSearchParams(d).toString(),
  'application/x-www-form-urlencoded',
];
const stringifierJson: BodyStringifier = (d) => [JSON.stringify(d), 'application/json'];

export class AuthService {
  async generateAuthorizeUrl(): Promise<string> {
    const authorizeUri = await AUTH_AUTHORIZE_URI();
    if (!authorizeUri) {
      throw new Error('authorization URI not defined');
    }

    const state = generateState();
    const [verifier, challenge] = await generateChallenge();
    const hasChallenge = !!verifier && !!challenge;

    stores.auth.verifier = hasChallenge ? verifier : null;
    stores.auth.state = state;

    const params = new URLSearchParams({
      client_id: await AUTH_CLIENT_ID(),
      redirect_uri: document.location.origin + ROUTES.AUTH,
      response_type: 'code',
      response_mode: 'fragment',
      scope: await AUTH_SCOPE(),
      state,
      ...(hasChallenge
        ? {
            code_challenge: challenge,
            code_challenge_method: 'S256',
          }
        : {}),
    });

    return `${authorizeUri}?${params.toString()}`;
  }

  async redeemCode(code: string): Promise<void> {
    switch (await AUTH_PROVIDER()) {
      case AuthServiceProvider.MICROSOFT:
        return await this.redeemCodeWithBodyStringifier(code, stringifierFormUrlenconded);
      case AuthServiceProvider.NOAUTH:
        return await this.redeemCodeWithBodyStringifier(code, stringifierJson);
      default:
        throw new Error('Unknown or unset auth service provider');
    }
  }

  private async redeemCodeWithBodyStringifier(code: string, bodyStringifier: BodyStringifier): Promise<void> {
    const { verifier } = stores.auth;
    stores.auth.reset();

    const params = {
      client_id: await AUTH_CLIENT_ID(),
      redirect_uri: document.location.origin + ROUTES.AUTH,
      grant_type: 'authorization_code',
      code_verifier: verifier,
      code,
      scope: await AUTH_SCOPE(),
    };

    const [body, contentType] = bodyStringifier(params);

    const [response] = await post<AuthTokenResponse>({
      url: await AUTH_TOKEN_URI(),
      contentType,
      body,
    });

    stores.auth.accessToken = response.access_token;
    stores.auth.refreshToken = response.refresh_token;
  }

  async refreshSession(): Promise<void> {
    const { refreshToken } = stores.auth;
    if (!refreshToken) {
      throw new Error('session expired: refresh token missing');
    }

    switch (await AUTH_PROVIDER()) {
      case AuthServiceProvider.MICROSOFT:
        return await this.refreshSessionWithBodyStringifier(refreshToken, stringifierFormUrlenconded);
      case AuthServiceProvider.NOAUTH:
        return await this.refreshSessionWithBodyStringifier(refreshToken, stringifierJson);
      default:
        throw new Error('Unknown or unset auth service provider');
    }
  }

  private async refreshSessionWithBodyStringifier(
    refreshToken: string,
    bodyStringifier: BodyStringifier
  ): Promise<void> {
    const params = {
      client_id: await AUTH_CLIENT_ID(),
      grant_type: 'refresh_token',
      scope: await AUTH_SCOPE(),
      refresh_token: refreshToken,
    };

    const [body, contentType] = bodyStringifier(params);

    const [response] = await post<AuthTokenResponse>({
      url: await AUTH_TOKEN_URI(),
      contentType,
      body,
    });

    stores.auth.accessToken = response.access_token;
    stores.auth.refreshToken = response.refresh_token;
  }

  expireSession() {
    stores.auth.reset();
    stores.notifications.add('sessionExpired');
  }

  async generateLogoutUrl(): Promise<string> {
    switch (await AUTH_PROVIDER()) {
      case AuthServiceProvider.MICROSOFT:
        return this.generateMicrosoftLogoutUrl();
      default:
        throw new Error('Unknown or unset auth service provider or no logout url for given provider');
    }
  }

  private async generateMicrosoftLogoutUrl(): Promise<string> {
    const params = new URLSearchParams({
      client_id: await AUTH_CLIENT_ID(),
      post_logout_redirect_uri: document.location.origin + ROUTES.LOGOUT,
    });

    return `${await AUTH_LOGOUT_URI()}?${params.toString()}`;
  }
}

export default new AuthService();

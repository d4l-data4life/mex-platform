import { Env } from '@stencil/core';
import dynamicEnv from './dynamic';

export enum AuthServiceProvider {
  MICROSOFT = 'MICROSOFT',
  NOAUTH = 'NOAUTH',
}

export const AUTH_PROVIDER = async () => Env.AUTH_PROVIDER ?? (await dynamicEnv).AUTH_PROVIDER;
export const AUTH_CLIENT_ID = async () => Env.AUTH_CLIENT_ID ?? (await dynamicEnv).AUTH_CLIENT_ID;
export const AUTH_AUTHORIZE_URI = async () => Env.AUTH_AUTHORIZE_URI ?? (await dynamicEnv).AUTH_AUTHORIZE_URI;
export const AUTH_TOKEN_URI = async () => Env.AUTH_TOKEN_URI ?? (await dynamicEnv).AUTH_TOKEN_URI;
export const AUTH_LOGOUT_URI = async () => Env.AUTH_LOGOUT_URI ?? (await dynamicEnv).AUTH_LOGOUT_URI;

export const AUTH_SCOPE = async () => `api://${await AUTH_CLIENT_ID()}/metadata.read`;
export const AUTH_CHALLENGE_VERIFIER_LENGTH: number = 64;
export const AUTH_CHALLENGE_VERIFIER_MASK: string =
  '(){}+-~.;@0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
export const AUTH_PERSIST_TOKENS: boolean = true;

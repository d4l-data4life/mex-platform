import { createStore } from '@stencil/store';
import { AUTH_PERSIST_TOKENS } from 'config';
import { Persistor } from './utils/persistor';
import createPersistedStore from './utils/persisted-store';
import { parseJwtPayload } from 'utils/auth';

interface StateType {
  accessToken?: string;
  refreshToken?: string;
}

interface PersistedStateType extends StateType {
  verifier?: string;
  state?: string;
  requestedRoute?: string;
  isReturning?: boolean;
}

const sessionPersistor = new Persistor('sessionStorage' in window ? sessionStorage : null);
const localPersistor = new Persistor('localStorage' in window ? localStorage : null);
const sessionStore = createPersistedStore<PersistedStateType>(sessionPersistor, 'auth', {});
const localStore = createPersistedStore<PersistedStateType>(localPersistor, 'auth', {});
const store = createStore<StateType>({});

class AuthStore {
  get verifier() {
    return sessionStore.get('verifier');
  }

  set verifier(verifier: string) {
    sessionStore.set('verifier', verifier);
  }

  get state() {
    return sessionStore.get('state');
  }

  set state(state: string) {
    sessionStore.set('state', state);
  }

  get accessToken() {
    return AUTH_PERSIST_TOKENS ? sessionStore.get('accessToken') : store.get('accessToken');
  }

  set accessToken(accessToken: string) {
    AUTH_PERSIST_TOKENS ? sessionStore.set('accessToken', accessToken) : store.set('accessToken', accessToken);
  }

  get refreshToken() {
    return AUTH_PERSIST_TOKENS ? sessionStore.get('refreshToken') : store.get('refreshToken');
  }

  set refreshToken(refreshToken: string) {
    AUTH_PERSIST_TOKENS ? sessionStore.set('refreshToken', refreshToken) : store.set('refreshToken', refreshToken);
  }

  set requestedRoute(route: string) {
    sessionStore.set('requestedRoute', route);
  }

  get requestedRoute() {
    return sessionStore
      .get('requestedRoute')
      ?.replace(/:/g, '')
      .replace(/[/]{2,}/g, '/');
  }

  get isAuthenticated() {
    return !!this.accessToken;
  }

  get isReturning() {
    return localStore.get('isReturning') || false;
  }

  set isReturning(isReturning: boolean) {
    localStore.set('isReturning', isReturning);
  }

  get userEmail() {
    return parseJwtPayload(this.accessToken)?.email;
  }

  get userId() {
    return parseJwtPayload(this.accessToken).sub;
  }

  resetSession(): void {
    sessionStore.reset();
    store.reset();
  }

  reset(): void {
    sessionStore.reset();
    localStore.reset();
    store.reset();
  }
}

export default new AuthStore();

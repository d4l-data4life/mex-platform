import { Component, Fragment, h, Host, Prop, State } from '@stencil/core';
import { ROUTES } from 'config';
import services from 'services';
import stores from 'stores';

@Component({
  tag: 'mex-view-auth',
  styleUrl: './auth.css',
})
export class AuthView {
  @Prop() url: URL;

  @State() hasErrored = false;
  @State() isRedirecting = false;

  get isProgressBarShown() {
    return this.isRedirecting && !stores.auth.isReturning;
  }

  resetHash() {
    history.replaceState({}, '', this.url.pathname);
  }

  showError(error: string) {
    this.hasErrored = true;
    stores.notifications.add('oauthError');
    console.error('oauth error', error);
  }

  async redeemCode(code: string) {
    const { requestedRoute } = stores.auth;
    await services.auth.redeemCode(code);
    stores.auth.isReturning = true;
    stores.router.push(requestedRoute ?? ROUTES.ROOT);
  }

  async redirectToLogin(force = false, timeout = stores.auth.isReturning ? 0 : 2500) {
    if (this.hasErrored && !force) {
      return;
    }

    this.isRedirecting = true;

    setTimeout(async () => {
      document.location.href = await services.auth.generateAuthorizeUrl();
    }, timeout);
  }

  async performLogin() {
    const params = new URLSearchParams(this.url.hash.slice(1));
    const error = params.get('error');
    const state = params.get('state');
    const code = params.get('code');

    (code || state || error) && this.resetHash();

    const isStateMismatch = !!state && state !== stores.auth.state;
    if (isStateMismatch) {
      this.showError('state mismatch');
      return this.redirectToLogin();
    }

    error && this.showError(error);

    if (!error && code) {
      try {
        return await this.redeemCode(code);
      } catch (error) {
        this.showError(error.message);
      }
    }

    this.redirectToLogin();
  }

  componentWillLoad() {
    this.performLogin();
  }

  render() {
    const { isProgressBarShown, isRedirecting, hasErrored } = this;

    return (
      <Host class="auth view view--center">
        {isProgressBarShown && (
          <Fragment>
            <div class="auth__text">{stores.i18n.t('auth.redirect')}</div>
            <div class="auth__progress" />
          </Fragment>
        )}
        {!isProgressBarShown && !hasErrored && <mex-logo class="auth__loader" loader />}
        {!isRedirecting && hasErrored && (
          <Fragment>
            <h3 class="auth__error">{stores.i18n.t('auth.loginFailed')}</h3>
            <button class="button button--primary" onClick={() => this.redirectToLogin(true, 0)}>
              {stores.i18n.t('auth.retryLogin')}
            </button>
          </Fragment>
        )}
      </Host>
    );
  }
}

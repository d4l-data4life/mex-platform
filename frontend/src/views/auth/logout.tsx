import { Component, getAssetPath, h, Host, State } from '@stencil/core';
import { ROUTES } from 'config';
import services from 'services';
import { href } from 'stencil-router-v2';
import stores from 'stores';

@Component({
  tag: 'mex-view-logout',
  styleUrl: './auth.css',
  assetsDirs: ['assets'],
})
export class LogoutView {
  @State() isRedirecting = false;

  async logout() {
    this.isRedirecting = true;
    const logoutUrl = await services.auth.generateLogoutUrl();
    stores.auth.resetSession();
    requestAnimationFrame(() => (document.location.href = logoutUrl));
  }

  async componentWillLoad() {
    if (stores.auth.isAuthenticated) {
      await this.logout();
    }
  }

  render() {
    const { t } = stores.i18n;
    return (
      <Host class={`auth view ${this.isRedirecting ? '' : 'auth--darkTheme'}`} data-test="logout">
        {this.isRedirecting && <mex-logo class="auth__loader" loader />}
        {!this.isRedirecting && (
          <div>
            <img
              aria-hidden="true"
              role="presentation"
              class="illustration auth__illustration"
              src={getAssetPath('./assets/logout.svg')}
              alt=""
            />
            <h3 class="auth__title" data-test="logout:title">
              {t('auth.loggedOut.title')}
            </h3>
            <a class="button button--primary" {...href(ROUTES.AUTH)} data-test="logout:login">
              {t('auth.loggedOut.login')}
            </a>
          </div>
        )}
      </Host>
    );
  }
}

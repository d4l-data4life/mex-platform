import { Component, h, Host } from '@stencil/core';
import {
  AUTH_PROVIDER,
  FEATURE_FLAG_LANGUAGE_SWITCHER,
  LANGUAGES,
  LANGUAGE_CODES,
  NAVIGATION_CONFIG,
  ROUTES,
} from 'config';
import { AuthServiceProvider } from 'config/auth';
import { href } from 'stencil-router-v2';
import stores from 'stores';

@Component({
  tag: 'mex-header',
  styleUrl: 'header.css',
})
export class HeaderComponent {
  #showLogout: boolean;

  async componentWillLoad() {
    this.#showLogout = (await AUTH_PROVIDER()) === AuthServiceProvider.MICROSOFT;
  }

  render() {
    const { t, language } = stores.i18n;

    return (
      <Host class="header" data-test="header">
        <a {...href(ROUTES.ROOT)} class="header__logo" title={t('navigation.pages.home')} data-test="header:logo">
          <mex-logo flag="beta" flagTooltip={t('navigation.beta')} />
        </a>

        {FEATURE_FLAG_LANGUAGE_SWITCHER() && (
          <ul class="header__languages" data-test="header:languages">
            {(NAVIGATION_CONFIG.LANGUAGE_SWITCHER_LANGUAGES ?? LANGUAGE_CODES).map((code) => (
              <li class={{ header__language: true, 'header__language--active': language === code }}>
                <button
                  onClick={() => (stores.i18n.language = code)}
                  title={t('navigation.language', { language: LANGUAGES[code] })}
                  aria-label={t('navigation.language', { language: LANGUAGES[code] })}
                  aria-current={language === code ? 'page' : null}
                  data-test="header:language"
                  data-test-context={code}
                  data-test-active={String(language === code)}
                >
                  {code.toUpperCase()}
                </button>
              </li>
            ))}
          </ul>
        )}

        {stores.auth.isAuthenticated && this.#showLogout && <mex-header-user class="header__user" />}
      </Host>
    );
  }
}

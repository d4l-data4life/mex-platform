import { Component, getAssetPath, h, Host } from '@stencil/core';
import { ROUTES } from 'config';
import { href } from 'stencil-router-v2';
import stores from 'stores';

@Component({
  tag: 'mex-view-not-found',
  styleUrl: './not-found.css',
  assetsDirs: ['assets'],
})
export class NotFoundView {
  render() {
    const { t } = stores.i18n;

    return (
      <Host class="view not-found" data-test="notFound">
        <div class="view__wrapper">
          <img
            aria-hidden="true"
            role="presentation"
            class="illustration not-found__illustration"
            src={getAssetPath('./assets/error.svg')}
            alt=""
          />

          <h4 class="not-found__title">{t('notFound.title')}</h4>
          <p class="not-found__text">{t('notFound.text')}</p>

          <a class="button" {...href(ROUTES.ROOT)} data-test="notFound:button">
            {t('notFound.button')}
          </a>
        </div>
      </Host>
    );
  }
}

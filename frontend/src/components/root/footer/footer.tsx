import { Component, getAssetPath, h, Host, Prop } from '@stencil/core';
import { APP_VERSION, NAVIGATION_CONFIG } from 'config';
import stores from 'stores';

@Component({
  tag: 'mex-footer',
  styleUrl: 'footer.css',
  assetsDirs: ['assets'],
})
export class FooterComponent {
  @Prop() isThin = false;

  render() {
    const { isThin } = this;
    return (
      <Host
        class={`footer${isThin ? ' footer--thin' : ''}`}
        style={{
          '--background-pattern-left': `url(${getAssetPath('./assets/header-background-pattern-1.svg')})`,
          '--background-pattern-right': `url(${getAssetPath('./assets/header-background-pattern-2.svg')})`,
        }}
        data-test="footer"
        role="contentinfo"
      >
        {!isThin && <p class="footer__pre">{stores.i18n.t('navigation.disclaimer')}</p>}
        <div class="footer__inner">
          <span class="footer__note">{stores.i18n.t('navigation.version', { version: APP_VERSION })}</span>
          <nav
            class="footer__navigation"
            data-test="footer:navigation"
            aria-label={stores.i18n.t('navigation.footer.label')}
          >
            <mex-links
              classes="footer__navigation-container"
              itemClasses="footer__navigation-item"
              items={NAVIGATION_CONFIG[stores.i18n.language]?.FOOTER ?? []}
            />
          </nav>
        </div>
      </Host>
    );
  }
}

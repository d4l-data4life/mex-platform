import { Component, h, Host } from '@stencil/core';
import { CardRow } from 'components/card/card';
import { ANALYTICS_CUSTOM_EVENTS, NAVIGATION_CONFIG } from 'config';
import services from 'services';
import stores from 'stores';
import { getSearchUrl } from 'utils/search';

@Component({
  tag: 'mex-view-home',
})
export class HomeView {
  get supportCardRows(): CardRow[] {
    return [
      {
        type: 'text',
        value: stores.i18n.t('support.text'),
      },
      {
        type: 'action',
        label: stores.i18n.t('support.contact'),
        value: `mailto:${NAVIGATION_CONFIG.SUPPORT_EMAIL}?subject=${encodeURIComponent(
          stores.i18n.t('support.subject') as string
        )}`,
        testAttr: 'support:button',
        onClick: () => services.analytics.trackEvent(...ANALYTICS_CUSTOM_EVENTS.HOME_SUPPORT_FORM),
      },
    ];
  }

  render() {
    const { language } = stores.i18n;
    return (
      <Host class="view">
        <mex-search
          expanded
          handleSearch={(value: string, field: string) => {
            stores.search.query = value;
            stores.search.focus = field;
            stores.router.push(getSearchUrl(true));
          }}
        />
        <mex-browse class="u-spacing-when-empty" />
        <div class="view__split view__wrapper">
          <mex-dashboard />
          <div>
            <mex-linklist
              tile={{ items: NAVIGATION_CONFIG[language]?.LEARNING_ENVIRONMENT ?? [] }}
              headline={stores.i18n.t('navigation.learningEnvironment')}
            />
            {!!NAVIGATION_CONFIG.SUPPORT_EMAIL && (
              <mex-card rows={this.supportCardRows} data-test="support" icon="support" />
            )}
          </div>
          {!!NAVIGATION_CONFIG[language]?.SERVICES.length && (
            <mex-servicelist items={NAVIGATION_CONFIG[language].SERVICES} />
          )}
        </div>
      </Host>
    );
  }
}

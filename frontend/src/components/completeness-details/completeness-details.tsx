import { Component, Fragment, h, Prop, State } from '@stencil/core';
import { ITEM_FIELDS_CONCATENATOR, URL_METADATA_COMPLETENESS_INFO } from 'config';
import stores from 'stores';
import { CompletenessData, translateFieldName } from 'utils/field';

@Component({
  tag: 'mex-completeness-details',
  styleUrl: 'completeness-details.css',
})
export class CompletenessDetailsComponent {
  @Prop() data: CompletenessData;
  @Prop() context: any;

  @State() activeCategory: string = this.missingFieldsCategories[0];

  get missingFieldsCategories() {
    const { missingFields } = this.data;
    return Object.keys(missingFields).filter((category) => missingFields[category]?.length);
  }

  render() {
    const { t } = stores.i18n;
    const { missingFieldsCategories } = this;
    const { scores, weights, missingFields } = this.data;

    return (
      <Fragment>
        <div class="completeness-details__total">
          <mex-icon-completeness class="completeness-details__bars" value={scores.total} monochrome />
          <strong class="u-underline-4">{t('item.completeness.score', { score: scores.total })}</strong>
        </div>
        <p class="completeness-details__intro" innerHTML={t('item.completeness.intro', weights)} />

        {!!missingFieldsCategories.length && (
          <Fragment>
            <h5 class="completeness-details__missing-headline">{t('item.completeness.missing')}</h5>
            <mex-tabs
              class="completeness-details__missing-tabs"
              items={missingFieldsCategories.map((category) => ({
                label: t(`item.completeness.tabs.${category}`, { count: missingFields[category].length }),
                isActive: this.activeCategory === category,
                data: category,
              }))}
              handleClick={(item) => (this.activeCategory = item.data as string)}
            />
            <div class="completeness-details__missing-fields">
              {missingFields[this.activeCategory]
                .map((field) => translateFieldName(field))
                .join(ITEM_FIELDS_CONCATENATOR.default)}
            </div>
          </Fragment>
        )}

        <a
          class="completeness-details__more"
          href={URL_METADATA_COMPLETENESS_INFO}
          rel="noopener noreferrer"
          target="_blank"
        >
          {t('item.completeness.more')}
        </a>

        <button
          class="completeness-details__close button"
          ref={(el) => (this.context.modalCloseButtonRef = el)}
          onClick={() => this.context.closeModal.emit()}
        >
          {t('item.completeness.close')}
        </button>
      </Fragment>
    );
  }
}

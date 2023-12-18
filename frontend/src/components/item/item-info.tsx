import { Component, Event, EventEmitter, Fragment, h, Prop } from '@stencil/core';
import { ModalData } from 'components/root/root';
import { FIELDS, FEATURE_FLAGS } from 'config';
import { SidebarFeature, SidebarFeatureConfig } from 'config/item';
import { Item } from 'services/item';
import stores from 'stores';
import { SearchResultsItem } from 'stores/search';
import {
  calculateCompleteness,
  translateFieldName,
  getConcatDisplayValue,
  getValue,
  hasValue,
  normalizeConceptId,
} from 'utils/field';

@Component({
  tag: 'mex-item-info',
  styleUrl: 'item-info.css',
})
export class ItemInfoComponent {
  modalCloseButtonRef: HTMLButtonElement;

  @Prop() item: SearchResultsItem | Item;
  @Prop() orientation: 'horizontal' | 'vertical' | 'auto' = 'vertical';
  @Prop() features: SidebarFeatureConfig[] = [];

  @Event() showModal: EventEmitter<ModalData>;
  @Event() closeModal: EventEmitter;

  get completeness() {
    return calculateCompleteness(this.item);
  }

  completenessDetailsModalFocusHandler(closeRef?: HTMLButtonElement) {
    (this.modalCloseButtonRef ?? closeRef)?.focus();
  }

  showCompletenessDetails() {
    this.showModal.emit({
      Contents: () => <mex-completeness-details data={this.completeness} context={this} />,
      handleSetFocus: this.completenessDetailsModalFocusHandler.bind(this),
    });
  }

  renderAccessControl(config: SidebarFeatureConfig) {
    const { item } = this;
    const fields = config.field.resolvesTo;

    return fields.map((field) => {
      const value = normalizeConceptId(getValue(field, this.item));
      const icon = value && config.iconMap?.find(({ conceptid }) => value === conceptid)?.icon;

      return (
        (icon === 'locked' || icon === 'unlocked') && (
          <div class="item-info__row">
            <mex-icon-lock
              class={{
                'u-underline-3': true,
                'u-underline--positive': icon === 'unlocked',
                'u-underline--negative': icon === 'locked',
              }}
              classes="item-info__icon icon--large"
              open={icon === 'unlocked'}
            />
            <div class="item-info__data">
              <strong>{getConcatDisplayValue([field], item)}</strong>
              {translateFieldName(field, item)}
            </div>
          </div>
        )
      );
    });
  }

  renderCompleteness() {
    const { completeness } = this;
    const completenessTotalScore = completeness.scores.total;

    return (
      <div
        class="item-info__row"
        onClick={() => FEATURE_FLAGS.METADATA_COMPLETENESS_DETAILS && this.showCompletenessDetails()}
        style={{ cursor: FEATURE_FLAGS.METADATA_COMPLETENESS_DETAILS ? 'pointer' : 'default' }}
      >
        <mex-icon-completeness value={completenessTotalScore} monochrome />
        <div class="item-info__data">
          <strong>{stores.i18n.t('item.completeness.score', { score: completenessTotalScore })}</strong>
          {translateFieldName(FIELDS.completeness)}
        </div>
      </div>
    );
  }

  renderDate(config: SidebarFeatureConfig) {
    const { item } = this;
    const fields = config.field?.resolvesTo ?? [];
    const displayedFields = fields
      .map((field) => (hasValue(field.cloneAndSetRawValueMode(), item) ? field.cloneAndSetRawValueMode() : field))
      .filter((field) => hasValue(field, item));
    const displayedDate = !!displayedFields.length && getConcatDisplayValue(displayedFields, item, null, false);

    return (
      displayedDate && (
        <div class="item-info__row">
          <mex-icon-calendar class="u-underline-3" classes="item-info__icon icon--large" arrows={fields.length > 1} />
          <div class="item-info__data">
            <strong>{displayedDate}</strong>
            {translateFieldName(config.field, item)}
          </div>
        </div>
      )
    );
  }

  renderFeature(config: SidebarFeatureConfig) {
    const { feature } = config;
    return (
      <Fragment>
        {feature === SidebarFeature.accessRestriction && this.renderAccessControl(config)}
        {feature === SidebarFeature.completeness && this.renderCompleteness()}
        {feature === SidebarFeature.date && this.renderDate(config)}
      </Fragment>
    );
  }

  render() {
    const { orientation, features } = this;

    return (
      <Fragment>
        <div class={`item-info item-info--${orientation}`}>{features.map((config) => this.renderFeature(config))}</div>
      </Fragment>
    );
  }
}

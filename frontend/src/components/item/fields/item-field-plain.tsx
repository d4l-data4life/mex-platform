import { Component, Fragment, h, Host, Prop } from '@stencil/core';
import { FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY } from 'config';
import { Field, FieldRenderer } from 'models/field';
import { Item } from 'services/item';
import stores from 'stores';
import {
  translateFieldName,
  hasValueChanged,
  translateFieldDescription,
  getConcatDisplayValue,
  getValues,
  getDisplayValues,
  getConcatenator,
  getRawValues,
  translateFieldValueDescription,
  normalizeKey,
  getLangAttrIfForeign,
} from 'utils/field';

@Component({
  tag: 'mex-item-field-plain',
  styleUrl: 'item-field.css',
})
export class ItemFieldPlainComponent {
  @Prop() field: Field;
  @Prop() item?: Item;
  @Prop() items?: Item[];
  @Prop() previousVersionItem?: Item;
  @Prop() highlightChanges = false;
  @Prop() refErrors?: string[];

  get renderer() {
    return this.field.renderer ?? FieldRenderer.plain;
  }

  get isLinked() {
    return !this.item || !!this.items;
  }

  isLink(value?: string) {
    const supportedProtocols = ['https', 'http'];
    return (
      !!value &&
      this.renderer === FieldRenderer.link &&
      supportedProtocols.some((protocol) => value.indexOf(`${protocol}://`) === 0)
    );
  }

  get descriptionTooltip() {
    return translateFieldDescription(this.field, this.item);
  }

  getValueTooltip(value: string) {
    const valueDescription = translateFieldValueDescription(value);
    return valueDescription !== value && valueDescription !== normalizeKey(value) ? valueDescription : null;
  }

  get isDiffHighlighted() {
    if (this.isLinked) {
      return this.highlightChanges && hasValueChanged(this.field, this.item, this.previousVersionItem);
    }

    return (
      this.highlightChanges &&
      this.field.resolvesTo.some((field) => hasValueChanged(field, this.item, this.previousVersionItem))
    );
  }

  get usePluralLabel() {
    const { item } = this;

    return this.isLinked
      ? this.items?.length > 1
      : !!item && this.field.resolvesTo.flatMap((field) => getValues(field, item)).length > 1;
  }

  getConcatDisplayValue(item: Item) {
    return getConcatDisplayValue((this.field.linkedField ?? this.field).resolvesTo, item);
  }

  getRawValues(item: Item) {
    const rawValues = getRawValues((this.field.linkedField ?? this.field).resolvesTo, item);
    return rawValues.length ? rawValues : [null];
  }

  getLang(item: Item, summarize = true) {
    return getLangAttrIfForeign((this.field.linkedField ?? this.field).resolvesTo, item, summarize);
  }

  getDisplayValues(item: Item) {
    const displayValues = getDisplayValues((this.field.linkedField ?? this.field).resolvesTo, item);
    return displayValues.length ? displayValues : [stores.i18n.t(FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY)];
  }

  render() {
    const {
      item,
      items,
      usePluralLabel,
      renderer,
      isLinked,
      field,
      descriptionTooltip,
      getValueTooltip,
      isDiffHighlighted,
      refErrors,
    } = this;
    const concatenator = getConcatenator([this.field]);
    const rawValues = !isLinked && item && this.getRawValues(item);
    const langAttrs = !isLinked && item && (this.getLang(item, false) as string[]);
    const displayValues = !isLinked && item && this.getDisplayValues(item);

    return (
      <Host
        class={`item-field item-field--plain ${
          renderer !== FieldRenderer.plain ? `item-field--plain-${renderer}` : ''
        }`}
      >
        <div class="item-field__tooltip">
          {!!descriptionTooltip && (
            <mex-tooltip text={descriptionTooltip} testAttr="item:field:tooltip" data-test-context={field.name} />
          )}
        </div>
        <div class="item-field__name">{translateFieldName(field, item, usePluralLabel)}</div>
        <div class={`item-field__value u-highlight ${isDiffHighlighted ? 'u-highlight-2' : ''}`}>
          {rawValues &&
            rawValues.map((rawValue, index) => {
              const valueTooltip = getValueTooltip(rawValue);

              return (
                <Fragment>
                  {!!index && concatenator}
                  {valueTooltip && (
                    <mex-tooltip key={index} class="item-field__value-tooltip" text={valueTooltip}>
                      <span slot="toggle" innerHTML={displayValues[index]} />
                    </mex-tooltip>
                  )}
                  {!valueTooltip &&
                    (this.isLink(rawValue) ? (
                      <a
                        key={index}
                        href={rawValue}
                        target="_blank"
                        rel="noopener noreferrer"
                        innerHTML={displayValues[index]}
                      />
                    ) : (
                      <span key={index} innerHTML={displayValues[index]} lang={langAttrs[index]} />
                    ))}
                </Fragment>
              );
            })}

          {isLinked &&
            items?.map((item, index) => (
              <Fragment>
                {!!index && concatenator}
                <span key={index} innerHTML={this.getConcatDisplayValue(item)} lang={this.getLang(item) as string} />
              </Fragment>
            ))}

          {isLinked && !items && <mex-placeholder text={'…'.repeat(20)} />}
          {items && !items.length && <span>{stores.i18n.t(FIELDS_EMPTY_VALUE_FALLBACK_I18N_KEY)}</span>}

          {!!refErrors?.length && (
            <mex-tooltip
              class="item-field__inline-error"
              text={stores.i18n.t('item.referenceError', {
                field: translateFieldName(field),
                identifiers: refErrors.join(', '),
              })}
            >
              <span class="item-field__inline-error-toggle" slot="toggle">
                <mex-icon-warning classes="icon--medium" />
              </span>
            </mex-tooltip>
          )}
        </div>
      </Host>
    );
  }
}

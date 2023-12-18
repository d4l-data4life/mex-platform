import { Component, Fragment, h, Prop, State } from '@stencil/core';
import { FEATURE_FLAGS, FIELDS, ITEM_CONFIG } from 'config';
import { Field, FieldRenderer } from 'models/field';
import { Item } from 'services/item';
import stores from 'stores';
import { getInvolvedFields, hasValue } from 'utils/field';
import { getSearchUrl } from 'utils/search';

@Component({
  tag: 'mex-item-fields',
})
export class ItemFieldsComponent {
  @Prop() item?: Item;
  @Prop() level: number = 1;
  @Prop() previousVersionItem?: Item;
  @Prop() highlightChanges = false;
  @Prop() context: 'details' | 'preview' | 'sidebar' = 'details';
  @Prop() fields?: Field[];

  @State() expandDebug = false;

  get entityType() {
    return this.item?.entityType;
  }

  get unusedFields(): Field[] {
    const { item } = this;

    const involvedFields = getInvolvedFields(this.entityType);

    return Object.values(FIELDS)
      .filter((field) => !involvedFields.includes(field))
      .filter((field) => hasValue(field, item))
      .filter((field) => field.renderer !== FieldRenderer.none);
  }

  get fieldsToRender(): Field[] {
    if (this.fields) {
      return this.fields;
    }

    return ITEM_CONFIG.DISPLAYED_FIELDS[this.entityType] ?? ITEM_CONFIG.DISPLAYED_FIELDS.default ?? [];
  }

  get isKnownEntityType() {
    return !!ITEM_CONFIG.DISPLAYED_FIELDS[this.item?.entityType];
  }

  renderField(field: Field) {
    const { item, previousVersionItem, highlightChanges, isKnownEntityType, level, context } = this;
    const props = {
      class: `item__field item-field--${context}`,
      key: field.name,
      field,
      item,
      previousVersionItem,
      highlightChanges,
    };

    const renderer = field.renderer ?? FieldRenderer.plain;
    const isLinkedField = renderer === FieldRenderer.entity || renderer === FieldRenderer.reference;

    if (level > 1 && isLinkedField) {
      return null;
    }

    if (!isKnownEntityType) {
      return hasValue(field, item) && renderer !== FieldRenderer.none ? <mex-item-field-plain {...props} /> : null;
    }

    if (context !== 'details') {
      switch (renderer) {
        case FieldRenderer.title:
        case FieldRenderer.description:
        case FieldRenderer.plain:
        case FieldRenderer.link:
        case FieldRenderer.time:
        case FieldRenderer.bullets:
          return <mex-item-field-plain {...props} />;
        case FieldRenderer.entity:
        case FieldRenderer.reference:
          return <mex-item-field-reference {...props} />;
        case FieldRenderer.none:
          return null;
        default:
          throw new Error(`Renderer missing for field ${field.name}`);
      }
    }

    switch (renderer) {
      case FieldRenderer.title:
        return <mex-item-field-title {...props} />;
      case FieldRenderer.description:
        return <mex-item-field-description {...props} />;
      case FieldRenderer.plain:
      case FieldRenderer.link:
      case FieldRenderer.time:
      case FieldRenderer.bullets:
        return <mex-item-field-plain {...props} />;
      case FieldRenderer.reference:
        return <mex-item-field-reference {...props} />;
      case FieldRenderer.entity:
        return <mex-item-field-entity {...props} />;
      case FieldRenderer.none:
        return null;
      default:
        throw new Error(`Renderer missing for field ${field.name}`);
    }
  }

  render() {
    const { item, context, fieldsToRender } = this;
    const { t } = stores.i18n;

    return (
      <Fragment>
        {item && context === 'details' && <mex-item-entity-headline item={item} />}

        {item && fieldsToRender.map((field) => this.renderField(field))}

        {item === undefined && <mex-placeholder lines={30} />}

        {context === 'details' && (
          <Fragment>
            {item === null && (
              <mex-not-found-state
                class="results__empty-state"
                testAttr="item:details:notFound"
                caption={t('item.notFound.title')}
                text={t('item.notFound.text')}
                buttonText={t('item.notFound.button')}
                buttonUrl={getSearchUrl(false, stores.search.page)}
              />
            )}

            {FEATURE_FLAGS.DEBUG_UNUSED_FIELDS && this.level === 1 && !!this.unusedFields.length && (
              <div class="item__debug">
                <mex-toggle
                  label="[DEBUG] show unused populated fields"
                  active={this.expandDebug}
                  toggleHandler={(active) => (this.expandDebug = active)}
                />
                <mex-accordion class="item-navigation__row" expanded={this.expandDebug}>
                  {this.unusedFields.map((field) => this.renderField(field))}
                </mex-accordion>
              </div>
            )}
          </Fragment>
        )}
      </Fragment>
    );
  }
}

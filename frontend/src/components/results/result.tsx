import { Component, Fragment, h, Prop } from '@stencil/core';
import { ROUTES, SEARCH_CONFIG } from 'config';
import { Field, FieldRenderer } from 'models/field';
import { href } from 'stencil-router-v2';
import { SearchResultsItem } from 'stores/search';
import { translateFieldName, hasValue, getConcatDisplayValue, getLangAttrIfForeign } from 'utils/field';

@Component({
  tag: 'mex-result',
  styleUrl: 'result.css',
})
export class ResultComponent {
  @Prop() item?: SearchResultsItem;

  get entityType() {
    return this.item?.entityType;
  }

  get displayedFields(): Field[] {
    const { item } = this;
    return (SEARCH_CONFIG.DISPLAYED_FIELDS[this.entityType] ?? SEARCH_CONFIG.DISPLAYED_FIELDS.default).filter(
      (field) =>
        field.renderer === FieldRenderer.title ||
        field.resolvesTo.some((resolvedField) => hasValue(resolvedField, item))
    );
  }

  get titleField(): Field {
    return this.displayedFields.find((field) => field.renderer === FieldRenderer.title);
  }

  renderField(field: Field) {
    const { item, entityType } = this;
    const resolvedFields = field.resolvesTo;

    switch (field.renderer) {
      case FieldRenderer.title:
        return (
          <div class="result__title">
            <mex-icon-entity
              class="result__icon"
              entityName={entityType}
              attrs={{ classes: 'icon--large u-underline-3' }}
            />
            <h5 class="u-highlight" lang={getLangAttrIfForeign(resolvedFields, item) as string}>
              <a
                class="result__title-link"
                {...(item.itemId ? href(ROUTES.ITEM.replace(':id', encodeURIComponent(item.itemId))) : {})}
                innerHTML={getConcatDisplayValue(resolvedFields, item, item.itemId)}
                data-test="result:link"
              />
            </h5>
          </div>
        );

      case FieldRenderer.description:
        return (
          <p
            class="result__description u-highlight"
            innerHTML={getConcatDisplayValue(resolvedFields, item, '')}
            lang={getLangAttrIfForeign(resolvedFields, item) as string}
          />
        );

      default:
        return (
          <li class="u-highlight">
            <strong>{translateFieldName(field, item)}:</strong>
            <span
              innerHTML={getConcatDisplayValue(resolvedFields, item)}
              lang={getLangAttrIfForeign(resolvedFields, item) as string}
            />
          </li>
        );
    }
  }

  render() {
    const { item, entityType, displayedFields, titleField } = this;

    return (
      <div class="result" data-test="result" data-test-key={item?.itemId}>
        {item ? (
          <Fragment>
            <div class="result__content">
              <ul class="result__fields">{displayedFields.map((field) => this.renderField(field))}</ul>
            </div>
            <aside
              class="result__info"
              aria-label={titleField ? getConcatDisplayValue(titleField.resolvesTo, item, item.itemId) : item.itemId}
            >
              <mex-item-info
                item={item}
                orientation="auto"
                features={SEARCH_CONFIG.SIDEBAR_FEATURES[entityType] ?? SEARCH_CONFIG.SIDEBAR_FEATURES.default}
              />
            </aside>
          </Fragment>
        ) : (
          <div class="result__content" data-test="result:placeholders">
            <mex-placeholder lines={SEARCH_CONFIG.DISPLAYED_FIELDS.default.length + 2} />
          </div>
        )}
      </div>
    );
  }
}

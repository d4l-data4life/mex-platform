import { Component, Fragment, h, Prop, State } from '@stencil/core';
import { ENTITY_TYPES, ITEM_CONFIG, ROUTES } from 'config';
import services from 'services';
import { Item } from 'services/item';
import stores from 'stores';
import { getValues } from 'utils/field';
import { Field } from 'models/field';

@Component({
  tag: 'mex-item-field-entity',
  styleUrl: 'item-field.css',
})
export class ItemFieldEntityComponent {
  @Prop() field: Field;
  @Prop() item?: Item;

  @State() linkedItems: Item[];

  get isLoading() {
    return !this.linkedItems?.length;
  }

  get values() {
    return this.field.resolvesTo
      .flatMap((field) => getValues(field, this.item) ?? [])
      .filter(Boolean)
      .map(({ fieldValue }) => fieldValue)
      .filter(Boolean);
  }

  isSupportedEntityType(item?: Item) {
    return ITEM_CONFIG.DEDICATED_VIEW_SUPPORTED_ENTITY_TYPES.includes(item?.entityType);
  }

  async componentDidLoad() {
    const { values } = this;
    this.linkedItems = await Promise.all(
      values.map((identifier) => stores.items.get(identifier) ?? services.item.resolveLink(identifier))
    );
    this.linkedItems.forEach((item, index) => item && stores.items.add(values[index], item));
  }

  renderIndicatorIconHtml(item: Item) {
    if (!ENTITY_TYPES[item.entityType].config?.icon) {
      return null;
    }

    const attrs = {
      class: 'u-underline-3',
      classes: 'icon--large',
    };

    return (
      <div class="item-field__entity-icon">
        <mex-icon-entity entityName={item.entityType} attrs={attrs} />
      </div>
    );
  }

  render() {
    return (
      <Fragment>
        {this.isLoading
          ? this.values.map((_, index) => (
              <div class="item-field item-field--entity" key={index}>
                <mex-placeholder lines={8} />
              </div>
            ))
          : this.linkedItems.map((item, index) => (
              <div class="item-field item-field--entity" key={index}>
                {item ? (
                  <Fragment>
                    <mex-item-entity-headline
                      item={item}
                      link={
                        this.isSupportedEntityType(item)
                          ? ROUTES.ITEM.replace(':id', encodeURIComponent(item.itemId))
                          : null
                      }
                    />
                    <mex-item-fields
                      item={item}
                      level={2}
                      data-test="item:details:entityDetails"
                      data-test-context={item.itemId}
                      context="preview"
                    />
                  </Fragment>
                ) : (
                  <div class="item-field__error">
                    <mex-icon-warning classes="icon--medium icon--inline" />
                    {stores.i18n.t('item.entityError', {
                      identifier: this.values[index],
                      interpolation: { escapeValue: false },
                    })}
                  </div>
                )}
              </div>
            ))}
      </Fragment>
    );
  }
}

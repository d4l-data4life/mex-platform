import { Component, h, Host, Prop } from '@stencil/core';
import { ENTITY_TYPES, FIELDS } from 'config';
import { Item } from 'services/item';
import { href } from 'stencil-router-v2';
import { formatValue, translateFieldValueDescription } from 'utils/field';

@Component({
  tag: 'mex-item-entity-headline',
  styleUrl: 'item-entity-headline.css',
})
export class ItemEntityHeadlineComponent {
  @Prop() item?: Item;
  @Prop() link?: string;

  get indicatorIconHtml() {
    const { entityType } = this.item;

    if (!ENTITY_TYPES[entityType]?.config?.icon) {
      return null;
    }

    const attrs = {
      class: 'u-underline-3',
      classes: 'icon--large',
    };

    return (
      <div class="item-entity-headline__icon">
        <mex-icon-entity entityName={entityType} attrs={attrs} />
      </div>
    );
  }

  render() {
    const { item, link } = this;

    return (
      !!item && (
        <Host class="item-entity-headline__wrapper">
          {this.indicatorIconHtml}

          <div class="item-entity-headline__row">
            {!!link ? (
              <a
                class="item-entity-headline__link"
                {...href(link)}
                data-test="item:details:entityLink"
                data-test-context={item.itemId}
              >
                <h4 class="item-entity-headline__title">
                  {formatValue([item.entityType], FIELDS.entityName)}
                  <mex-icon-arrow classes="icon--inline" />
                </h4>
              </a>
            ) : (
              <h4 class="item-entity-headline__title">{formatValue([item.entityType], FIELDS.entityName)}</h4>
            )}
            <p class="item-entity-headline__description">
              {translateFieldValueDescription(item.entityType, FIELDS.entityName)}
            </p>
          </div>
        </Host>
      )
    );
  }
}

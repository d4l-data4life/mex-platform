import { Component, h, Host, Prop, State } from '@stencil/core';
import { Field } from 'models/field';
import services from 'services';
import { Item } from 'services/item';
import stores from 'stores';
import { filterDuplicates, getValues } from 'utils/field';

@Component({
  tag: 'mex-item-field-reference',
  styleUrl: 'item-field.css',
})
export class ItemFieldReferenceComponent {
  @Prop() field: Field;
  @Prop() item?: Item;
  @Prop() previousVersionItem?: Item;
  @Prop() highlightChanges = false;

  @State() linkedItems: Item[];
  @State() refErrors: string[] = [];

  get renderedItems() {
    return this.linkedItems?.filter(Boolean);
  }

  get isLoading() {
    return !this.renderedItems;
  }

  get values(): string[] {
    return this.field.resolvesTo
      .flatMap((field) => getValues(field, this.item) ?? [])
      .filter(Boolean)
      .map(({ fieldValue }) => fieldValue)
      .filter(Boolean)
      .filter(filterDuplicates);
  }

  async componentDidLoad() {
    const values = this.values;
    this.linkedItems = await Promise.all(
      values.map((identifier) => stores.items.get(identifier) ?? services.item.resolveLink(identifier))
    );
    this.linkedItems.forEach((item, index) => item && stores.items.add(values[index], item));
    this.refErrors = values.filter((_, index) => !this.linkedItems[index]);
  }

  render() {
    const { field, item, previousVersionItem, highlightChanges, isLoading, renderedItems, refErrors } = this;
    const props = {
      class: 'item__field',
      key: field.name,
      field,
      item: isLoading ? null : item,
      items: renderedItems,
      previousVersionItem,
      highlightChanges,
      refErrors,
    };

    return (
      <Host class="item-field item-field--reference">
        <mex-item-field-plain {...props} />
      </Host>
    );
  }
}

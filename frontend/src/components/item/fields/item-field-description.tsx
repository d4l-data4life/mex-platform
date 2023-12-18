import { Component, h, Host, Prop } from '@stencil/core';
import { Field } from 'models/field';
import { Item } from 'services/item';
import { getConcatDisplayValue, getLangAttrIfForeign, hasValueChanged } from 'utils/field';

@Component({
  tag: 'mex-item-field-description',
  styleUrl: 'item-field.css',
})
export class ItemFieldDescriptionComponent {
  @Prop() field: Field;
  @Prop() item?: Item;
  @Prop() previousVersionItem?: Item;
  @Prop() highlightChanges = false;

  get isDiffHighlighted() {
    return (
      this.highlightChanges &&
      this.field.resolvesTo.some((field) => hasValueChanged(field, this.item, this.previousVersionItem))
    );
  }

  get displayValue() {
    return getConcatDisplayValue(this.field.resolvesTo, this.item);
  }

  get lang() {
    return getLangAttrIfForeign(this.field.resolvesTo, this.item) as string;
  }

  render() {
    const { displayValue, isDiffHighlighted, lang } = this;
    return (
      <Host class="item-field item-field--description">
        <p class={`item-field__value u-highlight ${isDiffHighlighted ? 'u-highlight-2' : ''}`} lang={lang}>
          <span innerHTML={displayValue} />
        </p>
      </Host>
    );
  }
}

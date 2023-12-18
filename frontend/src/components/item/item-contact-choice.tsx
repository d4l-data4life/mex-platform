import { Component, Event, EventEmitter, Prop, h } from '@stencil/core';
import { ModalData } from 'components/root/root';
import { SidebarFeatureConfig, SidebarFeatureConfigActionItem, SidebarFeatureContactAction } from 'config/item';
import { Field } from 'models/field';
import { Item } from 'services/item';
import stores from 'stores';
import { getValue } from 'utils/field';

@Component({
  tag: 'mex-item-contact-choice',
  styleUrl: 'item-contact-choice.css',
})
export class ItemContactChoiceComponent {
  @Prop() actions: SidebarFeatureConfigActionItem[];
  @Prop() config: SidebarFeatureConfig;
  @Prop() contactItem: Item;
  @Prop() contextItem: Item;

  @Prop() handleFormOpen?: () => void;
  @Prop() handleFormClose?: () => void;

  @Event() closeModal: EventEmitter;
  @Event() showModal: EventEmitter<ModalData>;

  get key() {
    return this.config.key;
  }

  getUrl(action: SidebarFeatureConfigActionItem): string {
    const { type, field } = action;
    const emailAddress = (field as Field)?.resolvesTo
      .map((field) => getValue(field, this.contactItem)?.trim())
      .find((value) => value?.includes('@'));

    if (type !== SidebarFeatureContactAction.email || !emailAddress) {
      return null;
    }

    return `mailto:${emailAddress}?subject=${encodeURIComponent(
      stores.i18n.t(`item.contact.${this.key}.email.subject`)
    )}`;
  }

  getClickHandler(action: SidebarFeatureConfigActionItem) {
    const { type, form } = action;
    if (type !== SidebarFeatureContactAction.form || !form) {
      return null;
    }

    return () => {
      this.handleFormOpen?.();
      this.showModal.emit({
        Contents: () => (
          <mex-form formId={form} formKey={this.key} recipient={this.contactItem} context={this.contextItem} embedded />
        ),
        handleClose: this.handleFormClose,
      });
    };
  }

  render() {
    const { t } = stores.i18n;
    const { actions, key } = this;

    return (
      <mex-modal-contents
        illustration="contact"
        caption={`item.contact.${key}.choiceTitle`}
        text={`item.contact.${key}.choiceText`}
        buttons={[
          {
            label: `item.contact.${key}.cancel`,
            clickHandler: () => this.closeModal.emit(),
            modifier: 'secondary',
          },
        ]}
      >
        <div class="item-contact-choice__tiles">
          {actions.map((action) => (
            <mex-tile
              class="item-contact-choice__tile"
              text={t(`item.contact.${key}.${action.type}.label`)}
              hint={t(`item.contact.${key}.${action.type}.text`)}
              modifiers={['text-left', 'button']}
              url={this.getUrl(action)}
              handleClick={this.getClickHandler(action)}
              icon="arrow"
            />
          ))}
        </div>
      </mex-modal-contents>
    );
  }
}

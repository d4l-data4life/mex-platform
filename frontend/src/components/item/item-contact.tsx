import { Component, Event, EventEmitter, h, Prop, State, Watch } from '@stencil/core';
import { CardRow } from 'components/card/card';
import { ModalData } from 'components/root/root';
import { ANALYTICS_CUSTOM_EVENTS } from 'config';
import { SidebarFeatureConfig, SidebarFeatureContactAction } from 'config/item';
import { Field } from 'models/field';
import services from 'services';
import { Item } from 'services/item';
import stores from 'stores';
import { getInvolvedFields, getValue, translateFieldName } from 'utils/field';

@Component({
  tag: 'mex-item-contact',
})
export class ItemContactComponent {
  @Prop() item?: Item;
  @Prop() config: SidebarFeatureConfig;
  @Prop() isOnLatestVersion: boolean = false;

  @State() sourceItem?: Item;

  @Watch('item')
  itemChangedHandler(newItem: Item, oldItem?: Item) {
    this.unsubscribeFromItemListener(oldItem);
    this.updateSourceItem(newItem);
  }

  @Event() showModal: EventEmitter<ModalData>;
  @Event() scrollingDisabled: EventEmitter<boolean>;

  get rowsFromDisplayField(): CardRow[] {
    const { sourceItem: item, config } = this;
    const { displayField } = config;
    if (!displayField) {
      return [];
    }

    const fields = displayField.resolvesTo;
    const values = fields.map((field) => getValue(field, item)?.trim());

    return fields.flatMap((field, index) => {
      const value = values[index];
      return value ? [{ type: 'value', label: translateFieldName(field, item), value }] : [];
    });
  }

  get rowsFromActions(): CardRow[] {
    const { sourceItem: item, config } = this;
    const { key } = config;
    if (!config.actions || !item) {
      return [];
    }

    const actions = config.actions.filter(
      ({ type, field, form }) =>
        (type === SidebarFeatureContactAction.email && (field as Field)?.isInitialized) ||
        (type === SidebarFeatureContactAction.form && form)
    );

    if (!actions.length) {
      return [];
    }

    const track = () => services.analytics.trackEvent(...ANALYTICS_CUSTOM_EVENTS.ITEM_CONTACT_FORM);
    const baseConfig = {
      type: 'action',
      label: stores.i18n.t(`item.contact.${key}.start`),
      testAttr: 'item:contact:button',
      onClick: track,
    };

    if (actions.length > 1) {
      return [
        {
          ...baseConfig,
          onClick: () => {
            track();

            this.showModal.emit({
              Contents: () => (
                <mex-item-contact-choice
                  config={config}
                  actions={actions}
                  contactItem={item}
                  contextItem={this.item}
                  handleFormOpen={this.handleFormOpen}
                  handleFormClose={this.handleFormClose}
                />
              ),
            });
          },
        },
      ];
    }

    const { type, field, form } = actions[0];
    const emailAddress = (field as Field)?.resolvesTo
      .map((field) => getValue(field, item)?.trim())
      .find((value) => value?.includes('@'));

    if (type === SidebarFeatureContactAction.email && emailAddress) {
      return [
        {
          ...baseConfig,
          value: `mailto:${emailAddress}?subject=${encodeURIComponent(
            stores.i18n.t(`item.contact.${key}.email.subject`)
          )}`,
        },
      ];
    }

    if (type === SidebarFeatureContactAction.form) {
      return [
        {
          ...baseConfig,
          onClick: () => {
            track();

            this.handleFormOpen();
            this.showModal.emit({
              Contents: () => <mex-form formId={form} formKey={key} recipient={item} context={this.item} embedded />,
              handleClose: this.handleFormClose,
            });
          },
        },
      ];
    }

    return [];
  }

  get rows() {
    return [...this.rowsFromDisplayField, ...this.rowsFromActions];
  }

  handleFormOpen = () => {
    this.scrollingDisabled.emit(true);
    window.onbeforeunload = () => '';
  };

  handleFormClose = () => {
    this.scrollingDisabled.emit(false);
    window.onbeforeunload = null;
  };

  updateSourceItem(item: Item = this.item) {
    const { field } = this.config;
    if (!field) {
      this.sourceItem = this.item;
      return;
    }

    const primaryBusinessIdentifier = getValue(field, item);
    if (!primaryBusinessIdentifier) {
      this.sourceItem = null;
      return;
    }

    stores.items.addListener(primaryBusinessIdentifier, (linkedItem) => (this.sourceItem = linkedItem));
    this.updateCacheIfWillNotPopulate(item, field, primaryBusinessIdentifier);
  }

  async updateCacheIfWillNotPopulate(item: Item, field: Field, businessIdentifier: string) {
    const involvedFields = getInvolvedFields(item.entityType, false);
    if (this.isOnLatestVersion && field.resolvesTo.every((field) => involvedFields.includes(field))) {
      return;
    }

    const sourceItem = await services.item.resolveLink(businessIdentifier);
    sourceItem && stores.items.add(businessIdentifier, sourceItem);
  }

  unsubscribeFromItemListener(item: Item) {
    const identifier = this.config.field && item && getValue(this.config.field, item);
    identifier && stores.items.removeListener(identifier);
  }

  componentWillLoad() {
    this.updateSourceItem();
  }

  disconnectedCallback() {
    this.unsubscribeFromItemListener(this.item);
  }

  render() {
    const { rows } = this;
    return !!rows.length && <mex-card rows={rows} data-test="item:contact" icon="email" />;
  }
}

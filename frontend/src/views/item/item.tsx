import { Component, h, Host, Listen, Prop, State, Watch } from '@stencil/core';
import { ROUTES } from 'config';
import services from 'services';
import { Item, Version } from 'services/item';
import stores from 'stores';
import { catchRetryableAction } from 'utils/error';
import { getSearchUrl } from 'utils/search';
import { ResponseError } from 'models/response-error';

@Component({
  tag: 'mex-view-item',
})
export class ItemView {
  #itemPromise?: Promise<Item>;

  @Prop() itemId: string;

  @State() item?: Item;
  @State() previousVersionItem?: Item;
  @State() latestVersionItem?: Item;
  @State() versions?: Version[];
  @State() expandVersions = false;

  @Listen('changeItemId')
  changeItemIdListener(event: CustomEvent) {
    const itemId = event.detail as string;
    itemId && stores.router.push(ROUTES.ITEM.replace(':id', encodeURIComponent(itemId)));
  }

  @Listen('toggleVersions')
  toggleVersionsListener(event: CustomEvent) {
    this.expandVersions = event.detail as boolean;

    if (!event.detail && !!this.newerVersion) {
      this.changeItemIdListener(new CustomEvent('changeItemId', { detail: this.versions[0].itemId }));
    }
  }

  @Watch('itemId')
  watchItemIdHandler() {
    this.populateData();
  }

  get hasCorrectVersions() {
    const { itemId } = this;
    return !!this.versions?.some((item) => item.itemId === itemId);
  }

  get newerVersion() {
    return this.versions?.[this.versions?.indexOf(this.versions?.find(({ itemId }) => itemId === this.itemId)) - 1];
  }

  get isVersionsExpanded() {
    return this.expandVersions || !!this.newerVersion;
  }

  get displayedVersions(): Version[] {
    const { versions, item } = this;

    if (!versions || versions.length) {
      return versions;
    }

    if (!item) {
      return null;
    }

    return [{ itemId: this.itemId, createdAt: new Date(this.item.createdAt) }];
  }

  isNotFoundError(error: ResponseError) {
    return error.status === 404;
  }

  resetData() {
    this.item = undefined;
    this.previousVersionItem = undefined;
    this.latestVersionItem = undefined;
    if (!this.hasCorrectVersions) {
      this.versions = undefined;
    }
  }

  populateData() {
    this.resetData();

    catchRetryableAction(async () => {
      this.#itemPromise = services.item.fetch(this.itemId);
      try {
        this.item = await this.#itemPromise;
      } catch (error) {
        if (!this.isNotFoundError(error)) {
          throw error;
        }

        this.item = null;
      }
    });

    catchRetryableAction(async () => {
      await this.populateVersions();
      await this.populateLatestVersionItem();
      await this.populatePreviousVersionItem();
    });
  }

  async populateVersions() {
    if (!this.hasCorrectVersions) {
      try {
        this.versions = await services.item.fetchVersions(this.itemId);
      } catch (error) {
        if (error.status !== 404) {
          throw error;
        }

        this.versions = [];
      }
    }
  }

  async populatePreviousVersionItem() {
    const previousVersion =
      this.versions?.[this.versions?.indexOf(this.versions?.find(({ itemId }) => itemId === this.itemId)) + 1];
    if (!previousVersion) {
      return;
    }

    this.previousVersionItem = await services.item.fetch(previousVersion.itemId);
  }

  async populateLatestVersionItem() {
    let item: Item;

    try {
      item = await this.#itemPromise;
    } catch (error) {
      if (!this.isNotFoundError(error)) {
        throw error;
      }

      return;
    }

    const latestVersion = this.versions?.[0];

    if (!latestVersion || latestVersion.itemId === item.itemId) {
      this.latestVersionItem = item;
      stores.items.add(item.itemId, item);
      return;
    }

    const cachedLatestVersionItem = stores.items.get(latestVersion.itemId);
    this.latestVersionItem = cachedLatestVersionItem ?? (await services.item.fetch(latestVersion.itemId));
    !cachedLatestVersionItem && stores.items.add(latestVersion.itemId, this.latestVersionItem);
  }

  async componentWillLoad() {
    await catchRetryableAction(async () => this.populateData());
  }

  render() {
    return (
      <Host class="view">
        <mex-search
          value={stores.search.query}
          searchFocus={stores.search.focus}
          handleSearch={(value: string, field: string) => {
            stores.search.query = value;
            stores.search.focus = field;
            stores.router.push(getSearchUrl(true));
          }}
        />
        <div class="view__wrapper">
          <mex-item-navigation
            itemId={this.itemId}
            item={this.item}
            versions={this.displayedVersions}
            expandVersions={this.isVersionsExpanded}
          />
          <mex-item
            item={this.item}
            previousVersionItem={this.previousVersionItem}
            latestVersionItem={this.latestVersionItem}
            highlightChanges={this.isVersionsExpanded}
          />
        </div>
      </Host>
    );
  }
}

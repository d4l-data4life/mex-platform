import { Component, Event, EventEmitter, Fragment, h, Prop, State, Watch } from '@stencil/core';
import { ModalData } from 'components/root/root';
import { FIELDS, SEARCH_DOWNLOAD_CONFIG, SEARCH_QUERY_EVERYTHING } from 'config';
import services from 'services';
import stores from 'stores';
import { SearchResultsItem } from 'stores/search';
import { formatResultsItemsAsCsv } from 'utils/search';

@Component({
  tag: 'mex-results-download',
  styleUrl: 'results-download.css',
})
export class ResultsDownloadComponent {
  modalContentButtonRefs: HTMLButtonElement[] = [];

  @Prop() disabled?: boolean = false;

  @State() isBusy: boolean = false;
  @State() progress?: number = null;

  @Event() showModal: EventEmitter<ModalData>;
  @Event() closeModal: EventEmitter;

  @Watch('isBusy')
  isBusyChangedHandler() {
    this.show();
  }

  @Watch('progress')
  progressChangedHandler() {
    this.show();
  }

  get isDisabled() {
    return this.disabled || this.isBusy || !this.config;
  }

  get hasResults() {
    return !!stores.search.count;
  }

  get entityTypeFacet() {
    const { facets } = stores.search;
    return facets?.find(({ axis }) => axis === FIELDS.entityName.linkedName);
  }

  get config() {
    const { entityTypeFacet, hasResults } = this;

    return SEARCH_DOWNLOAD_CONFIG.find(
      ({ entityType }) =>
        (!entityType && hasResults) ||
        entityTypeFacet?.buckets.some(({ value, count }) => value === entityType && !!count)
    );
  }

  createFileDownload(output: string) {
    const { config } = this;

    const byteArray = new TextEncoder().encode(output);
    const linkEl = document.createElement('a');
    linkEl.setAttribute('href', URL.createObjectURL(new Blob([byteArray], { type: 'text/csv' })));
    linkEl.setAttribute('download', `${config.key}.csv`);
    linkEl.style.display = 'none';

    document.body.appendChild(linkEl);
    linkEl.click();
    document.body.removeChild(linkEl);
  }

  async fetchResults(offset: number = 0, previousItems: SearchResultsItem[] = []) {
    const { config } = this;
    const { limit } = config;
    const { query } = stores.search;
    const filters = stores.filters.all
      .filter((filter) => !config.entityType || filter[0] !== FIELDS.entityName.linkedName)
      .concat(config.entityType ? [[FIELDS.entityName.linkedName, [config.entityType]]] : []);

    const response = await services.search.fetchResults({
      query: query || SEARCH_QUERY_EVERYTHING,
      offset,
      limit,
      filters,
      searchFocus: stores.search.focus,
      facets: [],
      fields: Object.values(FIELDS).filter((field) => !field.isVirtual),
      highlightFields: [],
    });

    if (!this.isBusy) {
      throw new Error('download aborted');
    }

    const count = response.numFound;
    const items = response.items.filter((item) => !config.entityType || item.entityType === config.entityType);

    this.progress = Math.min(Math.round(((offset + limit) / count) * 100), 100);

    if (count > offset + limit) {
      return await this.fetchResults(offset + limit, previousItems.concat(items));
    }

    return previousItems.concat(items);
  }

  async download() {
    if (this.isDisabled) {
      this.closeModal.emit();
      return;
    }

    const { config } = this;

    this.isBusy = true;
    this.progress = 0;

    try {
      const data = await this.fetchResults();
      const items = data; // circumvent js minification bug :(

      this.progress = 100;
      await new Promise((resolve) => requestAnimationFrame(() => setTimeout(resolve, 300)));

      const output = formatResultsItemsAsCsv(items);
      this.createFileDownload(output);

      stores.notifications.add(`download.${config.key}.success`);
    } catch (e) {
      console.error(e);
      this.isBusy && stores.notifications.add(`download.${config.key}.error`);
    } finally {
      this.isBusy = false;
      this.closeModal.emit();
    }
  }

  show() {
    const { config, isDisabled, progress } = this;

    this.showModal.emit({
      Contents: () => (
        <mex-modal-contents
          context={this}
          illustration="download"
          caption={`download.${config.key}.title`}
          text={progress === 100 ? `download.${config.key}.wait` : `download.${config.key}.text`}
          progress={progress}
          buttons={[
            {
              label: `download.${config.key}.cancel`,
              clickHandler: this.closeModal.emit,
              modifier: 'secondary',
              testAttr: 'results:download:modal:cancel',
            },
            {
              label: `download.${config.key}.start`,
              clickHandler: () => this.download(),
              disabled: isDisabled,
              testAttr: 'results:download:modal:start',
            },
          ]}
          data-test="results:download:modal"
        />
      ),
      handleSetFocus: (closeRef?: HTMLButtonElement) => (this.modalContentButtonRefs[0] ?? closeRef)?.focus(),
      handleClose: () => {
        this.isBusy = false;
        this.progress = null;
      },
    });
  }

  render() {
    const { config } = this;

    return (
      <Fragment>
        <button
          class="results-download__button"
          disabled={this.isDisabled}
          onClick={() => this.show()}
          title={config && stores.i18n.t(`download.${config.key}.title`)}
          data-test="results:download:button"
        >
          <mex-icon-download classes="icon--large" />
        </button>
      </Fragment>
    );
  }
}

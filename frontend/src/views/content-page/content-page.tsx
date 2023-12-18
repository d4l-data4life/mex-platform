import { Component, h, Host, Prop, State, Watch } from '@stencil/core';
import services from 'services';
import { ContentResponse } from 'services/content';
import stores from 'stores';
import { catchRetryableAction } from 'utils/error';

@Component({
  tag: 'mex-view-content-page',
})
export class ContentPageView {
  #contentPromise?: Promise<ContentResponse>;

  @Prop() pageId: string;

  @State() contentResponse?: ContentResponse;
  @State() isNotFound: boolean = false;

  @Watch('pageId')
  watchPageIdHandler() {
    this.populateData();
  }

  populateData() {
    catchRetryableAction(async () => {
      this.#contentPromise = services.content.fetch(this.pageId);

      try {
        this.contentResponse = await this.#contentPromise;
      } catch (error) {
        if (error.status !== 404) {
          throw error;
        }

        this.contentResponse = null;
        this.isNotFound = true;
      }
    });
  }

  componentWillLoad() {
    this.populateData();
  }

  render() {
    if (this.isNotFound) {
      return <mex-view-not-found />;
    }

    return (
      <Host class="view">
        {!!this.contentResponse && (
          <div class="view__wrapper">
            <mex-content-page content={this.contentResponse[stores.i18n.language]} />
          </div>
        )}
      </Host>
    );
  }
}

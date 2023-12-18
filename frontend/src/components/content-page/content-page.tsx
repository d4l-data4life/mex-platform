import { Component, Event, EventEmitter, Host, h, Prop, State } from '@stencil/core';
import { Content } from 'services/content';
import device from 'utils/device';

@Component({
  tag: 'mex-content-page',
  styleUrl: 'content-page.css',
})
export class ContentPageComponent {
  #contentEl: HTMLMexContentPageContentElement;

  @Prop() content: Content;

  @State() isViewportMobile: boolean = false;
  @State() activeAnchors: string[] = [];
  @State() scrollTarget?: string;

  @Event() staticPageTitleChanged: EventEmitter;
  @Event() stickyFooterEnabled: EventEmitter;

  constructor() {
    this.handleViewportChange = this.handleViewportChange.bind(this);
  }

  handleViewPortIntersect = (displayedHeadings: string[]) => {
    this.activeAnchors = [...displayedHeadings];
    this.scrollTarget = null;
  };

  handleAnchorClick = (target: string) => {
    this.scrollTarget = target.length > 1 && target[0] === '#' ? target.slice(1) : '';
  };

  handleViewportChange(isViewportMobile: boolean) {
    this.isViewportMobile = isViewportMobile;
    this.stickyFooterEnabled.emit(!isViewportMobile);
  }

  componentWillLoad() {
    device.mobileViewportChanges.addListener(this.handleViewportChange);
  }

  disconnectedCallback() {
    device.mobileViewportChanges.removeListener(this.handleViewportChange);
  }

  render() {
    const { title = '', toc = [] } = this.content;
    this.staticPageTitleChanged.emit(title);

    return (
      <Host class="content-page">
        <mex-content-page-nav
          tocTree={toc}
          activeAnchors={this.activeAnchors}
          contentEl={this.#contentEl}
          isMediaTablet={!this.isViewportMobile}
          handleAnchorClick={this.handleAnchorClick}
        />
        <mex-content-page-content
          ref={(el) => (this.#contentEl = el)}
          content={this.content}
          scrollTarget={this.scrollTarget}
          isMediaTablet={!this.isViewportMobile}
          handleViewPortIntersect={this.handleViewPortIntersect}
        />
      </Host>
    );
  }
}

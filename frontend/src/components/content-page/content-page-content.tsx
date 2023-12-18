import { Component, h, Host, Prop, Watch } from '@stencil/core';
import {
  Content,
  ContentBlock,
  ContentBlockType,
  ContentUnitHeading,
  ContentUnitText,
  ContentUnitInfobox,
  ContentUnitTable,
  ContentUnitImage,
  ContentUnitList,
  ContentUnitMarkdown,
  ContentUnitCode,
  ContentUnitEntityTypeHeadline,
  ContentUnitFieldDescription,
  ContentUnitCompletenessDocumentationTable,
} from 'services/content';
import { getPageClasses } from 'utils/content';
import {
  ContentBlockCode,
  ContentBlockCompletenessDocumentationTable,
  ContentBlockEntityTypeHeadline,
  ContentBlockFieldDescription,
  ContentBlockFigure,
  ContentBlockHeading,
  ContentBlockInfobox,
  ContentBlockList,
  ContentBlockMarkdown,
  ContentBlockTable,
  ContentBlockText,
} from 'components/content-blocks/content-blocks';

@Component({
  tag: 'mex-content-page-content',
  styleUrl: 'content-page-content.css',
  assetsDirs: ['assets'],
})
export class ContentPageContentComponent {
  #intersectionObserver: IntersectionObserver;
  #headingEls: HTMLElement[] = [];
  #hasPersistedDisplayedHeading: boolean = false;
  #displayedHeadings: string[] = [];
  #containerEl: HTMLElement;

  @Watch('scrollTarget')
  scrollToTarget(newScrollTarget) {
    if (newScrollTarget === null) {
      return;
    }

    if (this.isMediaTablet) {
      const targetAnchorEl = this.#headingEls.find(({ id }) => id === newScrollTarget);
      if (targetAnchorEl) {
        this.#containerEl.scrollTop = targetAnchorEl.offsetTop;
      }
    }
  }

  @Prop() content: Content;
  @Prop() scrollTarget: string;
  @Prop() isMediaTablet?: boolean = true;
  @Prop() handleViewPortIntersect?: (displayedHeadings: string[]) => void;

  getBlock(block: ContentBlock) {
    const { type = '', content = {} } = block;

    switch (type) {
      case ContentBlockType.heading:
        return (
          <ContentBlockHeading content={content as ContentUnitHeading} ref={(el) => el && this.#headingEls.push(el)} />
        );
      case ContentBlockType.text:
        return <ContentBlockText content={content as ContentUnitText} />;
      case ContentBlockType.infobox:
        return <ContentBlockInfobox content={content as ContentUnitInfobox} />;
      case ContentBlockType.table:
        return <ContentBlockTable content={content as ContentUnitTable} />;
      case ContentBlockType.image:
        return <ContentBlockFigure content={content as ContentUnitImage} />;
      case ContentBlockType.list:
        return <ContentBlockList content={content as ContentUnitList} />;
      case ContentBlockType.markdown:
        return <ContentBlockMarkdown content={content as ContentUnitMarkdown} />;
      case ContentBlockType.line:
        return <hr />;
      case ContentBlockType.code:
        return <ContentBlockCode content={content as ContentUnitCode} />;
      case ContentBlockType.entityTypeHeadline:
        return (
          <ContentBlockEntityTypeHeadline
            content={content as ContentUnitEntityTypeHeadline}
            ref={(el) => el && this.#headingEls.push(el)}
          />
        );
      case ContentBlockType.fieldDescription:
        return (
          <ContentBlockFieldDescription
            content={content as ContentUnitFieldDescription}
            ref={(el) => el && this.#headingEls.push(el)}
          />
        );
      case ContentBlockType.completenessDocumentationTable:
        return (
          <ContentBlockCompletenessDocumentationTable content={content as ContentUnitCompletenessDocumentationTable} />
        );
      default:
        console.error('Unsupported content block type: ' + type);
    }
  }

  componentShouldUpdate(_newValue, _oldValue, propName) {
    return propName !== 'scrollTarget';
  }

  componentWillLoad() {
    this.#intersectionObserver = new IntersectionObserver(
      (entries) => {
        entries.forEach(({ isIntersecting, target }) => {
          if (isIntersecting) {
            if (this.#hasPersistedDisplayedHeading) {
              this.#displayedHeadings = [`#${target.id}`];
              this.#hasPersistedDisplayedHeading = false;
            } else {
              this.#displayedHeadings.push(`#${target.id}`);
            }
          } else {
            if (this.#displayedHeadings.length > 1) {
              this.#displayedHeadings = this.#displayedHeadings.filter((slug) => slug !== `#${target.id}`);
            } else {
              this.#hasPersistedDisplayedHeading = true;
            }
          }
          this.handleViewPortIntersect(this.#displayedHeadings);
        });
      },
      { threshold: 0.5 }
    );
  }

  componentDidRender() {
    this.#intersectionObserver.disconnect();

    this.#headingEls.forEach((anchorEl) => {
      this.#intersectionObserver.observe(anchorEl);
    });
  }

  componentWillRender() {
    this.#headingEls = [];
    this.#displayedHeadings = [];
    this.#hasPersistedDisplayedHeading = false;
  }

  render() {
    const { blocks = [], alignment = '' } = this.content;

    if (this.#containerEl) {
      this.#containerEl.scrollTop = 0;
    }

    return (
      <Host ref={(el) => (this.#containerEl = el)} class={getPageClasses(['content', `content--${alignment}`])}>
        <div id="top-of-the-page" ref={(el) => this.#headingEls.push(el)} />
        {blocks.map((block) => this.getBlock(block))}
      </Host>
    );
  }
}

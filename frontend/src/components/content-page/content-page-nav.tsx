import { Component, h, Host, Prop, State, Watch } from '@stencil/core';
import { href } from 'stencil-router-v2';
import { TocTree } from 'services/content';
import stores from 'stores';
import { getPageClass, getPageClasses } from 'utils/content';

@Component({
  tag: 'mex-content-page-nav',
  styleUrl: 'content-page-nav.css',
})
export class ContentPageNavComponent {
  #isAnyAnchorActive: boolean = false;
  #navEl?: HTMLElement;
  #currentActiveAnchorEl?: HTMLAnchorElement;

  @State() isExpanded: boolean = true;

  @Prop() tocTree: TocTree;
  @Prop() activeAnchors?: string[] = [];
  @Prop() contentEl?: HTMLMexContentPageContentElement;
  @Prop() isMediaTablet?: boolean = true;
  @Prop() handleAnchorClick?: (slug: string) => void;

  @Watch('tocTree')
  handleTocTreeChange() {
    this.isExpanded = true;
  }

  ensureAnchorElVisibility(el: HTMLAnchorElement) {
    if (!el || el === this.#currentActiveAnchorEl || !this.#navEl) {
      return;
    }

    const { top: navTop, bottom: navBottom } = this.#navEl.getBoundingClientRect();
    const { top, bottom } = el.getBoundingClientRect();
    const safetyOffset = 50;
    const upperFoldOffset = top - navTop - safetyOffset;
    const lowerFoldOffset = navBottom - bottom - safetyOffset;

    upperFoldOffset < 0 && this.#navEl.scrollTo({ top: this.#navEl.scrollTop + upperFoldOffset, behavior: 'smooth' });
    lowerFoldOffset < 0 && this.#navEl.scrollTo({ top: this.#navEl.scrollTop - lowerFoldOffset, behavior: 'smooth' });

    this.#currentActiveAnchorEl = el;
  }

  getTocHTML(tree: TocTree, depth: number = 0) {
    const { isExpanded } = this;

    return (
      <div class={getPageClasses(['toc', `toc--depth-${depth}`])}>
        {tree.map((node) => {
          const { value, children } = node;
          const showActive = !this.#isAnyAnchorActive && this.activeAnchors.includes(value?.[1]);

          const url = value?.[1];
          const isAnchor = url?.[0] === '#';

          this.#isAnyAnchorActive = this.#isAnyAnchorActive || showActive;

          return (
            <ul class={getPageClass('toc-list')}>
              <li>
                {value && (
                  <a
                    class={getPageClasses([
                      'toc-link',
                      `toc-link--depth-${depth}`,
                      ...(showActive || (depth === 0 && !!children && !isExpanded) ? ['toc-link--active'] : []),
                    ])}
                    onClick={() => {
                      if (depth === 0 && children) {
                        this.isExpanded = !isExpanded;
                      } else {
                        isAnchor && this.handleAnchorClick?.(value[1]);
                      }
                    }}
                    {...(isAnchor ? (this.isMediaTablet ? {} : { href: url }) : href(`/${url}`))}
                    ref={(el) => showActive && this.ensureAnchorElVisibility(el)}
                  >
                    <span class={getPageClass('toc-link-text')}>{value[0]}</span>
                    {depth === 0 && (
                      <mex-icon-chevron
                        classes={`icon--inline icon--mirrorable ${
                          !!children && isExpanded ? 'icon--mirrored-vertical' : ''
                        }`}
                      />
                    )}
                  </a>
                )}

                {!!children && (
                  <mex-accordion expanded={depth !== 0 || isExpanded}>
                    {this.getTocHTML(children, depth + 1)}
                  </mex-accordion>
                )}
              </li>
            </ul>
          );
        })}
      </div>
    );
  }

  render() {
    this.#isAnyAnchorActive = false;
    const rootNode = this.tocTree.find(({ children }) => children);
    rootNode && rootNode?.value?.[1] !== '#' && (rootNode.value[1] = '#top-of-the-page');

    return (
      <Host>
        {!!this.tocTree.length && (
          <nav class={getPageClass('nav')} ref={(el) => (this.#navEl = el)}>
            <div class={`u-underline-2 ${getPageClass('nav-header')}`}>{stores.i18n.t('content.navigation.title')}</div>
            {this.getTocHTML(this.tocTree)}
          </nav>
        )}
      </Host>
    );
  }
}

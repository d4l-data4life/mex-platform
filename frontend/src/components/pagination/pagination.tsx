import { Component, h, Prop } from '@stencil/core';
import { SEARCH_PAGINATION_START } from 'config';
import stores from 'stores';

@Component({
  tag: 'mex-pagination',
  styleUrl: 'pagination.css',
})
export class PaginationComponent {
  @Prop() range: number[];
  @Prop() current?: number;
  @Prop() disabled = false;
  @Prop() renderSlot = false;
  @Prop() handleClick?: (num: number) => void;
  @Prop() ariaLabelAttr?: string;
  @Prop() testAttr?: string;

  get index() {
    return this.range.indexOf(this.current);
  }

  get hasPrev() {
    return this.index > 0;
  }

  get hasNext() {
    const { index } = this;
    return index !== -1 && index < this.range.length - 1;
  }

  render() {
    const { testAttr, current, range } = this;

    return (
      <nav
        class="pagination"
        role="navigation"
        aria-label={this.ariaLabelAttr ?? stores.i18n.t('navigation.label')}
        data-test={testAttr}
      >
        <button
          aria-label={`${stores.i18n.t('navigation.back')}. ${stores.i18n.t('item.navigation.resultIndexAria', {
            count: range.length,
            index: current + SEARCH_PAGINATION_START,
          })}`}
          class="pagination__item"
          onClick={() => this.handleClick?.(current - 1)}
          disabled={this.disabled || !this.hasPrev}
          title={stores.i18n.t('navigation.back')}
          data-test={testAttr && `${testAttr}:button`}
          data-test-key={testAttr && 'prev'}
        >
          <mex-icon-arrow classes="icon--mirrored-horizontal" />
        </button>
        {!this.renderSlot &&
          this.range.map((num) => (
            <button
              key={num}
              class={{ pagination__item: true, 'pagination__item--active': current === num }}
              onClick={() => this.handleClick?.(num)}
              disabled={this.disabled}
              aria-current={current === num ? 'page' : null}
              data-test={testAttr && `${testAttr}:button`}
              data-test-key={testAttr && num}
              data-test-active={testAttr && String(current === num)}
            >
              {num}
            </button>
          ))}
        {this.renderSlot && (
          <div class="pagination__item pagination__info">
            <slot />
          </div>
        )}
        <button
          aria-label={`${stores.i18n.t('navigation.forward')}. ${stores.i18n.t('item.navigation.resultIndexAria', {
            count: range.length,
            index: current + SEARCH_PAGINATION_START,
          })}`}
          class="pagination__item"
          onClick={() => this.handleClick?.(current + 1)}
          disabled={this.disabled || !this.hasNext}
          title={stores.i18n.t('navigation.forward')}
          data-test={testAttr && `${testAttr}:button`}
          data-test-key={testAttr && 'next'}
        >
          <mex-icon-arrow />
        </button>
      </nav>
    );
  }
}

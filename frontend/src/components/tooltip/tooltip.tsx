import { Component, h, Host, Prop, State } from '@stencil/core';
import stores from 'stores';

let overallTooltipCount: number = 0;

@Component({
  tag: 'mex-tooltip',
  styleUrl: 'tooltip.css',
})
export class TooltipComponent {
  #id: string;

  @Prop() text?: string;
  @Prop() classes = '';
  @Prop() method: 'focus' | 'hover' = 'focus';
  @Prop() handleClick?: (event?: MouseEvent | TouchEvent) => void;
  @Prop() testAttr?: string;

  @State() verticalOrientation: 'up' | 'down' = 'up';
  @State() horizontalOrientation: 'left' | 'center' | 'right' = 'center';
  @State() maxWidth: number;
  @State() isVisible = false;

  get isClickable() {
    return this.method === 'focus' || !!this.handleClick;
  }

  alignTooltip(event: FocusEvent | MouseEvent | TouchEvent | KeyboardEvent) {
    const { left, top, width, height } = (event.target as HTMLButtonElement).getBoundingClientRect();
    const { innerWidth: vWidth, innerHeight: vHeight } = window;
    const x = left + width / 2;
    const y = top + height / 2;
    this.horizontalOrientation = x / vWidth < 1 / 3 ? 'right' : x / vWidth > 2 / 3 ? 'left' : 'center';
    this.verticalOrientation = y / vHeight <= 0.5 ? 'down' : 'up';
    this.maxWidth = this.getMaxWidth(left);
  }

  getMaxWidth(left: number) {
    const { horizontalOrientation } = this;
    const total = window.innerWidth;
    const leftPerc = left / total;
    const spaceLeft = Math.floor(leftPerc * 100) - 5;
    const spaceRight = Math.floor((1 - leftPerc) * 100) - 5;

    if (horizontalOrientation === 'left') {
      return spaceLeft;
    }

    if (horizontalOrientation === 'right') {
      return spaceRight;
    }

    return Math.min(spaceLeft, spaceRight) * 2;
  }

  show(event: FocusEvent | MouseEvent | TouchEvent | KeyboardEvent) {
    this.alignTooltip(event);
    this.isVisible = true;
    this.preventTextSelect();
  }

  hide() {
    this.isVisible = false;
    this.preventTextSelect();
  }

  toggle(event: KeyboardEvent) {
    this.isVisible ? this.hide() : this.show(event);
  }

  onKeyUp(event: KeyboardEvent) {
    const { key } = event;
    [' ', 'Enter'].includes(key) && this.toggle(event);
    key === 'Escape' && this.hide();
  }

  preventTextSelect() {
    if (this.method === 'focus') {
      document.addEventListener('selectstart', this.handleSelectStart);
      window.setTimeout(() => document.removeEventListener('selectstart', this.handleSelectStart), 200);
    }
  }

  handleSelectStart(event: Event) {
    event.preventDefault();
  }

  componentWillLoad() {
    overallTooltipCount++;
    this.#id = `tooltip-${overallTooltipCount}`;
  }

  render() {
    const id = this.#id;
    const { method, isVisible, testAttr } = this;

    return (
      <Host
        class={`tooltip tooltip--method-${method} ${this.isClickable ? 'tooltip--clickable' : ''} ${
          isVisible ? 'tooltip--visible' : ''
        } tooltip--${this.verticalOrientation}  tooltip--${this.horizontalOrientation} ${this.classes}`}
        data-test={testAttr}
        data-test-active={testAttr && String(isVisible)}
      >
        {isVisible && (
          <div
            class="tooltip__backdrop"
            onClick={() => this.hide()}
            onMouseEnter={method === 'hover' ? () => this.hide() : null}
            onTouchStart={method === 'hover' ? () => this.hide() : null}
          />
        )}

        <button
          class="tooltip__toggle"
          onFocus={method === 'focus' ? (event) => this.show(event) : null}
          onBlur={method === 'focus' ? () => this.hide() : null}
          onKeyUp={method === 'focus' ? (event) => this.onKeyUp(event) : null}
          onClick={(event) => {
            method === 'focus' && this.show(event);
            this.handleClick?.(event);
          }}
          onMouseEnter={method === 'hover' ? (event) => this.show(event) : null}
          onMouseLeave={method === 'hover' ? () => this.hide() : null}
          onTouchStartCapture={method === 'hover' ? (event) => this.show(event) : null}
          onContextMenu={(event) => event.preventDefault()}
          tabIndex={method === 'focus' ? 0 : -1}
          aria-hidden={String(method === 'hover')}
          aria-label={stores.i18n.t('tooltip.label')}
          aria-describedby={id}
          aria-expanded={String(this.isVisible)}
          data-test={testAttr && `${testAttr}:button`}
        >
          <slot name="toggle">
            <mex-icon-info classes="icon--medium" />
          </slot>
        </button>
        <div class="tooltip__container" style={{ '--tooltip-width': `${this.maxWidth ?? 0}vw` }}>
          <div class="tooltip__content" id={id} role="tooltip" aria-hidden={String(!isVisible)}>
            <slot name="content">
              <span class="tooltip__text">{this.text}</span>
            </slot>
          </div>
        </div>
      </Host>
    );
  }
}

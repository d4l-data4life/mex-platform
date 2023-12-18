import { Component, Fragment, h, Prop } from '@stencil/core';
import { href } from 'stencil-router-v2';

@Component({
  tag: 'mex-tile',
  styleUrl: 'tile.css',
})
export class TileComponent {
  @Prop() text?: string;
  @Prop() hint?: string;
  @Prop() url?: string;
  @Prop() isBusy?: boolean;
  @Prop() testAttr?: string;
  @Prop() testKeyAttr?: number;
  @Prop() icon?: 'arrow' | 'drilldown';
  @Prop() modifiers: string[] = [];
  @Prop() disabled = false;
  @Prop() handleClick?: () => void;

  get linkAttributes() {
    const isExternal = ['http:', 'https:', 'mailto:'].some((protocol) => this.url?.indexOf(protocol) === 0);
    const attrs = isExternal ? { href: this.url } : href(this.url);

    return {
      ...attrs,
      onClick: (event: MouseEvent) => {
        this.handleClick?.();
        (attrs as any).onClick?.(event);
      },
    };
  }

  get contents() {
    const { text, hint, isBusy, icon } = this;
    const iconProps = { class: 'tile__icon', classes: 'icon--medium' };

    return (
      <Fragment>
        <div class="tile__text">
          <span>
            {text ?? ''}
            {isBusy && <mex-placeholder lines={2} />}
          </span>
        </div>
        <div class="tile__footer">
          {!!hint && <span class="tile__hint">{hint}</span>}
          {icon === 'arrow' && <mex-icon-arrow {...iconProps} />}
          {icon === 'drilldown' && <mex-icon-arrow-drilldown {...iconProps} />}
        </div>
      </Fragment>
    );
  }

  render() {
    const { contents, testAttr, testKeyAttr, handleClick, disabled, modifiers } = this;

    const classes = {
      tile: true,
      'tile--link': !!(this.url || this.handleClick),
      ...modifiers.reduce((obj, modifier) => Object.assign(obj, { [`tile--${modifier}`]: true }), {}),
    };

    return this.url ? (
      <a class={classes} {...this.linkAttributes} data-test={testAttr} data-test-key={testKeyAttr}>
        {contents}
      </a>
    ) : handleClick ? (
      <button
        class={classes}
        onClick={!disabled && handleClick}
        data-test={testAttr}
        data-test-key={testKeyAttr}
        disabled={disabled}
      >
        {contents}
      </button>
    ) : (
      <div class={classes} data-test={testAttr} data-test-key={testKeyAttr}>
        {contents}
      </div>
    );
  }
}

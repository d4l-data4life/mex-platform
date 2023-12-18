import { Component, h, Fragment, Prop } from '@stencil/core';
import stores from 'stores';

@Component({
  tag: 'mex-tag',
  styleUrl: 'tag.css',
  scoped: true,
})
export class TagComponent {
  @Prop() text = '';
  @Prop() closable = true;
  @Prop() clickable = false;
  @Prop() closeTitle?: string;
  @Prop() testAttr?: string;

  @Prop() handleClose?: () => void;
  @Prop() handleClick?: () => void;

  get classes() {
    const { clickable, closable } = this;
    return ['tag', clickable && 'tag--clickable', closable && 'tag--closable'].filter(Boolean).join(' ');
  }

  render() {
    const { clickable, closable, text, classes, testAttr } = this;

    return (
      <Fragment>
        {clickable ? (
          <button class={classes} onClick={() => this.handleClick?.()} data-test={testAttr && `${testAttr}:button`}>
            {text}
          </button>
        ) : (
          <span class={classes} data-test={testAttr && `${testAttr}:text`}>
            {text}
          </span>
        )}
        {closable && (
          <button
            class="tag__close"
            onClick={() => this.handleClose?.()}
            title={this.closeTitle ?? stores.i18n.t('tag.close')}
            data-test={testAttr && `${testAttr}:close`}
          >
            <mex-icon-close classes="icon--small" />
          </button>
        )}
      </Fragment>
    );
  }
}

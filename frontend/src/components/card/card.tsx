import { Component, Fragment, h, Prop } from '@stencil/core';
import { href } from 'stencil-router-v2';

export interface CardRow {
  type: string;
  label?: string;
  value?: string;
  testAttr?: string;
  onClick?: (event?: MouseEvent) => void;
}

@Component({
  tag: 'mex-card',
  styleUrl: 'card.css',
})
export class CardComponent {
  @Prop() rows: CardRow[];
  @Prop() icon?: 'email' | 'support';

  render() {
    const { icon, rows } = this;

    return (
      <div class="card">
        {icon === 'email' && <mex-icon-email class="card__icon u-underline-3" classes="icon--large" />}
        {icon === 'support' && <mex-icon-support class="card__icon u-underline-3" classes="icon--large" />}

        {rows.map(({ type, label, value, testAttr, onClick }, index) => (
          <div class={`card__row card__row--${type}`} key={index} data-test={type === 'value' && testAttr}>
            {type === 'value' && (
              <Fragment>
                <div class="card__label">{label}</div>
                <div class="card__value">{value}</div>
              </Fragment>
            )}
            {type === 'text' && <div class="card__text">{value}</div>}
            {type === 'action' && value?.[0] === '/' && (
              <a class="button" {...href(value)} data-test={testAttr} onClick={onClick}>
                {label}
              </a>
            )}
            {type === 'action' && value?.startsWith('mailto:') && (
              <a class="button" href={value} data-test={testAttr} onClick={onClick}>
                {label}
              </a>
            )}
            {type === 'action' && !value && !!onClick && (
              <button class="button" data-test={testAttr} onClick={onClick}>
                {label}
              </button>
            )}
          </div>
        ))}
      </div>
    );
  }
}

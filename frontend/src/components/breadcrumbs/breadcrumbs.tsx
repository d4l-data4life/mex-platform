import { Component, Fragment, h, Prop } from '@stencil/core';
import { href } from 'stencil-router-v2';

export interface Breadcrumb {
  label: string;
  url?: string;
  icon?: string;
  testAttr?: string;
  testKeyAttr?: number;
  handleClick?: () => void;
}

@Component({
  tag: 'mex-breadcrumbs',
  styleUrl: 'breadcrumbs.css',
})
export class BreadcrumbNavigationComponent {
  @Prop() items: Breadcrumb[];

  getContent({ label, icon }: Breadcrumb) {
    const iconClasses = 'breadcrumbs__item-icon icon--inline icon--no-vertical-offset';
    const iconClassesEntity = `${iconClasses} icon--large u-underline-3`;

    return (
      <Fragment>
        {icon === 'back' && <mex-icon-arrow classes={`${iconClasses} icon--medium icon--mirrored-horizontal`} />}
        {icon?.includes('entity:') && (
          <mex-icon-entity entityName={icon.split(':').pop()} attrs={{ classes: iconClassesEntity }} />
        )}
        <span class="breadcrumbs__item-label">{label}</span>
      </Fragment>
    );
  }

  render() {
    const { items, getContent } = this;

    return (
      <ol class="breadcrumbs">
        {items.map((item, index) => (
          <li key={index} class="breadcrumbs__item" data-test={item.testAttr} data-test-key={item.testKeyAttr}>
            {item.url && <a {...href(item.url)}>{getContent(item)}</a>}
            {item.handleClick && <button onClick={() => item.handleClick()}>{getContent(item)}</button>}
            {!item.url && !item.handleClick && <span>{getContent(item)}</span>}
          </li>
        ))}
      </ol>
    );
  }
}

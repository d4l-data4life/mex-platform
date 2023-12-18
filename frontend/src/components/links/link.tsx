import { Component, Fragment, h, Prop } from '@stencil/core';
import { href } from 'stencil-router-v2';

@Component({
  tag: 'mex-link',
})
export class LinkComponent {
  @Prop() url: string;
  @Prop() label: string;
  @Prop() showUrl?: boolean;
  @Prop() testAttr?: string;

  get linkLabel() {
    return this.showUrl ? this.url : this.label;
  }
  get isExternal(): boolean {
    const { url } = this;
    return !url.startsWith('/') && !url.startsWith(document.location.origin);
  }
  render() {
    const { label, linkLabel, isExternal, showUrl, url, testAttr } = this;
    return (
      <Fragment>
        {showUrl && (
          <Fragment>
            {label}
            <br />
          </Fragment>
        )}
        {isExternal && (
          <a href={url} target="_blank" rel="noopener noreferrer" data-test={testAttr}>
            {linkLabel}
          </a>
        )}
        {!isExternal && (
          <a {...href(url)} data-test={testAttr}>
            {linkLabel}
          </a>
        )}
      </Fragment>
    );
  }
}

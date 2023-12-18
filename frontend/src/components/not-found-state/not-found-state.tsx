import { Component, h, Host, Prop } from '@stencil/core';
import { href } from 'stencil-router-v2';

@Component({
  tag: 'mex-not-found-state',
  styleUrl: 'not-found-state.css',
})
export class NotFoundStateComponent {
  @Prop() caption: string;
  @Prop() text?: string;
  @Prop() buttonText?: string;
  @Prop() buttonUrl?: string;
  @Prop() testAttr?: string;

  render() {
    const { caption, text, buttonText, buttonUrl, testAttr } = this;

    return (
      <Host class="not-found-state" data-test={testAttr}>
        <mex-illustration-search-empty class="not-found-state__illustration" />

        <h4 class="not-found-state__title">{caption}</h4>
        {text && <p class="not-found-state__text">{text}</p>}

        {buttonText && buttonUrl && (
          <a class="button" {...href(buttonUrl)} data-test={testAttr && `${testAttr}:button`}>
            {buttonText}
          </a>
        )}
      </Host>
    );
  }
}

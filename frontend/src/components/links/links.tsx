import { Component, h, Prop } from '@stencil/core';
import { LinkItem } from 'config/navigation';

@Component({
  tag: 'mex-links',
})
export class LinksComponent {
  @Prop() items: LinkItem[];
  @Prop() classes?: string;
  @Prop() itemClasses?: string;
  @Prop() showUrl?: boolean;

  render() {
    const { classes, itemClasses, items, showUrl } = this;

    return (
      <ul class={classes}>
        {items.map((item, index) => (
          <li key={index} class={itemClasses}>
            <mex-link {...item} showUrl={showUrl} />
          </li>
        ))}
      </ul>
    );
  }
}

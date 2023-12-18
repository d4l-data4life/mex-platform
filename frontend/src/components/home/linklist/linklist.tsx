import { Component, h, Prop } from '@stencil/core';
import { LinkItem } from 'config/navigation';

export interface LinklistTile {
  items: LinkItem[];
}

@Component({
  tag: 'mex-linklist',
  styleUrl: 'linklist.css',
})
export class LinklistComponent {
  @Prop() headline: string;
  @Prop() tile: LinklistTile;

  render() {
    const { headline, tile } = this;

    return (
      <div class="linklist">
        <h2 class="linklist__title u-underline-2">{headline}</h2>
        <mex-links items={tile.items} classes="linklist__list" itemClasses="linklist__item" />
      </div>
    );
  }
}

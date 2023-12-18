import { Component, h, Host, Prop } from '@stencil/core';

export interface TabsItem {
  label: string;
  isActive: boolean;
  data?: unknown;
}

@Component({
  tag: 'mex-tabs',
  styleUrl: 'tabs.css',
})
export class TabsComponent {
  @Prop() items: TabsItem[];
  @Prop() handleClick?: (item: TabsItem) => void;
  @Prop() testAttr?: string;

  render() {
    return (
      <Host class="tabs">
        {this.items.map((item, index) => (
          <button
            class={{ tabs__item: true, 'tabs__item--active': item.isActive }}
            key={index}
            onClick={() => this.handleClick?.(item)}
            role="tab"
            aria-selected={String(item.isActive)}
            data-test={this.testAttr}
            data-test-active={this.testAttr && String(item.isActive)}
          >
            {item.label}
          </button>
        ))}
      </Host>
    );
  }
}

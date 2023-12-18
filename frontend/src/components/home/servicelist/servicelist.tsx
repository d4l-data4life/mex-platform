import { Component, h, Host, Prop } from '@stencil/core';
import { LinkItem } from 'config/navigation';
import stores from 'stores';

@Component({
  tag: 'mex-servicelist',
  styleUrl: 'servicelist.css',
})
export class ServiceListComponent {
  @Prop() items: LinkItem[];

  render() {
    return (
      <Host class="servicelist">
        <nav>
          <h4 class="servicelist__title u-underline-2">{stores.i18n.t('navigation.services')}</h4>
          <mex-links items={this.items} classes="servicelist__list" itemClasses="servicelist__link" showUrl />
        </nav>
      </Host>
    );
  }
}

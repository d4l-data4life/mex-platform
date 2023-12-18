import { Component, h, Host, State } from '@stencil/core';
import { href } from 'stencil-router-v2';
import stores from 'stores';
import { ROUTES } from 'config';

@Component({
  tag: 'mex-header-user',
  styleUrl: 'header-user.css',
})
export class HeaderUserComponent {
  #containerRef: HTMLMexDropdownElement;
  #userEmail: string;

  @State() isExpanded: boolean;

  componentWillLoad() {
    this.#userEmail = stores.auth.userEmail;
  }

  render() {
    return (
      <Host class="header-user">
        <mex-dropdown
          ref={(el) => (this.#containerRef = el)}
          class="header-user__toggle"
          toggleClass="header-user__select"
          orientation="right"
          handleExpand={() => (this.isExpanded = true)}
          handleCollapse={() => (this.isExpanded = false)}
          testAttr="header:user"
        >
          <div slot="label">
            <mex-icon-user classes="icon--large header-user__icon" />
          </div>
          <ul class="header-user__menu" slot="options">
            {this.#userEmail && (
              <li class="header-user__item">
                <span>{this.#userEmail}</span>
              </li>
            )}
            <li class="header-user__item">
              <a
                {...href(ROUTES.LOGOUT)}
                tabIndex={this.isExpanded ? 0 : -1}
                onBlur={() => {
                  requestAnimationFrame(() => {
                    !this.#containerRef?.contains(document.activeElement) && this.#containerRef.collapse();
                  });
                }}
              >
                {stores.i18n.t('navigation.logout')}
              </a>
            </li>
          </ul>
        </mex-dropdown>
      </Host>
    );
  }
}

import { Component, h } from '@stencil/core';
import stores from 'stores';

@Component({
  tag: 'mex-notifications',
  styleUrl: 'notifications.css',
})
export class NotificationsComponent {
  render() {
    const { items } = stores.notifications;

    return (
      !!items.length && (
        <div class="notifications" role="status" aria-live="assertive">
          {items.map((key) => (
            <div class={`notifications__item notifications__item--${key.split(':')[1] ?? 'info'}`} key={key}>
              {stores.i18n.t(`notifications.${key.split(':')[0]}`)}
            </div>
          ))}
        </div>
      )
    );
  }
}

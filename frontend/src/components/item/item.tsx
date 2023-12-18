import { Component, Fragment, h, Prop } from '@stencil/core';
import { ITEM_CONFIG } from 'config';
import { SidebarFeature } from 'config/item';
import { Item } from 'services/item';
import { aggregateSidebarFeatures } from 'utils/config';

@Component({
  tag: 'mex-item',
  styleUrl: 'item.css',
})
export class ItemComponent {
  @Prop() item?: Item;
  @Prop() previousVersionItem?: Item;
  @Prop() latestVersionItem?: Item;
  @Prop() highlightChanges = false;

  get sidebarFeatures() {
    return ITEM_CONFIG.SIDEBAR_FEATURES[this.item?.entityType] ?? ITEM_CONFIG.SIDEBAR_FEATURES.default;
  }

  get aggregatedSidebarFeatures() {
    return aggregateSidebarFeatures(this.sidebarFeatures);
  }

  render() {
    const { item, latestVersionItem, previousVersionItem, highlightChanges, aggregatedSidebarFeatures } = this;

    return (
      <Fragment>
        <div class="item" data-test="item">
          <mex-item-fields
            class="item__fields"
            data-test="item:details"
            item={item}
            previousVersionItem={previousVersionItem}
            highlightChanges={highlightChanges}
            context="details"
          />
          {item !== null && (
            <aside class="item__sidebar" data-test="item:sidebar">
              {aggregatedSidebarFeatures.map(({ feature, configs }) => (
                <Fragment>
                  {feature === 'itemInfo' && item && (
                    <mex-item-info
                      class="item__sidebar-feature item__info"
                      item={item}
                      orientation="auto"
                      features={configs}
                    />
                  )}
                  {feature === SidebarFeature.contactForm &&
                    configs.map((config) => (
                      <mex-item-contact
                        class="item__sidebar-feature"
                        item={latestVersionItem}
                        isOnLatestVersion={item?.itemId === latestVersionItem?.itemId}
                        config={config}
                      />
                    ))}
                  {feature === SidebarFeature.displayField && (
                    <mex-item-fields
                      class="item__sidebar-feature"
                      data-test="item:sidebar:fields"
                      item={item}
                      fields={configs.map(({ field }) => field)}
                      previousVersionItem={previousVersionItem}
                      highlightChanges={highlightChanges}
                      context="sidebar"
                    />
                  )}
                </Fragment>
              ))}
            </aside>
          )}
        </div>
        {item && <mex-item-related-results class="item__related-results" item={item} />}
      </Fragment>
    );
  }
}

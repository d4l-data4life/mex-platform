import { Component, Prop, h, Fragment } from '@stencil/core';
import { ENTITY_TYPES } from 'config';
import { EntityTypeIcon } from 'config/entity-types';

@Component({
  tag: 'mex-icon-entity',
})
export class IconEntityComponent {
  @Prop() entityName: string;
  @Prop() attrs: object = {};

  render() {
    const { entityName, attrs } = this;
    const icon = ENTITY_TYPES[entityName].config?.icon;

    return (
      !!icon && (
        <Fragment>
          {icon === EntityTypeIcon.source && <mex-icon-source {...attrs} />}
          {icon === EntityTypeIcon.resource && <mex-icon-resource {...attrs} />}
          {icon === EntityTypeIcon.datum && <mex-icon-datum {...attrs} />}
          {icon === EntityTypeIcon.platform && <mex-icon-platform {...attrs} />}
        </Fragment>
      )
    );
  }
}

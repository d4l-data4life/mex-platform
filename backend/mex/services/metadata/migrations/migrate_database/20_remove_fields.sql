ALTER TABLE "item_values" DROP CONSTRAINT "fk_field_defs_values";

DROP VIEW "items_with_links";

CREATE OR REPLACE FUNCTION "items_with_links"(linkFieldNames TEXT[])
RETURNS TABLE(source_item_id TEXT, source_created_at timestamptz(6), source_owner TEXT, source_entity_name TEXT, source_business_id_field_name TEXT, source_business_id TEXT, source_link_field_name TEXT, target_business_id TEXT, target_item_id TEXT)
LANGUAGE plpgsql AS
$func$
BEGIN
    RETURN QUERY
select z.source_item_id, z.source_created_at, z.source_owner, z.source_entity_name, z.source_business_id_field_name, z.source_business_id, z.source_link_field_name, z.target_business_id, z.target_item_id
from (
	select y.*, row_number() over (partition by y.source_business_id, y."target_business_id", y."source_link_field_name" order by y.source_created_at desc) as nr
	from (
		select x.*, iwbi_target.item_id as target_item_id
		from (
		    select iwbi_sources.item_id as source_item_id, iwbi_sources.created_at as source_created_at, iwbi_sources.owner as source_owner, iwbi_sources.entity_name as source_entity_name, iwbi_sources.business_id_field_name as source_business_id_field_name, iwbi_sources.business_id as source_business_id, civ.field_name as "source_link_field_name", civ.field_value as "target_business_id"
		    from current_item_values civ
		    inner join items_with_business_id iwbi_sources
		      on iwbi_sources.item_id = civ.item_id
		    where civ.field_name = ANY ($1)
		) x
		left outer join items_with_business_id iwbi_target
		  on iwbi_target.business_id = x.target_business_id
	) y
) z
where z.nr = 1;
END;
$func$;

DROP TABLE "field_defs";

CREATE OR REPLACE FUNCTION next_migration_version() RETURNS integer
LANGUAGE plpgsql IMMUTABLE AS
$$
BEGIN
    return 21;
END;
$$;

CREATE OR REPLACE VIEW "search_config" AS (
    SELECT sco."id", sco."name", sco."obj_type", scf."field_name"
    FROM "search_config_objects" sco
    INNER JOIN "search_config_fields" scf
    ON sco."id" = scf."object_id"
);

CREATE OR REPLACE VIEW "current_item_values" AS (
	SELECT x."id", x."item_id", x."field_name", x."field_value", x."place", x."revision", x."language"
    FROM "item_values" x
	WHERE x."deleted" = false
);

CREATE OR REPLACE VIEW "items_with_business_id" AS (
    SELECT i."id" as "item_id", i."created_at", i."owner", i."entity_name", coalesce(i."business_id", '') as "business_id", i."business_id_field_name"
    FROM "items" i
    WHERE i."business_id" IS NOT NULL
);

CREATE OR REPLACE VIEW "latest_items_with_business_id" AS (
    select x."item_id", x."created_at" , x."owner", x."entity_name", x."business_id", x."business_id_field_name"
    from (
        select *, row_number() over (partition by "business_id", "entity_name" order by "created_at" desc) as "nr"
        from "items_with_business_id"
    ) x
    where x."nr" = 1
);

CREATE OR REPLACE VIEW "items_nullable_business_id" AS (
    SELECT i."id" as "item_id", i."created_at", i."owner", i."entity_name", i."business_id", i."business_id_field_name"
    FROM "items" i
);

CREATE OR REPLACE VIEW "items_with_links" AS (
select z.source_item_id, z.source_created_at, z.source_owner, z.source_entity_name, z.source_business_id_field_name, z.source_business_id, z.source_link_field_name, z.target_business_id, z.target_item_id
from (
	select y.*, row_number() over (partition by source_business_id, "target_business_id", "source_link_field_name" order by source_created_at desc) as nr
	from (
		select x.*, iwbi_target.item_id as target_item_id
		from (
		    select iwbi_sources.item_id as source_item_id, iwbi_sources.created_at as source_created_at, iwbi_sources.owner as source_owner, iwbi_sources.entity_name as source_entity_name, iwbi_sources.business_id_field_name as source_business_id_field_name, iwbi_sources.business_id as source_business_id, civ.field_name as "source_link_field_name", civ.field_value as "target_business_id"
		    from current_item_values civ
		    inner join items_with_business_id iwbi_sources
		      on iwbi_sources.item_id = civ.item_id
		    --- Since the table field_defs does not contain linked fields, the below line only works if no linked fields are of type 'link'
		    where civ.field_name in (select fd."name" from field_defs fd where fd.kind = 'link')
		) x
		left outer join items_with_business_id iwbi_target
		  on iwbi_target.business_id = x.target_business_id
	) y
) z
where z.nr = 1
);

CREATE OR REPLACE FUNCTION "f_latest_items_with_business_id_bounded"(horizon timestamptz(6))
returns table (item_id text, created_at timestamptz(6), owner text, entity_name text, business_id text)
language plpgsql as
$func$
begin
    return query
    select x."item_id", x."created_at" , x."owner", x."entity_name", x."business_id"
    from (
        select z.*, row_number() over (partition by z."business_id", z."entity_name" order by z."created_at" desc) as "nr"
        from "items_with_business_id" z
        where z.created_at <= horizon
    ) x
    where x."nr" = 1;
end;
$func$;

CREATE OR REPLACE FUNCTION next_migration_version() RETURNS integer
LANGUAGE plpgsql IMMUTABLE AS
$$
BEGIN
    return 17;
END;
$$;


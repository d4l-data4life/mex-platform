-- name: DbGetItemsCount :one
SELECT COUNT(*) FROM items;

-- name: DbGetItem :many
SELECT * FROM items
WHERE id = $1;

-- name: DbGetItemValues :many
SELECT * FROM current_item_values
WHERE item_id = $1
ORDER BY field_name ASC, place ASC;

-- name: DbListItemValues :many
SELECT * FROM current_item_values
ORDER BY item_id ASC, field_name ASC, place ASC;

-- name: DbListItems :many
SELECT * FROM items_nullable_business_id
WHERE item_id >= $1 ORDER BY item_id ASC LIMIT $2;

-- name: DbListItemsOfType :many
SELECT * FROM items_nullable_business_id
WHERE item_id >= $1 AND entity_name = $2 ORDER BY item_id ASC LIMIT $3;

-- name: DbCreateItem :one
INSERT INTO items (created_at, id, owner, entity_name, business_id_field_name, business_id, hash)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DbDeleteItem :exec
DELETE FROM items where id = $1;

-- name: DbDeleteItems :execrows
DELETE FROM items where id = ANY(@item_ids::text[]);

-- name: DbDeleteAllItems :exec
DELETE FROM items;

-- name: DbCreateItemValue :one
INSERT INTO item_values (created_at, id, revision, deleted, field_name, field_value, language, place, item_id)
VALUES ($1, $2, DEFAULT, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- Perform an UPSERT:
-- - Either INSERT a new row with initial count 1, or, if already present
-- - UPDATE the row by incrementing the count
-- name: DbIncrementItemCounter :exec
INSERT INTO item_view_counts (created_at, item_id, user_id, counts)
VALUES ($1, $2, $3, $4)
ON CONFLICT (item_id, user_id) DO UPDATE
    SET counts = item_view_counts.counts + EXCLUDED.counts;

-- name: DbCreateRelation :one
INSERT INTO relations (created_at, id, owner, deleted, type, source_item_id, target_item_id, info_item_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: DbGetRelation :one
SELECT * FROM relations
WHERE id = $1
LIMIT 1;

-- name: DbGetRelationTargetsForSourceAndType :many
SELECT target_item_id FROM relations
WHERE source_item_id = $1 and type = $2;

-- name: DbListRelations :many
SELECT * FROM relations;

-- name: DbCreateRelationsFromBusinessIDs :execrows
INSERT INTO "relations" ("created_at", "id", "deleted", "owner", "type", "source_item_id", "target_item_id", "info_item_id")
SELECT
	NOW()              AS "created_at",
	uuid_generate_v4() AS "id",
	false              AS "deleted",
	$1                 AS "owner",
	$2                 AS "type",
	"source_item_id"   AS "source_item_id",
	"target_item_id"   AS "target_item_id",
	null               AS "info_item_id"
FROM (
	select civ.item_id as source_item_id, liwbi.item_id as target_item_id
	from current_item_values civ
	inner join latest_items_with_business_id liwbi
		on liwbi.business_id = civ.field_value and civ.field_name = $4
	where civ.item_id = $3 and liwbi.item_id <> $3
) unused1;

-- name: DbComputeVersions :many
select y.item_id, y.created_at, y.version from (
	select x.*, row_number() over () as version
	from (
		select iwbi2.*
		from items_with_business_id iwbi1
		inner join items_with_business_id iwbi2
			on iwbi1.business_id = iwbi2.business_id
		where iwbi1.item_id = $1
		order by iwbi2.created_at
	) x
) y
order by y.version desc;

-- name: DbAggregationCandidateValues :many
select civ.field_name, y.partition_business_id, civ.place, civ.field_value, civ.language, civ.item_id
from (select distinct row_number()
                      over (partition by x.partition_business_id order by x.candidate_created_at desc) as date_rank,
                      x.candidate_id                                                                   as candidate_id,
                      x.partition_business_id                                                          as partition_business_id
      from (select i.id as candidate_id, i.created_at as candidate_created_at, iv.field_value as partition_business_id
            from items i
                     join item_values iv
                          on i.id = iv.item_id and iv.field_name = sqlc.arg(partition_field)
            where i.business_id = sqlc.arg(business_id)) x) y
         join current_item_values civ
              on y.candidate_id = civ.item_id and civ.field_name <> sqlc.arg(source_id_field)
where y.date_rank = 1
order by civ.field_name asc, civ.item_id asc, y.partition_business_id asc, civ.place asc, civ.field_value asc;

-- name: DbListItemsWithBusinessId :many
select * from items_with_business_id order by business_id, created_at asc;

-- name: DbListItemsForBusinessId :many
select * from items_with_business_id where business_id = $1 order by created_at asc;

-- name: DbFollowRelationsTwoSteps :many
select distinct r1."type" as relation_type_1 , r2.type as relation_type_2, r2.target_item_id
from relations r2
inner join relations r1
  on r2.source_item_id = r1.target_item_id
where r1.source_item_id = $1;

-- name: DbLatestItemValuesByBusinessId :many
select civ.*
from latest_items_with_business_id liwbi
inner join current_item_values civ
  on civ.item_id = liwbi.item_id
where liwbi.business_id = $1;

-- name: DbLatestItemByBusinessId :one
select i.id as item_id, i."owner", i.entity_name , i.business_id_field_name, i.business_id, TO_CHAR(i.created_at, 'YYYY-MM-DD"T"HH24:MI:SSZ') as created_at
from items i
inner join latest_items_with_business_id liwbi
	on i.id =liwbi.item_id
where liwbi.business_id = $1;

-- name: DbImputeBusinessId :execrows
update items
set business_id = val.business_id
from
	(SELECT i."id" as "item_id", iv."field_value" AS "business_id"
    FROM "items" i
    INNER JOIN "item_values" iv
    ON iv."field_name" = i."business_id_field_name" AND iv."item_id" = i."id"
	) as val
where items.id = val.item_id and items.business_id IS NULL;

-- name: DbListHashesPresentSimple :many
SELECT hash
FROM items
WHERE hash = ANY(@hashes::text[]);

-- name: DbListHashesPresentLatestOnly :many
SELECT c.hash FROM (
    SELECT i2.hash, row_number() over (partition by i2.business_id order by i2.created_at desc) as date_rank
    FROM items i2
    WHERE i2.business_id in (
        SELECT i1.business_id
        FROM items i1
        WHERE i1.hash = ANY (@hashes::text[]) AND i1.business_id IS NOT NULL
    )
 ) c
WHERE c.date_rank = 1 and c.hash = ANY (@hashes::text[]);

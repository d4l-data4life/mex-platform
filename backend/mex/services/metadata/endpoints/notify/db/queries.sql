-- name: DbGetItemValues :many
SELECT field_name, field_value FROM current_item_values
WHERE item_id = $1
ORDER BY field_name ASC, place ASC;

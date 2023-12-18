-- The following functions are used to derive items by their business ID which are linked to form a tree (for example, org structures).
-- In this context, we refer to items as nodes and the item's business ID is referred to as node ID.

-- This function returns all non-root nodess, that is, nodes, that have at least one outgoing edge.
create or replace function f_tree_descendants(entity_name text, link_field_name text)
returns table (node_id text, parent_node_id text)
language plpgsql as
$func$
begin
    return query
    select source_business_id as node_id, target_business_id as parent_node_id
    from items_with_links iwl
    where iwl.source_entity_name = entity_name and iwl.source_link_field_name = link_field_name;
end
$func$;

-- This function returns all root nodes, that is all nodes that have no outgoing edge.
create or replace function f_tree_roots(arg_entity_name text, link_field_name text)
returns table (node_id text)
language plpgsql as
$func$
begin
    return query
        select liwbi.business_id  as node_id
        from latest_items_with_business_id liwbi
        where liwbi.entity_name = arg_entity_name
            except
        select iwl.source_business_id as node_id
        from items_with_links iwl
        where iwl.source_entity_name = arg_entity_name and iwl.source_link_field_name = link_field_name;
end;
$func$;

-- Use recursive feature to traverse the tree structure and maintain a depth counter.
create or replace function f_tree(entity_name text, link_field_name text)
returns table (node_id text, parent_node_id text, depth int)
language plpgsql as
$func$
begin
    return query
        WITH RECURSIVE tree AS (
            SELECT 0 AS depth, tr.node_id, null as parent_node_id
            FROM f_tree_roots(entity_name, link_field_name) tr
                UNION ALL
            SELECT tree.depth + 1, td.node_id, td.parent_node_id
            FROM f_tree_descendants(entity_name, link_field_name) td
            JOIN tree
                ON td.parent_node_id = tree.node_id
        )
        SELECT tree.node_id, tree.parent_node_id, tree.depth
        FROM tree;
end;
$func$;


create or replace function f_tree_root_path(entity_name text, link_field_name text, start_node_id text)
returns table (node_id text, parent_node_id text, depth int)
language plpgsql as
$func$
begin
    return query
        with recursive root_path(node_id, parent_node_id, depth) as (
            select t.node_id, t.parent_node_id, t.depth from f_tree(entity_name, link_field_name) t where t.node_id = start_node_id
                union
            select t.node_id, t.parent_node_id, t.depth
            from  f_tree(entity_name, link_field_name) t
            inner join root_path rp
                on rp.parent_node_id = t.node_id
        )
        select * from root_path;
end;
$func$;


CREATE OR REPLACE FUNCTION next_migration_version() RETURNS integer
LANGUAGE plpgsql IMMUTABLE AS
$$
BEGIN
    return 17;
END;
$$;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "entity_types" (
    "created_at" timestamptz NOT NULL,
    "name"       text        NOT NULL,
    "config"     jsonb,

    PRIMARY KEY ("name")
);

CREATE TABLE IF NOT EXISTS "items" (
    "created_at"  timestamptz NOT NULL,
    "id"          text        NOT NULL,
    "owner"       text        NOT NULL,
    "entity_name" text        NOT NULL,
    "business_id" text,
    "business_id_field_name" text,

    PRIMARY KEY ("id"),
    CONSTRAINT "fk_entity_types_items" FOREIGN KEY ("entity_name") REFERENCES "entity_types"("name")
);

DROP INDEX IF EXISTS items_business_id_idx;
CREATE INDEX IF NOT EXISTS items_business_id_idx ON "items" USING btree ("business_id");


CREATE TABLE IF NOT EXISTS "relations" (
    "created_at"       timestamptz NOT NULL,
    "id"               text        NOT NULL,
    "deleted"          boolean     NOT NULL,
    "owner"            text        NOT NULL,
    "source_item_id"   text        NOT NULL,
    "type"             text        NOT NULL,
    "target_item_id"   text        NOT NULL,
    "info_item_id"     text,

    PRIMARY KEY ("id"),

    CONSTRAINT "fk_items_left_relations"     FOREIGN KEY ("source_item_id") REFERENCES "items"("id") ON DELETE CASCADE,
    CONSTRAINT "fk_items_right_relations"    FOREIGN KEY ("target_item_id") REFERENCES "items"("id") ON DELETE CASCADE,
    CONSTRAINT "fk_items_metadata_relations" FOREIGN KEY ("info_item_id")   REFERENCES "items"("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "field_defs" (
    "created_at"       timestamptz NOT NULL,
    "name"             text        NOT NULL,
    "kind"             text        NOT NULL,
    "display_id"       text,
    "index_def"        jsonb       NOT NULL,

    PRIMARY KEY ("name")
);


CREATE TABLE IF NOT EXISTS "item_values" (
    "created_at"      timestamptz NOT NULL,
    "id"              text        NOT NULL,
    "revision"        serial      NOT NULL,
    "deleted"         boolean     NOT NULL,
    "language"        text,
    "field_name"      text        NOT NULL,
    "field_value"     text        NOT NULL,
    "place"           integer     NOT NULL,
    "item_id"         text        NOT NULL,
    "revision_comment" text,

    PRIMARY KEY ("id","revision"),

    CONSTRAINT "fk_items_values"      FOREIGN KEY ("item_id")    REFERENCES "items"("id") ON DELETE CASCADE,
    CONSTRAINT "fk_field_defs_values" FOREIGN KEY ("field_name") REFERENCES "field_defs"("name")
);

DROP INDEX IF EXISTS item_values_field_name_idx;
CREATE INDEX IF NOT EXISTS item_values_field_name_idx ON "item_values" USING btree ("field_name");

DROP INDEX IF EXISTS item_values_item_id_idx;
CREATE INDEX IF NOT EXISTS item_values_item_id_idx ON "item_values" USING btree ("item_id");


CREATE TABLE IF NOT EXISTS "item_view_counts" (
    "created_at" timestamptz NOT NULL,
    "item_id"    text        NOT NULL,
    "user_id"    text        NOT NULL,
    "counts"     integer     NOT NULL,

    PRIMARY KEY ("item_id","user_id"),

    CONSTRAINT "fk_items_item_view_counts" FOREIGN KEY ("item_id") REFERENCES "items"("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "search_config_objects" (
    "id"    text    NOT NULL,
    "name"  text    NOT NULL,
    "obj_type"  text    NOT NULL,

    PRIMARY KEY ("id"),

    UNIQUE ("name", "obj_type")
);

CREATE TABLE IF NOT EXISTS "search_config_fields" (
    "object_id"     text    NOT NULL,
    "field_name"    text    NOT NULL,

    PRIMARY KEY ("object_id", "field_name"),

    CONSTRAINT "fk_search_config_objects_fields" FOREIGN KEY ("object_id") REFERENCES "search_config_objects"("id") ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION next_migration_version() RETURNS integer
LANGUAGE plpgsql IMMUTABLE AS
$$
BEGIN
    return 2;
END;
$$;

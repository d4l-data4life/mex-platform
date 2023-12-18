ALTER TABLE "items" DROP CONSTRAINT "fk_entity_types_items";

DROP TABLE "entity_types";

CREATE OR REPLACE FUNCTION next_migration_version() RETURNS integer
LANGUAGE plpgsql IMMUTABLE AS
$$
BEGIN
    return 19;
END;
$$;

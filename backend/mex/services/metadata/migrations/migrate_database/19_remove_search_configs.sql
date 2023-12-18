DROP VIEW "search_config";

DROP TABLE "search_config_fields";
DROP TABLE "search_config_objects";

CREATE OR REPLACE FUNCTION next_migration_version() RETURNS integer
LANGUAGE plpgsql IMMUTABLE AS
$$
BEGIN
    return 20;
END;
$$;

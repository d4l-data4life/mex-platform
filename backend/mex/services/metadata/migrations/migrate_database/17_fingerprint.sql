ALTER TABLE "items" ADD COLUMN "hash" text;

DROP INDEX IF EXISTS items_hash_idx;
CREATE INDEX IF NOT EXISTS items_hash_idx ON "items" USING btree ("hash");


CREATE OR REPLACE FUNCTION next_migration_version() RETURNS integer
LANGUAGE plpgsql IMMUTABLE AS
$$
BEGIN
    return 18;
END;
$$;

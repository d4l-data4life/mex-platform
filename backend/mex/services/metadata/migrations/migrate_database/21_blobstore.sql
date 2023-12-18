create table blob_store (blob_name text, blob_type text, blob_oid oid not null, primary key (blob_name, blob_type));

CREATE OR REPLACE FUNCTION next_migration_version() RETURNS integer
LANGUAGE plpgsql IMMUTABLE AS
$$
BEGIN
    return 22;
END;
$$;

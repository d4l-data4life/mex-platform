CREATE SCHEMA IF NOT EXISTS "mex";

SET search_path TO "mex",public;

CREATE OR REPLACE FUNCTION next_migration_version() RETURNS integer
AS $$
BEGIN
	return 1;
END;$$
LANGUAGE plpgsql IMMUTABLE;

DROP TABLE scratch_auth;
ALTER TABLE users DROP COLUMN name;
ALTER TABLE users ADD scratch_name user_name;
CREATE INDEX scratch_name_idx ON users (scratch_name);

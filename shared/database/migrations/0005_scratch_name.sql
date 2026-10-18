ALTER TABLE users
ALTER COLUMN scratch_name SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT scname_unique UNIQUE (scratch_name);

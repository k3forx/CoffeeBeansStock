-- 匿名認証のロールバック: NOT NULL 制約とインデックスを復元
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ALTER COLUMN password_hash SET NOT NULL;
ALTER TABLE users ALTER COLUMN name SET NOT NULL;

ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;

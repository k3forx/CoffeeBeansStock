ALTER TABLE usage_history DROP CONSTRAINT IF EXISTS check_usage_type;
ALTER TABLE usage_history DROP COLUMN usage_type;

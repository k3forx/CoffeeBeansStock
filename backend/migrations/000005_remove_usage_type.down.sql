ALTER TABLE usage_history ADD COLUMN usage_type VARCHAR(50) NOT NULL DEFAULT 'manual';
ALTER TABLE usage_history ADD CONSTRAINT check_usage_type CHECK (usage_type IN ('manual', 'quick_button'));
ALTER TABLE usage_history ALTER COLUMN usage_type DROP DEFAULT;

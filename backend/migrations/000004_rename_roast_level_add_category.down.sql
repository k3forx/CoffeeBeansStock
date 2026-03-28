-- roast_detail が NULL の行を roast_level から補完
UPDATE coffee_beans SET roast_detail = CASE roast_level
    WHEN 'shallow' THEN 'light'
    WHEN 'medium' THEN 'medium'
    WHEN 'medium_deep' THEN 'city'
    WHEN 'deep' THEN 'french'
END WHERE roast_detail IS NULL;

ALTER TABLE coffee_beans DROP CONSTRAINT IF EXISTS check_roast_level;
ALTER TABLE coffee_beans DROP CONSTRAINT IF EXISTS check_roast_detail;
ALTER TABLE coffee_beans DROP CONSTRAINT IF EXISTS check_roast_level_detail_consistency;
ALTER TABLE coffee_beans DROP COLUMN roast_level;
ALTER TABLE coffee_beans ALTER COLUMN roast_detail SET NOT NULL;
ALTER TABLE coffee_beans RENAME COLUMN roast_detail TO roast_level;
ALTER TABLE coffee_beans ADD CONSTRAINT check_roast_level
    CHECK (roast_level IN ('light', 'cinnamon', 'medium', 'high', 'city', 'full_city', 'french', 'italian'));

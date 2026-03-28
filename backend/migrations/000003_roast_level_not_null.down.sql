ALTER TABLE coffee_beans DROP CONSTRAINT IF EXISTS check_roast_level;
ALTER TABLE coffee_beans ALTER COLUMN roast_level DROP NOT NULL;

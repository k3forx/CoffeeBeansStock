-- Set default roast level for existing NULL values
UPDATE coffee_beans SET roast_level = 'medium' WHERE roast_level IS NULL;

-- Make roast_level NOT NULL
ALTER TABLE coffee_beans ALTER COLUMN roast_level SET NOT NULL;

-- Add CHECK constraint for valid roast levels
ALTER TABLE coffee_beans ADD CONSTRAINT check_roast_level
    CHECK (roast_level IN ('light', 'cinnamon', 'medium', 'high', 'city', 'full_city', 'french', 'italian'));

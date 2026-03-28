-- 1. 旧CHECK制約を削除
ALTER TABLE coffee_beans DROP CONSTRAINT IF EXISTS check_roast_level;

-- 2. roast_level → roast_detail にリネーム
ALTER TABLE coffee_beans RENAME COLUMN roast_level TO roast_detail;

-- 3. roast_detail を NULL 許可に変更
ALTER TABLE coffee_beans ALTER COLUMN roast_detail DROP NOT NULL;

-- 4. roast_level カラムを新規追加（大分類）
ALTER TABLE coffee_beans ADD COLUMN roast_level VARCHAR(50);

-- 5. 既存データをバックフィル（roast_detail から roast_level を導出）
UPDATE coffee_beans SET roast_level = CASE
    WHEN roast_detail IN ('light', 'cinnamon') THEN 'shallow'
    WHEN roast_detail IN ('medium', 'high') THEN 'medium'
    WHEN roast_detail IN ('city', 'full_city') THEN 'medium_deep'
    WHEN roast_detail IN ('french', 'italian') THEN 'deep'
END;

-- 6. roast_level に NOT NULL 制約
ALTER TABLE coffee_beans ALTER COLUMN roast_level SET NOT NULL;

-- 7. 新CHECK制約を追加
ALTER TABLE coffee_beans ADD CONSTRAINT check_roast_level
    CHECK (roast_level IN ('shallow', 'medium', 'medium_deep', 'deep'));

ALTER TABLE coffee_beans ADD CONSTRAINT check_roast_detail
    CHECK (roast_detail IS NULL OR roast_detail IN (
        'light', 'cinnamon', 'medium', 'high', 'city', 'full_city', 'french', 'italian'));

ALTER TABLE coffee_beans ADD CONSTRAINT check_roast_level_detail_consistency
    CHECK (
        roast_detail IS NULL
        OR (roast_level = 'shallow' AND roast_detail IN ('light', 'cinnamon'))
        OR (roast_level = 'medium' AND roast_detail IN ('medium', 'high'))
        OR (roast_level = 'medium_deep' AND roast_detail IN ('city', 'full_city'))
        OR (roast_level = 'deep' AND roast_detail IN ('french', 'italian'))
    );

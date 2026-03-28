-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255),
    password_hash VARCHAR(255),
    name VARCHAR(100),
    low_stock_threshold INTEGER NOT NULL DEFAULT 100,
    notification_enabled BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Coffee beans table (master data without purchase info)
CREATE TABLE coffee_beans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    origin VARCHAR(100),
    roast_level VARCHAR(50) NOT NULL,
    current_stock INTEGER NOT NULL CHECK (current_stock >= 0),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_coffee_beans_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT check_roast_level
        CHECK (roast_level IN ('light', 'cinnamon', 'medium', 'high', 'city', 'full_city', 'french', 'italian'))
);

-- Purchase history table (separated from coffee_beans)
CREATE TABLE purchase_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    coffee_bean_id UUID NOT NULL,
    user_id UUID NOT NULL,
    purchase_date DATE NOT NULL,
    purchase_price DECIMAL(10, 2),
    purchase_store VARCHAR(200),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_purchase_history_coffee_bean
        FOREIGN KEY (coffee_bean_id)
        REFERENCES coffee_beans(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_purchase_history_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Usage history table
CREATE TABLE usage_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    coffee_bean_id UUID NOT NULL,
    user_id UUID NOT NULL,
    usage_date DATE NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    usage_type VARCHAR(50) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_usage_history_coffee_bean
        FOREIGN KEY (coffee_bean_id)
        REFERENCES coffee_beans(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_usage_history_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT check_usage_type
        CHECK (usage_type IN ('manual', 'quick_button'))
);

-- Indexes for users (email index removed for anonymous auth)

-- Indexes for coffee_beans
CREATE INDEX idx_coffee_beans_user_id ON coffee_beans(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_coffee_beans_updated_at ON coffee_beans(updated_at DESC) WHERE deleted_at IS NULL;

-- Indexes for purchase_history
CREATE INDEX idx_purchase_history_coffee_bean_id ON purchase_history(coffee_bean_id);
CREATE INDEX idx_purchase_history_user_id ON purchase_history(user_id);
CREATE INDEX idx_purchase_history_purchase_date ON purchase_history(purchase_date DESC);
CREATE INDEX idx_purchase_history_coffee_bean_date ON purchase_history(coffee_bean_id, purchase_date DESC);

-- Indexes for usage_history
CREATE INDEX idx_usage_history_coffee_bean_id ON usage_history(coffee_bean_id);
CREATE INDEX idx_usage_history_user_id ON usage_history(user_id);
CREATE INDEX idx_usage_history_usage_date ON usage_history(usage_date DESC);
CREATE INDEX idx_usage_history_coffee_bean_date ON usage_history(coffee_bean_id, usage_date DESC);

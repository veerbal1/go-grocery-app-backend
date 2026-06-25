-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS list (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS list_item (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    unit VARCHAR(10) NOT NULL,
    price INT NOT NULL DEFAULT 0,
    list_id INT REFERENCES list(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_list_user_id ON list(user_id);
CREATE INDEX IF NOT EXISTS idx_list_item_list_id ON list_item(list_id);
CREATE INDEX IF NOT EXISTS idx_users_email on users(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_list_item_list_id;
DROP INDEX IF EXISTS idx_list_user_id;

DROP TABLE IF EXISTS list_item;
DROP TABLE IF EXISTS list;
-- +goose StatementEnd
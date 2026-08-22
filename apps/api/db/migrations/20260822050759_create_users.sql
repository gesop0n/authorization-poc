-- +goose Up

-- NOTE: https://github.com/pressly/goose?utm_source=chatgpt.com#:~:text=More%20complex%20statements%20(PL/pgSQL)%20that%20have%20semicolons%20within%20them%20must%20be%20annotated%20with%20%2D%2D%20%2Bgoose%20StatementBegin%20and%20%2D%2D%20%2Bgoose%20StatementEnd%20to%20be%20properly%20recognized.%20For%20example%3A
-- +goose StatementBegin
-- updated_at を自動で更新する
CREATE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TABLE users (
    id UUID PRIMARY KEY,
    display_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,


    CONSTRAINT users_status_check
        CHECK (status IN ('active', 'suspended')),

    -- DB側でも空白の名前を拒否する
    CONSTRAINT users_display_name_not_blank
        CHECK (btrim(display_name) <> '')
);

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE users;
DROP FUNCTION set_updated_at();

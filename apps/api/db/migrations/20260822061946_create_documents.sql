-- +goose Up
CREATE TABLE documents (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL,
    owner_user_id UUID NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Project削除時にDocumentも削除
    CONSTRAINT documents_project_id_fkey
        FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,

    -- Userが参照されている場合は削除を拒否
    CONSTRAINT documents_owner_user_id_fkey
        FOREIGN KEY (owner_user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    -- DB側でも空白の名前を拒否する
    CONSTRAINT documents_title_not_blank
        CHECK (btrim(title) <> '')
);

CREATE INDEX documents_project_id_idx
    ON documents (project_id);

CREATE INDEX documents_owner_user_id_idx
    ON documents (owner_user_id);

CREATE TRIGGER documents_set_updated_at
BEFORE UPDATE ON documents
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE documents;

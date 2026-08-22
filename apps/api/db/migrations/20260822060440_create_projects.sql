-- +goose Up
CREATE TABLE projects (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    workspace_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT projects_workspace_id_fkey
        FOREIGN KEY (workspace_id)
        REFERENCES workspaces(id)
        ON DELETE CASCADE,

    -- DB側でも空白の名前を拒否する
    CONSTRAINT projects_name_not_blank
        CHECK (btrim(name) <> ''),

    CONSTRAINT projects_status_check
        CHECK (status IN ('active', 'archived'))
);

CREATE INDEX projects_workspace_id_idx
    ON projects (workspace_id);

CREATE TRIGGER projects_set_updated_at
BEFORE UPDATE ON projects
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE projects;

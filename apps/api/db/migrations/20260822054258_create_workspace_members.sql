-- +goose Up
CREATE TABLE workspace_members (
    workspace_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (workspace_id, user_id),

    -- Workspace削除時にメンバーも削除
    CONSTRAINT workspace_members_workspace_id_fkey
        FOREIGN KEY (workspace_id)
        REFERENCES workspaces(id)
        ON DELETE CASCADE,

    -- Userが参照されている場合は削除を拒否
    CONSTRAINT workspace_members_user_id_fkey
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT workspace_members_role_check
        CHECK (role IN ('owner', 'admin', 'member'))
);

CREATE TRIGGER workspace_members_set_updated_at
BEFORE UPDATE ON workspace_members
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE INDEX workspace_members_user_id_idx
    ON workspace_members (user_id);


-- +goose Down
DROP TABLE workspace_members;

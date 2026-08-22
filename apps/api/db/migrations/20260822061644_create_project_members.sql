-- +goose Up
CREATE TABLE project_members (
    project_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (project_id, user_id),

    -- Project削除時にメンバーも削除
    CONSTRAINT project_members_project_id_fkey
        FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE,

    -- Userが参照されている場合は削除を拒否
    CONSTRAINT project_members_user_id_fkey
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT project_members_role_check
        CHECK (role IN ('admin', 'editor', 'viewer'))
);

CREATE INDEX project_members_user_id_idx
    ON project_members (user_id);

CREATE TRIGGER project_members_set_updated_at
BEFORE UPDATE ON project_members
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE project_members;

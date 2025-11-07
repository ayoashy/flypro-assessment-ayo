-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS expense_reports (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    title VARCHAR(200) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    total DECIMAL(10, 2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_reports_user_id ON expense_reports(user_id);
CREATE INDEX IF NOT EXISTS idx_reports_status ON expense_reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_deleted_at ON expense_reports(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_reports_deleted_at;
DROP INDEX IF EXISTS idx_reports_status;
DROP INDEX IF EXISTS idx_reports_user_id;
DROP TABLE IF EXISTS expense_reports;
-- +goose StatementEnd




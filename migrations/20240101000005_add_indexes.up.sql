-- +goose Up
-- +goose StatementBegin
-- Additional composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_expenses_user_status ON expenses(user_id, status);
CREATE INDEX IF NOT EXISTS idx_expenses_user_category ON expenses(user_id, category);
CREATE INDEX IF NOT EXISTS idx_reports_user_status ON expense_reports(user_id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_reports_user_status;
DROP INDEX IF EXISTS idx_expenses_user_category;
DROP INDEX IF EXISTS idx_expenses_user_status;
-- +goose StatementEnd




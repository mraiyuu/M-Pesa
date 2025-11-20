-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS access_tokens(
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(20),
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW (),
    updated_at TIMESTAMPTZ DEFAULT NOW
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS access_tokens;
-- +goose StatementEnd

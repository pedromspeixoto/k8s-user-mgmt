-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_files (
    id           INT NOT NULL AUTO_INCREMENT,
    user_id      VARCHAR(255) NOT NULL,
    file_id      VARCHAR(255) NOT NULL,
    file_type    VARCHAR(255) NOT NULL,
    file_content LONGBLOB NOT NULL,
    created_at   DATETIME(3) NULL,
    updated_at   DATETIME(3) NULL,
    deleted_at   DATETIME(3) NULL,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_files;
-- +goose StatementEnd
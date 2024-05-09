-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
    ADD COLUMN isDraft bool not null default false;
ALTER TABLE events
    ADD COLUMN isMachineGenerated bool not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
    DROP COLUMN isDraft;
ALTER TABLE events
    DROP COLUMN isMachineGenerated;
-- +goose StatementEnd

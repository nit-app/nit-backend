-- +goose Up
-- +goose StatementBegin
ALTER TABLE events
ADD COLUMN plainDescription VARCHAR(2048);

UPDATE events SET plainDescription = 'an absolutely free event for absolutely free people, come and have fun'
WHERE uuid = 'f70dd14f-8e24-11ee-8542-fa163e445fa2';

UPDATE events SET plainDescription = 'a paid event, come and have fun if you can afford it'
WHERE uuid = 'b141fa84-8e25-11ee-8542-fa163e445fa2';

ALTER TABLE events
ALTER COLUMN plainDescription SET NOT NULL,
ADD CHECK ( LENGTH(events.plainDescription) > 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events
DROP COLUMN plainDescription;
-- +goose StatementEnd

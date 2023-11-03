-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    "uuid"       uuid PRIMARY KEY,
    title        VARCHAR(512)  NOT NULL CHECK ( LENGTH(title) > 0 ),
    description  text          NOT NULL,
    priceLow     INTEGER       NOT NULL CHECK ( priceLow >= 0 AND priceLow <= priceHigh ),
    priceHigh    INTEGER       NOT NULL CHECK ( priceHigh >= 0 AND priceHigh >= priceLow ),
    ageLimitLow  SMALLINT      NOT NULL CHECK ( ageLimitLow >= 0 ),
    ageLimitHigh SMALLINT CHECK ( ageLimitHigh IS NULL OR ageLimitHigh >= ageLimitLow ),
    location     VARCHAR(2048) NOT NULL,
    ownerInfo    VARCHAR(2048) NOT NULL,
    favCount     INTEGER       NOT NULL CHECK (favCount >= 0) DEFAULT 0,
    createdAt    TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modifiedAt   TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deletedAt    TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE event_external_links
(
    "uuid"  uuid PRIMARY KEY,
    url     varchar(1024) NOT NULL CHECK ( LENGTH(url) > 0 ),
    addedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_event_uuid FOREIGN KEY ("uuid") REFERENCES events ("uuid")
);

CREATE TABLE event_schedule
(
    "uuid"   uuid PRIMARY KEY,
    beginsAt TIMESTAMP WITH TIME ZONE NOT NULL,
    endsAt   TIMESTAMP WITH TIME ZONE NOT NULL,
    addedAt  TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_event_uuid FOREIGN KEY ("uuid") REFERENCES events ("uuid")
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event_external_links;
DROP TABLE event_schedule;
DROP TABLE events CASCADE;
-- +goose StatementEnd

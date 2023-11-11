-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    "uuid"         UUID PRIMARY KEY,
    title          VARCHAR(512)  NOT NULL CHECK ( LENGTH(title) > 0 ),
    description    TEXT          NOT NULL,
    priceLow       INTEGER       NOT NULL CHECK ( priceLow >= 0 AND priceLow <= priceHigh ),
    priceHigh      INTEGER       NOT NULL CHECK ( priceHigh >= 0 AND priceHigh >= priceLow ),
    ageLimitLow    SMALLINT      NOT NULL CHECK ( ageLimitLow >= 0 ),
    ageLimitHigh   SMALLINT CHECK ( ageLimitHigh IS NULL OR ageLimitHigh >= ageLimitLow ),
    location       VARCHAR(2048) NOT NULL,
    ownerInfo      VARCHAR(2048) NOT NULL,
    hasCertificate BOOLEAN       NOT NULL                       DEFAULT FALSE,
    favCount       INTEGER       NOT NULL CHECK (favCount >= 0) DEFAULT 0,
    createdAt      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modifiedAt     TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deletedAt      TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE event_external_links
(
    "linkUuid"  uuid PRIMARY KEY,
    "eventUuid" uuid          NOT NULL,
    title       varchar(512)  NOT NULL CHECK ( LENGTH(title) > 0 ),
    url         varchar(1024) NOT NULL CHECK ( LENGTH(url) > 0 ),
    addedAt     TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_event_uuid FOREIGN KEY ("eventUuid") REFERENCES events ("uuid")
);

CREATE TABLE event_schedule
(
    scheduleUuid uuid PRIMARY KEY,
    "eventUuid"  uuid                     NOT NULL,
    beginsAt     TIMESTAMP WITH TIME ZONE NOT NULL,
    endsAt       TIMESTAMP WITH TIME ZONE NOT NULL,
    addedAt      TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_event_uuid FOREIGN KEY ("eventUuid") REFERENCES events ("uuid")
);

CREATE TABLE event_tags
(
    "uuid" uuid,
    tag    varchar(512),

    CONSTRAINT fk_event_uuid FOREIGN KEY ("uuid") REFERENCES events ("uuid")
)

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event_external_links;
DROP TABLE event_schedule;
DROP TABLE event_tags;
DROP TABLE events;
-- +goose StatementEnd

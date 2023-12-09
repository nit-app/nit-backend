-- +goose Up
-- +goose StatementBegin
INSERT INTO events (uuid, title, description, pricelow, pricehigh, agelimitlow, agelimithigh, location, ownerinfo)
VALUES ('f70dd14f-8e24-11ee-8542-fa163e445fa2',
        'Test free event',
        'Description of a test event',
        0,
        0,
        0,
        0,
        '56.822523, 60.605641 улица Декабристов, 77Б, Екатеринбург, Свердловская область, 620063',
        'ООО Тестовые Решения');

INSERT INTO event_external_links ("linkUuid", "eventUuid", title, url, addedat)
VALUES (gen_random_uuid(), 'f70dd14f-8e24-11ee-8542-fa163e445fa2',
        'Telegram',
        'https://t.me/.___adfasdfasdfass',
        CURRENT_TIMESTAMP);

INSERT INTO event_schedule (scheduleuuid, "eventUuid", beginsat, endsat)
VALUES (gen_random_uuid(), 'f70dd14f-8e24-11ee-8542-fa163e445fa2', '2024-06-06T12:00:00+05:00',
        '2024-06-06T18:00:00+05:00');

INSERT INTO event_tags (uuid, tag)
VALUES ('f70dd14f-8e24-11ee-8542-fa163e445fa2', 'test'),
       ('f70dd14f-8e24-11ee-8542-fa163e445fa2', 'лекция'),
       ('f70dd14f-8e24-11ee-8542-fa163e445fa2', 'митап'),
       ('f70dd14f-8e24-11ee-8542-fa163e445fa2', 'квиз');

INSERT INTO events (uuid, title, description, pricelow, pricehigh, agelimitlow, agelimithigh, location, ownerinfo)
VALUES ('b141fa84-8e25-11ee-8542-fa163e445fa2',
        'Test paid event',
        'Description of a test event11',
        100,
        500,
        14,
        100,
        '56.822523, 60.605641 улица Декабристов, 77Б, Екатеринбург, Свердловская область, 620063',
        'ООО Тестовые Решения');

INSERT INTO event_external_links ("linkUuid", "eventUuid", title, url, addedat)
VALUES (gen_random_uuid(), 'b141fa84-8e25-11ee-8542-fa163e445fa2',
        'VK',
        'https://vk.com/durov',
        CURRENT_TIMESTAMP);

INSERT INTO event_schedule (scheduleuuid, "eventUuid", beginsat, endsat)
VALUES (gen_random_uuid(), 'b141fa84-8e25-11ee-8542-fa163e445fa2', '2024-06-06T12:00:00+05:00',
        '2024-06-06T18:00:00+05:00');

INSERT INTO event_tags (uuid, tag)
VALUES ('b141fa84-8e25-11ee-8542-fa163e445fa2', 'test'),
       ('b141fa84-8e25-11ee-8542-fa163e445fa2', 'хакатон'),
       ('b141fa84-8e25-11ee-8542-fa163e445fa2', 'командное'),
       ('b141fa84-8e25-11ee-8542-fa163e445fa2', 'соревнование');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM event_tags
WHERE uuid IN ('b141fa84-8e25-11ee-8542-fa163e445fa2', 'f70dd14f-8e24-11ee-8542-fa163e445fa2');
DELETE
FROM event_schedule
WHERE "eventUuid" IN ('b141fa84-8e25-11ee-8542-fa163e445fa2', 'f70dd14f-8e24-11ee-8542-fa163e445fa2');
DELETE
FROM event_external_links
WHERE "eventUuid" IN ('b141fa84-8e25-11ee-8542-fa163e445fa2', 'f70dd14f-8e24-11ee-8542-fa163e445fa2');
DELETE
FROM events
WHERE uuid IN ('b141fa84-8e25-11ee-8542-fa163e445fa2', 'f70dd14f-8e24-11ee-8542-fa163e445fa2');
-- +goose StatementEnd

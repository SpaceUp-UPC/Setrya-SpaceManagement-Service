CREATE TABLE IF NOT EXISTS spaces (
    id              VARCHAR(64) PRIMARY KEY,
    owner_id        VARCHAR(64) NOT NULL,
    name            VARCHAR(160) NOT NULL,
    description     TEXT NOT NULL,
    location_json   JSON NOT NULL,
    dimensions_json JSON NOT NULL,
    status          VARCHAR(24) NOT NULL,
    created_at      DATETIME(6) NOT NULL,
    updated_at      DATETIME(6) NOT NULL,
    INDEX idx_spaces_owner_created (owner_id, created_at)
);

CREATE TABLE IF NOT EXISTS linked_iot_devices (
    id          VARCHAR(64) PRIMARY KEY,
    space_id    VARCHAR(64) NOT NULL,
    device_id   VARCHAR(64) NOT NULL UNIQUE,
    device_type VARCHAR(80) NOT NULL,
    status      VARCHAR(24) NOT NULL,
    linked_at   DATETIME(6) NOT NULL,
    INDEX idx_linked_devices_space (space_id),
    CONSTRAINT fk_linked_devices_space FOREIGN KEY (space_id) REFERENCES spaces(id)
);

CREATE TABLE IF NOT EXISTS reviews (
    id         VARCHAR(64) PRIMARY KEY,
    space_id   VARCHAR(64) NOT NULL,
    author_id  VARCHAR(64) NOT NULL,
    rating     TINYINT UNSIGNED NOT NULL,
    comment    TEXT NOT NULL,
    created_at DATETIME(6) NOT NULL,
    INDEX idx_reviews_space_created (space_id, created_at),
    CONSTRAINT fk_reviews_space FOREIGN KEY (space_id) REFERENCES spaces(id),
    CONSTRAINT chk_reviews_rating CHECK (rating BETWEEN 1 AND 5)
);

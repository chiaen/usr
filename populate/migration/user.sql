
CREATE TABLE `user_records` (
  `uid`                    CHAR(36)     NOT NULL,
  `email`                  VARCHAR(255) NOT NULL,
  `password`               CHAR(128) NOT NULL,
  PRIMARY KEY (`uid`)
);

CREATE TABLE `refresh_tokens` (
  `uid`         CHAR(36)     NOT NULL,
  `token`       CHAR(64)     NOT NULL,
  `issued_at`   BIGINT       NOT NULL,
  `valid_until` BIGINT       NOT NULL,
  PRIMARY KEY (`token`)
)

CREATE TABLE `user_has_interests` (
  `user_id`     CHAR(36)  NOT NULL,
  `interest_id`  CHAR(36) NOT NULL,
  PRIMARY KEY (`user_id`, `interest_id`)
)

CREATE TABLE `interests` (
  `id`      VARCHAR(36)  NOT NULL,
  `name`    TINYTEXT,
  PRIMARY KEY (`id`)
)
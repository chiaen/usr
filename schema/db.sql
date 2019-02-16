
CREATE TABLE `user_records` (
  `uid`                    CHAR(36)     NOT NULL,
  `email`                  VARCHAR(255) NOT NULL,
  `password`               CHAR(128) NOT NULL,
  PRIMARY KEY (`uid`),
  INDEX `idx_user_records_email` (`email` ASC)
);

CREATE TABLE `user_has_interests` (
  `user_id`     CHAR(36)  NOT NULL,
  `interest_id`  INT NOT NULL,
  `published` TINYINT NOT NULL,
  PRIMARY KEY (`user_id`, `interest_id`),
  INDEX `idx_user_has_interests_uid_interest` (`user_id` ASC, `interest_id` ASC)
);

CREATE TABLE `interests` (
  `id`      INT  NOT NULL AUTO_INCREMENT,
  `name`    TINYTEXT,
  PRIMARY KEY (`id`)
);
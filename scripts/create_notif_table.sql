CREATE TABLE if not exists `notifications`
(
    `id`         int  NOT NULL AUTO_INCREMENT,
    `num`        long  NOT NULL,
    `user_login` long NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;


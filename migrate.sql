-- user
CREATE TABLE `user` (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(40) UNIQUE NOT NULL,
    `email` VARCHAR(40) UNIQUE NOT NULL,
    `enable` boolean DEFAULT true,
    `hashed_password` VARCHAR(512),
    `team_id` INT(60) NOT NULL,
    `forget_hash` VARCHAR(512) DEFAULT NULL,
    `register_hash` VARCHAR(512) DEFAULT NULL
)DEFAULT CHARSET=utf8;

-- team
CREATE TABLE `team` (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(40) UNIQUE NOT NULL,
    `enable` boolean DEFAULT true,
    `is_admin` boolean DEFAULT false,
    `hashed_password` VARCHAR(256)
)DEFAULT CHARSET=utf8;

-- question
CREATE TABLE `question` (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(40) UNIQUE NOT NULL,
    `flag` VARCHAR(128),
    `genre` VARCHAR(40),
    `score` INT(60),
    `author_id` INT(60),
    `publish_start_time` TIMESTAMP NOT NULL DEFAULT NOW(),
    `sentence` TEXT
)DEFAULT CHARSET=utf8;

CREATE TABLE `answer` (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT(60),
    `question_id` INT(60),
    `flag` VARCHAR(128),
    `create_time` TIMESTAMP(4) NOT NULL DEFAULT CURRENT_TIMESTAMP(4)
)DEFAULT CHARSET=utf8;

CREATE TABLE `question_file` (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `question_id` INT(60),
    `name` VARCHAR(40),
    `file_hash` VARCHAR(256)
)DEFAULT CHARSET=utf8;

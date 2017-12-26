-- user
CREATE TABLE user (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(40) UNIQUE NOT NULL,
    `email` VARCHAR(40) UNIQUE NOT NULL,
    `enable` INT(1) DEFAULT '0',
    `hashed_password` VARCHAR(512)
)DEFAULT CHARSET=utf8;

CREATE TABLE user_forget (
    `id` INT(60)PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT(60),
    `url_hash` VARCHAR(256),
    `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)DEFAULT CHARSET=utf8;

CREATE TABLE user_register (
    `id` INT(60)PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT(60),
    `url_hash` VARCHAR(256),
    `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)DEFAULT CHARSET=utf8;

-- team
CREATE TABLE team (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(40) UNIQUE NOT NULL,
    `enable` INT(1) DEFAULT '1',
    `hashed_password` VARCHAR(256)
)DEFAULT CHARSET=utf8;

CREATE TABLE team_member (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `team_id` INT(60),
    `user_id` INT(60)
)DEFAULT CHARSET=utf8;

-- question
CREATE TABLE question (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(40) UNIQUE NOT NULL,
    `flag` VARCHAR(40),
    `score` INT(60),
    `sentence` TEXT
)DEFAULT CHARSET=utf8;

CREATE TABLE answer (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT(60),
    `question_id` INT(60),
    `flag` VARCHAR(40),
    `create_time` DATE NOT NULL
)DEFAULT CHARSET=utf8;

CREATE TABLE genre_type (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `name` VARCHAR(40) UNIQUE NOT NULL
)DEFAULT CHARSET=utf8;

CREATE TABLE question_genre (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `question_id` INT(60),
    `genre_type_id` INT(60)
)DEFAULT CHARSET=utf8;

CREATE TABLE question_file (
    `id` INT(60) PRIMARY KEY AUTO_INCREMENT,
    `question_id` INT(60),
    `name` VARCHAR(40),
    `file_hash` VARCHAR(256)
)DEFAULT CHARSET=utf8;

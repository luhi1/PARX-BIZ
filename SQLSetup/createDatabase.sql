-- Active: 1712548880115@@127.0.0.1@3306@fblacnp
DROP SCHEMA `fblacnp`;
CREATE SCHEMA `fblacnp`;
use fblacnp;

-- User Tables
-- PA == Program Area;
CREATE TABLE `Program_Areas`
(
    `id`           BIGINT AUTO_INCREMENT,
    `name`   VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
);
insert into Program_Areas(name)
values
    ('Architectural Design'),
    ('Biomedical'),
    ('Business'),
    ('Computer Science'),
    ('Cybersecurity'),
    ('Engineering'),
    ('Graphic Design');
CREATE TABLE `Users`
(
    `id`      BIGINT UNSIGNED NOT NULL,
    `username` VARCHAR(255)       NOT NULL,
    `password`    VARCHAR(255)       NOT NULL,
    `real_name` VARCHAR(255)      NOT NULL,
    `program_area`     BIGINT not null,
    PRIMARY KEY (`id`),
    FOREIGN KEY (program_area) references Program_Areas(id)
);
insert into Users(id, username, `password`, real_name, program_area)
values
    (1, 'test', 'n4bQgYhMfWWaL-qgxVrQFaO_TxsrC4Is0V1sFbDwCgg=', 'test', 1);



-- Partners Tables
CREATE TABLE `Representatives`
(
    `id`           BIGINT AUTO_INCREMENT,
    `email` VARCHAR(255) NOT NULL,
    `phone` BIGINT NOT NULL,
    PRIMARY KEY (`id`)
);
insert into Representatives( email, phone)
values
    ('cathy@gmail.com', 2245551234);

CREATE TABLE `Partner_Types`
(
    `id`           BIGINT AUTO_INCREMENT,
    `name`   VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
);
insert into Partner_Types(`name`)
values
    ('Non-Profit'),
    ('Charity'),
    ('Student-Run Organization'),
    ('Cooperative'),
    ('Event Coordinator'),
    ('Business'),
    ('Local Business');

CREATE TABLE `Partners`
(
    `id`        BIGINT UNSIGNED   NOT NULL AUTO_INCREMENT,
    `name`      VARCHAR(255)      NOT NULL,
    `representative` BIGINT NOT NULL,
    `type` BIGINT NOT NULL,
    `active` BOOLEAN not null default 0,
    PRIMARY KEY (`id`),
    FOREIGN KEY (representative) references Representatives(id),
    FOREIGN KEY (`type`) references Partner_Types(id)
);
insert into Partners(`name`, representative, `type`, `active`)
values
    ('TestOrg', 1, 1, 1),
    ('TestOrg2', 1, 1, 1),
    ('TestOrg3', 1, 1, 1),
    ('TestOrg4', 1, 1, 1),
    ('TestOrg5', 1, 1, 1),
    ('TestOrg6', 1, 1, 1),
    ('TestOrg7', 1, 1, 1),
    ('TestOrg8', 1, 1, 1),
    ('TestOrg9', 1, 1, 1),
    ('TestOrg10', 1, 1, 1),
    ('TestOrg11', 1, 1, 1),
    ('TestOrg12', 1, 1, 1),
    ('TestOrg13', 1, 1, 1),
    ('TestOrg14', 1, 1, 1),
    ('TestOrg15', 1, 1, 1),
    ('TestOrg16', 1, 1, 1),
    ('TestOrg17', 1, 1, 1),
    ('TestOrg18', 1, 1, 1),
    ('TestOrg19', 1, 1, 1),
    ('TestOrg20', 1, 1, 1),
    ('TestOrg21', 1, 1, 1),
    ('TestOrg22', 1, 1, 1),
    ('TestOrg23', 1, 1, 1),
    ('TestOrg24', 1, 1, 1);
CREATE TABLE `Resources`
(
    `id`        BIGINT UNSIGNED   NOT NULL AUTO_INCREMENT,
    `partner`      BIGINT UNSIGNED   NOT NULL,
    `info` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (partner) references Partners(id)
);
insert into Resources(partner, info)
values
    (1, 'We help test your code'),
    (1, 'We help 2 test your code'),
    (3, 'We help test your code off the jawn'),
    (4, 'We help test your code'),
    (5, 'We help 2 test your code'),
    (6, 'We help test your code off the jawn'),
    (7, 'We help test your code'),
    (8, 'We help 2 test your code'),
    (9, 'We help test your code off the jawn'),
    (10, 'We help test your code'),
    (11, 'We help 2 test your code'),
    (12, 'We help test your code off the jawn'),
    (13, 'We help test your code'),
    (14, 'We help 2 test your code'),
    (15, 'We help test your code off the jawn'),
    (16, 'We help test your code'),
    (17, 'We help 2 test your code'),
    (18, 'We help test your code off the jawn'),
    (19, 'We help test your code'),
    (20, 'We help 2 test your code'),
    (21, 'We help test your code off the jawn'),
    (22, 'We help test your code'),
    (23, 'We help 2 test your code'),
    (24, 'We help test your code off the jawn'),
    (24, 'We help test your code off the jawn');


select Users.username, Users.password, Users.real_name,
       Program_Areas.name
       from Users join Program_Areas on Users.program_area = Program_Areas.id;
select Resources.info, Representatives.email,
       Representatives.phone, Partner_Types.name, Partners.active
       from Resources join Partners on Resources.partner = Partners.id
                        join Representatives on Partners.representative = Representatives.id
                        join Partner_Types on Partners.type = Partner_Types.id;
select Users.id, Users.real_name, Program_Areas.name from Users join Program_Areas on Users.program_area = Program_Areas.id where Users.username = 'test' && Users.password = 'n4bQgYhMfWWaL-qgxVrQFaO_TxsrC4Is0V1sFbDwCgg='
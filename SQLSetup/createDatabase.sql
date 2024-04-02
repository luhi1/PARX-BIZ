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
insert into users(id, username, `password`, real_name, program_area)
values
    (1, 'test', 'n4bQgYhMfWWaL-qgxVrQFaO_TxsrC4Is0V1sFbDwCgg=', 'test', 1);



-- Partners Tables
CREATE TABLE `Representatives`
(
    `id`           BIGINT AUTO_INCREMENT,
    `name`   VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(255) NOT NULL,
    PRIMARY KEY (`id`)
);
insert into Representatives(`name`, email, phone)
values
    ('Cathy', 'cathy@gmail.com', '224-555-1234');

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
    PRIMARY KEY (`id`),
    FOREIGN KEY (representative) references Representatives(id),
    FOREIGN KEY (`type`) references Partner_Types(id)
);
insert into Partners(`name`, representative, `type`)
values
    ('TestOrg', 1, 1);
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
    (1, 'We help test your code');
select Users.username, Users.password, Users.real_name,
       Program_Areas.name
       from Users join Program_Areas on Users.program_area = Program_Areas.id;
select Resources.info, Representatives.name, Representatives.email,
       Representatives.phone, Partner_Types.name
       from Resources join Partners on Resources.partner = Partners.id
                        join Representatives on Partners.representative = Representatives.id
                        join Partner_Types on Partners.type = Partner_Types.id;
select users.id, users.real_name, Program_Areas.name from users join program_areas on users.program_area = program_areas.id where users.username = 'test' && users.password = 'n4bQgYhMfWWaL-qgxVrQFaO_TxsrC4Is0V1sFbDwCgg='
/*
Connect to your mysql DB using a terminal client:
$ mysql -h 127.0.0.1 -P <port> -u root -p

Then select the database first:
$ USE <database name>;

Finally, execute this file there:
$ source <path/to/create-tables.sql>
*/

DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
    `id` INT AUTO_INCREMENT NOT NULL,
    `name` VARCHAR(128) NOT NULL,
    `completed` BOOLEAN,
    PRIMARY KEY (`id`)
);

INSERT INTO `task` (`name`, `completed`) VALUES
('Complete Go Challenges', false),
('Clean the dishes', false),
('Learn GO', true);

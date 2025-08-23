CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE ON snippetboxdb.* TO 'web'@'localhost';
FLUSH PRIVILEGES;

ALTER USER 'web'@'localhost' IDENTIFIED BY '1234';
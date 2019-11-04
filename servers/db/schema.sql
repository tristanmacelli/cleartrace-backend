drop table if exists users;

create table users (
    ID integer not null auto_increment primary key, -- add auto increment
    email varchar(320),
    passHash varbinary(80),
    username varchar(255),
    firstname varchar(200),
    lastname varchar(200),
    photoURL varchar(300)
)

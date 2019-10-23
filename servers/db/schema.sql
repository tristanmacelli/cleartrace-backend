drop table if exists users;

create table users (
    ID integer,
    email varchar(320),
    passHash varchar(100),
    username varchar(255),
    firstname varchar(200),
    lastname varchar(200),
    photoURL varchar(300)
)

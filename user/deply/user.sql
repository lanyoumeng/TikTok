CREATE DATABASE IF NOT EXISTS toktik_user;


USE toktik_user;

create table if not exists user
(
    id               bigint unsigned auto_increment
        primary key,
    created_at       datetime(3)  null,
    updated_at       datetime(3)  null,
    deleted_at       datetime(3)  null,
    username         varchar(40)  not null,
    password         longtext     not null,
    avatar           varchar(255) not null,
    background_image varchar(255) not null,
    is_follow        tinyint(1)   not null,
    signature        varchar(255) null,
    constraint username
        unique (username)
);

create table if not exists user_count
(
    id              bigint unsigned auto_increment
        primary key,
    created_at      datetime(3)     null,
    updated_at      datetime(3)     null,
    deleted_at      datetime(3)     null,
    user_id         bigint unsigned null,
    follow_count    bigint          null,
    follower_count  bigint          null,
    total_favorited bigint          null,
    work_count      bigint          null,
    favorite_count  bigint          null,
    constraint fk_user_count_user
        foreign key (user_id) references user (id)
);

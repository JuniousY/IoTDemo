CREATE
DATABASE IF NOT EXISTS `demo`
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE
`demo`;

create table product
(
    id          int auto_increment comment 'id'
        primary key,
    name        varchar(30)                        not null comment '产品名称',
    status      int      default 0                 not null comment '状态 0正常，1已下线，-1删除',
    create_time datetime default CURRENT_TIMESTAMP not null comment '创建时间',
    update_time datetime default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP comment '更新时间'
) comment '产品';

create table device
(
    id          bigint auto_increment comment 'id'
        primary key,
    product_id  int                                   not null comment '产品id',
    name        varchar(20)                           not null comment '设备名称',
    secret      varchar(50)                           not null comment '密码',
    info        varchar(64) default ''                not null comment '设备信息',
    status      int         default 0                 not null comment '状态 0未激活，1已激活，-1删除',
    is_online   int         default 0                 not null comment '是否在线 0下线 1在线',
    create_time datetime    default CURRENT_TIMESTAMP not null,
    update_time datetime    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
) comment '设备';

create index product_id_index
    on device (product_id);

---
INSERT INTO demo.product (id, name)
VALUES ('1', 'prod001');


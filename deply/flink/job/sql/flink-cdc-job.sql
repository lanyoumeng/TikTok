-- 定义MySQL源表
CREATE TABLE MySQLSource (
     id BIGINT,
     name STRING,
     password STRING,
     avatar STRING,
     background_image STRING,
     signature STRING,
     create_at TIMESTAMP,
     update_at TIMESTAMP,
     deleted_at TIMESTAMP,
     PRIMARY KEY (id) NOT ENFORCED
) WITH (
    'connector' = 'mysql-cdc',
    'hostname' = 'mysql',
    'port' = '3306',
    'username' = 'root',
    'password' = '123456',
    'database-name' = 'user',
    'table-name' = 'user',
    'server-id' = '5401'
);

-- 定义Kafka接收表
CREATE TABLE KafkaSink (
   id BIGINT,
   name STRING,
   password STRING,
   avatar STRING,
   background_image STRING,
   signature STRING,
   create_at TIMESTAMP,
   update_at TIMESTAMP,
   deleted_at TIMESTAMP
) WITH (
    'connector' = 'kafka',
    'topic' = 'flink-user',
    'properties.bootstrap.servers' = 'kafka:9094',
    'format' = 'json'
);

-- 从MySQL源表插入数据到Kafka接收表
INSERT INTO KafkaSink
SELECT id, name, password, avatar, background_image, signature, create_at, update_at, deleted_at
FROM MySQLSource;

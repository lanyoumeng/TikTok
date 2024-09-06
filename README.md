# 抖音项目

## 介绍

**项目介绍：**使用Kratos框架开发的微服务短视频平台，用户可通过此平台观看视频和投稿，并进行点赞和评论，关注作者并发送消息

**api文档**：https://console-docs.apipost.cn/preview/599a053c1085c274/9b93e21dcae61608

**技术栈：**Apisix，Docker，Kratos，Etcd，Kafka，Prometheus ，Grafana，Jaeger，ELK，Grpc，Redis，FFmpeg，Gorm，Oss，Jwt

**技术亮点：**

1. ELK+Kafka实现日志的收集和可视化处理
2. 使用Kafka进行流量削峰和异步处理，显著减少了响应时间
3. 利用阿里云Oss存储视频和封面，使用FFmpeg压缩视频和提取封面
4. 通过Jaeger和Prometheus对服务进行追踪监控，使用Grafana提供可视化界面
5. 使用Redis缓存信息，提高性能和响应时间。使用Flink CDC+Kafka保证Redis和Mysql的最终一致性
6. Apisix作为微服务网关，统一入口，对外提供http服务，内部使用grpc通信。基于Apisix实现负载均衡和流量控制
7. 基于Etcd作为服务的注册和发现中心，nacos作为配置中心，使用JWT验证身份、 Bcrypt加密密码、Validator验证数据合法性
8. 采用Docker Compose进行部署，简化环境变量和配置管理，提高可移植性

------



## 目录

以video为例

```
├── flink
│   ├── flink-1.18.0
│   └── job
└── video
    ├── api
    ├── cmd
    │   └── video
    ├── configs
    ├── internal
    │   ├── biz
    │   ├── conf
    │   ├── data
    │   ├── mocks
    │   │   └── mrepo
    │   ├── pkg
    │   │   ├── aliyun
    │   │   ├── errno
    │   │   ├── model
    │   │   └── vkafka
    │   ├── server
    │   └── service
    ├── pkg
    │   ├── token
    │   └── tool
    ├── store
    │   └── video
    └── third_party
```





------



## 技术选型

| 名称                                                         | 用途                                                         | 功能                                            | 端口                                                         |
| ------------------------------------------------------------ | ------------------------------------------------------------ | ----------------------------------------------- | ------------------------------------------------------------ |
| grpc/ kratos                                                 | 微服务服务框架                                               | ✓                                               |                                                              |
| Gorm                                                         | ORM框架                                                      | ✓                                               |                                                              |
| jwt                                                          | 鉴权                                                         | ✓                                               |                                                              |
| Mysql                                                        | 持久化数据库                                                 | ✓                                               |                                                              |
| Redis                                                        | 缓存                                                         | ✓                                               |                                                              |
| Fink cdc                                                     | mysql-redis同步                                              | ✓                                               | 控制面板：8081                                               |
| oss                                                          | 存储                                                         | ✓                                               |                                                              |
| kafka                                                        | 消息队列                                                     | 日志收集✓ 、    发布视频异步✓、评论异步消息异步 | kafdrop(9000:9000)×                                          |
| etcd                                                         | 服务注册中心                                                 | ✓                                               | etcd集群分开部署                                             |
| nacos                                                        | 配置中心                                                     | ✓                                               | http://localhost:8848/nacos 初始账号密码为 nacos/nacos Data Id 必须加.yaml user.yaml |
| swagger                                                      | API文档                                                      |                                                 | apifox直接导出                                               |
| FFmpeg                                                       | 视频压缩、裁剪视频封面                                       | ✓                                               |                                                              |
| APISIX(UI： [Dashboard](https://apisix.apache.org/zh/docs/dashboard/USER_GUIDE/)) | API 网关，专门用于对 API 流量进行管理、路由和流量控制 API 的动态路由、限流、鉴权、监控 | ✓                                               | 9091被Prometheus采集 <br />dashboard （3001:3001）           |
| Prometheus ，grafana                                         | 监控和可观测性软件工具：Prometheus收集系统的性能指标， 然后Grafana创建仪表板进行可视化展示， | ✓                                               | "9090:9090" <br />3000:3000                                  |
| Jaeger                                                       | Tracing链路追踪工具                                          | ✓                                               | ui:16686                                                     |
| Elasticsearch(UI:elastic), Logstash, Kibana                  | Logging日志收集和可视化工具                                  | ✓                                               | ui:5601                                                      |
| docker-compose                                               | 部署                                                         | ✓                                               |                                                              |
|                                                              |                                                              |                                                 |                                                              |
| K8s                                                          | 部署                                                         |                                                 |                                                              |
| CI/CD                                                        | github-action                                                |                                                 |                                                              |



------



## 数据库设计

通过索引加快访问速度

![image-20240906231045066](https://raw.githubusercontent.com/lanyoumeng/Drawing-bed/main/docs/202409062310409.png)

------



## 缓存设计

读：

先从redis读，没有去mysql去取，读取到的数据写入redis+过期时间，否则，从redis返回

写：

使用flinkcdc+分布式锁解决一致性问题

mysql-->flink cdc-->kafka-->订阅处理redis



### 缓存设计

| 服务                               | 数据信息                                                    | 数据类型                | **key**                               | **value**                                                    | 创建时机 | 更新（刷新过期时间）                                         | 删除                           |
| ---------------------------------- | ----------------------------------------------------------- | ----------------------- | ------------------------------------- | ------------------------------------------------------------ | -------- | ------------------------------------------------------------ | ------------------------------ |
| user                               | count <br />较频繁，可以不设置过期时间，或加锁              | hash无对应mysql表       | count::user_id                        | count 计数信息                                               | 查询为空 | 其他服务更新<br />1.视频发布2.点赞<br />3.关注               | 1.过期时间                     |
|                                    | user                                                        | hash                    | user::user_id                         | user信息                                                     | 查询为空 |                                                              | 1.过期时间2.mysql表更新/删除   |
|                                    |                                                             |                         |                                       |                                                              |          |                                                              |                                |
|                                    | video_info<br />视频获赞数和跟随数是rpc调用的               | hash                    | videoInfo::video_id                   | video_info 静态信息type Video struct {    Id       int64     AuthorId int64     Title    string     PlayUrl  string     CoverUrl string } | 查询为空 |                                                              | 1.过期时间2.video表更新/删除   |
| video                              | 系统发布的所有视频 id 用于视频流                            | zset<br />无对应mysql表 | videoAll                              | score: video发布时间戳  <br /> float64 <br />member: video_id |          | 不过期，发布视频时插入                                       | 项目不涉及删除视频             |
|                                    | 用户发布视频的 id <br />用于发布列表                        | set<br />无对应mysql表  | publishVids::user_id                  | video_ids                                                    | 查询为空 | video表创建一条记录：<br />添加对应key-member  video表不会更新<br /> video表删除一条记录：<br />移除对应key-member | 1.过期时间                     |
|                                    |                                                             |                         |                                       |                                                              |          |                                                              |                                |
|                                    | 1.一条点赞信息                                              | string无对应mysql表     | `fav::<videoId>::<userId>`            | true/false                                                   | 查询为空 | favorite表创建/更新记录，操作为：<br />点赞：变为true取消点赞：变为false | 1.过期时间2.favorite表删除记录 |
| **favor**                          | 2.用户点赞视频列表<br />size=用户的点赞数量                 | set<br />无对应mysql表  | `userFavList::<userId>`               | video_id                                                     | 查询为空 | favorite表创建/更新记录，操作为： <br />点赞操作 -->cnt+1<br />取消点赞-->cnt-1 <br />favorite表删除记录,点赞为true-->cnt-1 | 1.过期时间                     |
| 较频繁，可以不设置过期时间，或加锁 | 3.视频的获赞数                                              | string无对应mysql表     | `videoFavCnt::<videoId>`              | int64                                                        | 查询为空 | 点赞操作 -->cnt+1<br />取消点赞-->cnt-1 <br />favorite表删除记录,点赞为true-->cnt-1 | 1.过期时间                     |
|                                    | 4.视频作者的获赞数TotalFavorited                            | string无对应mysql表     | userTotalFavorited::authorId          | int64                                                        | 查询为空 | 点赞操作 -->cnt+1<br />取消点赞-->cnt-1 <br />favorite表删除记录,点赞为true-->cnt-1 | 1.过期时间                     |
|                                    |                                                             |                         |                                       |                                                              |          |                                                              |                                |
| **comment**                        | 评论用户                                                    | zset                    | comment::video_id                     | score: 评论发布时间戳 member: 序列化为json（字符串的话，分隔符会和content冲突）<br />Id        <br />UserId   <br />Content  <br />CreateDate | 查询为空 | comment表创建一条记录：添加对应comment--member  <br />comment表不会更新 <br />comment表删除一条记录：移除对应comment--member | 1.过期时间                     |
|                                    |                                                             |                         |                                       |                                                              |          |                                                              |                                |
|                                    | 关注用户的id                                                | set                     | follow::user_id                       | 关注用户的ids                                                | 查询为空 | relation表创建一条记录：添加对应key-member  relation表不会更新 relation表删除一条记录：移除对应key-member | 1.过期时间                     |
| relation                           | 粉丝的id                                                    | set                     | follower::user_id                     | 粉丝的ids                                                    | 查询为空 |                                                              | 1.过期时间                     |
|                                    | 好友的id <br />通过求关注和粉丝的交集可以得到，不用特地缓存 |                         |                                       |                                                              |          |                                                              |                                |
|                                    |                                                             |                         |                                       |                                                              |          |                                                              |                                |
| message                            | 历史消息                                                    | zset                    | message::history::小user_id+大user_id | score: 消息发布时间戳 member: 序列化为json：user_id<br />target_id<br />content | 查询为空 | message表创建一条记录：添加对应key-member  message表不会更新 messaget表删除一条记录：移除对应key-member | 1.过期时间                     |
|                                    |                                                             |                         |                                       |                                                              |          |                                                              |                                |
|                                    |                                                             |                         |                                       |                                                              |          |                                                              |                                |



//Zset 只存储视频 ID 和时间戳，假如value过大怎么解决

将大 Zset 拆分为多个小 Zset，根据视频发布的时间进行分片存储。例如，可以按天、周、月或其他时间周期进行分片。



### 分布式锁

对favorite redis添加了etcd分布式锁，防止缓存击穿

在internal/pkg/favkafka/redis.go

------





## 服务实现

基于Etcd作为服务的注册和发现中心，nacos作为配置中心，

rpc内部使用rpc调用获取相关信息



### User

1. 使用JWT验证身份、 Bcrypt加密密码、Validator验证数据合法性



- 因为创建和登录用户时，userId是没有的，所以都是用username查询  加索引

- is_follow  在video服务查找  user服务不提供

### video

#### 发布

1. 通过kafka异步发布视频，减少响应时间
2. filetype库校验文件类型，使用ffmpeg进行视频压缩和封面提取
3. 阿里云OSS存储视频和封面

#### 视频流

1. 使用zset存储视频id，方便根据时间获取视频信息
2. sync.Map+errgroup.WithContext(ctx) 利用协程并发组装视频流



### favorite

#### 点赞

1. 对缓存添加了etcd 分布式锁，防止缓存击穿

2. action_type

   1-点赞，2-取消点赞



### comment

1. action_type

   1-发布评论，2-删除评论



### relation

1. action_type

   1-关注，2-取消关注

2. 获取关注列表
   获取关注用户的ids 然后rpc获取用户信息
   再获取每个用户的is_follow字段 即登录用户是否关注了该用户

3. 获取好友列表
   好友列表，关注和粉丝的交集

   

### message

1. 前端是轮询获得消息
2. 时间顺序获取消息记录
3. 

------




## app

使用前端提供的app

app bug：

1. Avatar用户头像和BackgroundImage 用户个人页顶部大图默认
2. 视频流刷新不及时

------

## 消息队列



| 服务模块  | 用途                           |
| --------- | ------------------------------ |
| video     | 异步上传视频，缩短响应时间     |
| sociality | 异步执行社交操作，缩短响应时间 |
| chat      | 异步发送聊天消息，缩短响应时间 |



### topic

一共有3类topic

1.日志

2.redis同步

3.video的视频发布





------



## ELK

- 通过logstash收集和聚合微服务的日志数据传输给kafka

![image-20240906232118199](https://raw.githubusercontent.com/lanyoumeng/Drawing-bed/main/docs/202409062321544.png)

------



## 测试

1.自动生成mock

```
//go:generate mockgen -destination=../mocks/mrepo/video.go -package=mrepo . VideoRepo
```

------



## 遇到问题

1. v.log.Error 用于err    v.log.Debug用于调试

2. 需要关闭防火墙

3. 异步写入kafka时，ctx不要用上层传入的，  因为上层直接结束了，ctx相当于到期了

   重新context.Background()

4. 视频投稿接口的post的body是multipart/form-data类型，  而 gRPC 的 `protobuf` 不支持直接处理这种格式。生成的video_http.pb.go文件中的对应函数是不能用的，要进行修改

## 

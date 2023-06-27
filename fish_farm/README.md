#### `golang`框架

- gin

#### 微服务模块

**本项目微服务模块分为两种**

**对外暴露接口**：前端可通过网关来调用，对前端发来的请求进行处理

**内部调用**：前端无法通过网关来直接调用，而是被其它微服务模块通过`rpc`来调用，可以与数据库进行交互

**对外暴露接口**

- goods-web：商品微服务
- order-web：订单微服务
- oss-web：文件操作微服务
- user-web：用户微服务
- userop-web：用户操作微服务

**内部调用**

- goods_srv：商品微服务
- inventory_srv：库存微服务
- order_srv：订单微服务
- user_srv：用户信息微服务
- userop_srv：用户操作微服务



##### 微服务使用的技术栈

- 注册中心：服务注册、服务健康检查
  - consul
- 配置中心：统一管理模块配置信息
  - nacos
- 网关：服务发现、动态路由、负载均衡、限流
  - koog
- 熔断：服务被调用时出错进行熔断降级
  - sentinel
- 服务通信：微服务之间的调用
  - grpc
- 链路追踪：用于服务排查
  - jaeger
- 消息队列：消息分发、应用解耦
  - rocketmq



##### 服务编写使用的库

- go-password-encoder：密码加密与解密
- jwt：加密数据作为`token`使用
- zap：将日志转为文件进行存储
- viper：配置文件管理
- base64Captcha：图片验证码
- redsync：基于`redis`分布式事务锁
- gorm：将`msyql`中的表映射为`go`中的`struct`
- grpc：调用`.protoc`生成的`go`文件
- validator：表单验证
- es：与`elasticsearch`进行交互
- alipay：集成`支付宝`支付



##### 第三方服务

- oss：阿里云`oss`存储服务
  - 用来存储图片
- sms：阿里云`sms`短信发送服务
  - 用于用户注册与登入
- alipay：阿里云支付服务
  - 可以使用`支付宝`进行支付



##### 数据库选择

- mysql：存储用户商品等服务模块数据
- redis：存储 `短信验证码`、token、session
- elasticsearch：存储商品索引后的消息，用于搜索优化
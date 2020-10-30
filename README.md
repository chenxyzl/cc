# cc
for cc chat test

#运行方法
启动bin目录下对应的二进制文件，访问http://127.0.0.1:33001/

#cmd
入口

#analssis
词频分析

#framework
核心框架[本测试比较简单，时间也比较紧，所以本处目前只有seesion和uuid]

#logic
业务逻辑层[目前包含用户管理和聊天室管理，预留了一些扩展接口，时间教紧，未完善]

#msg_dispatcher
消息分发层[目前只包含消息包和事件，分发还在logic里，未完善]

#util
常用工具包[目前只包含了敏感词过滤--过滤算法几年前写的，有点看不懂了～～]

#质量管理
目前只有单元测试-go test[完整的需要一整套cicd支持,嵌入到jerkins/gitlab runner等]

#扩展性
预留了很多,未完善的也有很多[如rpc注册分发，抽出server(目前UserManger承担了这个职能)， gate，持久化-分布式注册发现等--计划做，被别的事情耽误了]

#客户端技术
cocos create 2.2

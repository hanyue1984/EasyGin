# 简化Gin框架(EasyGin)

基于Gin框架清晰目录，路由分离，增加redis以及ORM的支持
使用以下技术库

1.Go基础 [教程地址](https://www.topgoer.com/go基础/)

2.Gin [教程地址](https://www.topgoer.com/gin框架/简介.html) 

3.GORM [教程地址](https://www.topgoer.com/数据库操作/gorm/gorm介绍.html) 

4.搭配后台 [下载地址](https://gitee.com/icecoldmoon/easyGinAdmin.git) 

搭建好后需要用后台创建一个管理员，然后修改数据库中的user表把admin设置为true roles设置成{root}这样就拥有了最高权限能得到所有的后台界面

后台中有初始化的数据主要数据是menu跟menu_mate表中的数据否则看不到菜单

工具简单封装了Redis
使用用例
```
//声明一个工具包中的redis客户端
var UsersRedis tools.RedisClient

//连接配置包中的要用的redis服务器配置该接口就算调用多次也只会连接一个类似单例这里指的是User如果你换个名字还会创建一个连接
UsersRedis.Connect("User", Config.AppConfig.RedisCommon)

//获取数据并且转化成结构体
result := UsersRedis.Get(ctx, key, models.Users{})
```
## 目录说明
```lua
启动默认是dev环境如果自己需要添加则启动时需要 -env
主目录
├── app.go -- main入口
└── app -- 应用逻辑入口
     ├── config 配置文件
     ├── controllers 控制器目录
     ├── models 模型目录
     ├── routes 路由目录
     └── services 服务目录
└── common -- 公共组件包
     └── middleware 公共中间件
     └── tools 公共工具组件
            └──redis 简单封装
```
## 联系方式
有不明白的或者有什么意见欢迎联系我 WX:icecoldmoon QQ:5101111
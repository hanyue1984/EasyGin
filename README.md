# 简化Gin框架(EasyGin)

基于Gin框架清晰目录，路由分离，增加redis以及ORM的支持
使用以下技术库

1.Gin 教程地址 https://www.kancloud.cn/gopher_go/go/570005

2.GORM 教程地址 https://www.topgoer.com/数据库操作/gorm/gorm介绍.html

工具简单封装了Redis
使用用例
```
var UsersRedis tools.RedisClient
UsersRedis.Connect("User", Config.AppConfig.RedisCommon)
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
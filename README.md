## Nomo

取自**No**t only Me**mo**。通过Lark、微信将思考碎片记录到Notion或者飞书文档

## 编译部署
### 编译

修改`cmd/nomo/.env`设置正确的MySQL地址，然后运行build.sh即可
```
./build.sh [platform] [env]
```
**platform：** go支持夸平台编译，platform默认是编译为当前平台，如果需要linux平台，则指定为linux即可
***env：** 用来选择打包到output目录选择的`.env`文件，`cmd/nomo`目录下存在几个`.env`文件，根据指定的名字打包对应的文件到ouput目录

### 运行
编译后会生成一个`output`目录，运行只需要一个二进制文件和`.env文`件，全部都在`bin`目录下，修改`.env`文件设置MySQL和相关的环境，然后直接运行二进制即可
```
output
├── bin
│   ├── .env  # 配置文件，MySQL等配置
│   └── nomo  # 可执行程序
├── conf   
│   ├── openhex.crt # TLS证书，用于https
│   └── openhex.key # TLS私钥，用于https
├── run.sh    # 启动脚本
└── run_wx.sh # 微信云启动脚本
```
进入bin目录，直接运行即可
```
./nomo
```

### 部署
理论上部署取决于环境，当前仓库提供了两种方式：[微信云](https://cloud.weixin.qq.com/cloudrun)和Ansible部署
- 微信云托管是个Docker环境需要Dockerfile，当前仓库提供了一份，需要在微信云部署的只需要fork一份仓库，然后修改`cmd/nomo/.env`为自己的微信云托管环境的配置，就可以直接部署
- Ansible部署更加通用一些，可以在任何云主机或者物理机环境部署，`deploy`目录实现了一个ansible部署的方式，通过sysmted托管进程，确保机器重启自动启动服务

不想自己维护的同学，也可以直接使用我在[腾讯云](https://cloud.tencent.com/)部署的一套服务：[https://nomo.openhex.cn/api/v1](https://nomo.openhex.cn/api/v1)。使用方式可以参考：[如何使用飞书机器人打造一个私人的Flomo?](https://openhex.cn/%E5%A6%82%E4%BD%95%E4%BD%BF%E7%94%A8%E9%A3%9E%E4%B9%A6%E6%9C%BA%E5%99%A8%E4%BA%BA%E6%89%93%E9%80%A0%E4%B8%80%E4%B8%AA%E7%A7%81%E4%BA%BA%E7%9A%84flomo)

## RoadMap
- [x] 不同租户Lark机器人支持
- [x] 支持多种Notion页面主题，比如flat类型以及database类型
- [x] 支持微信云部署
- [x] 支持腾讯云部署
- [x] 支持飞书Doc存储Memo
- [ ] 支持微信订阅号发送Memo

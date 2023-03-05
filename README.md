<h1 align="center">Welcome to Nomo 👋</h1>
<p>
</p>

> Nomo(Not only Memo)，通过Lark、微信将思考碎片记录到Notion或者飞书文档

## 编译部署
### 编译

修改`cmd/nomo/.env`设置正确的MySQL地址，然后运行build.sh即可
```
./build.sh [platform] [env]
```
**platform：** go支持夸平台编译，platform默认是编译为当前平台，如果需要linux平台，则指定为linux即可
***env：** 用来选择打包到output目录选择的`.env`文件，`cmd/nomo`目录下存在几个`.env`文件，根据指定的名字打包对应的文件到ouput目录

### 运行
编译后会生成一个`output`目录，运行只需要一个二进制文件和`.env``文件，全部都在`bin`目录下，修改`.env`文件设置MySQL和相关的环境
`conf`目录下的`crt`和`key`证书为SSL证书，如果需要以`https`模式启动请将自己域名对应的证书放到该目录并且修改`.env`
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

不想自己维护的同学，也可以直接使用我在[腾讯云](https://cloud.tencent.com/)部署的一套服务：[https://nomo.openhex.cn/api/v1](https://nomo.openhex.cn/api/v1)。使用方式可以参考：[如何使用飞书机器人打造一个私人的Flomo?](https://blog.openhex.cn/posts/35d22c04-5518-4871-9812-832af9e8d5fa)

## RoadMap
- [x] 不同租户Lark机器人支持
- [x] 支持多种Notion页面主题，比如flat类型以及database类型
- [x] 支持微信云部署
- [x] 支持腾讯云部署
- [x] 支持飞书Doc存储Memo
- [x] 支持微信订阅号发送Memo

---
## Build and Deploy
### Build
Modify cmd/nomo/.env to set the correct MySQL address, and then run build.sh.
```
./build.sh [platform] [env]
```
**platform: **Go supports cross-platform compilation. The default platform is the current platform. If you need to compile for Linux, specify linux.
**env: **Used to select the corresponding .env file to be packaged into the output directory. There are several .env files under the cmd/nomo directory. Package the corresponding file to the output directory according to the specified name.

### Run
After compilation, a output directory will be generated. To run, you only need an executable file and an .env file, all of which are in the bin directory. Modify the .env file to set up MySQL and related environments. 
The certificates in the crt and key directories under the conf directory are SSL certificates. If you need to start in https mode, put the certificate corresponding to your domain name in the directory and modify the .env file.

```
output
├── bin
│   ├── .env  # Configuration file, MySQL and other configurations
│   └── nomo  # Executable program
├── conf   
│   ├── openhex.crt # TLS certificate, used for https
│   └── openhex.key # TLS private key, used for https
├── run.sh    # Start script
└── run_wx.sh # WeChat Cloud start script
```
Go to the bin directory and run it directly.
```
./nomo
```
### Deployment
In theory, deployment depends on the environment. The current repository provides two ways: WeChat Cloud and Ansible deployment.

- [WeChat Cloud](https://cloud.weixin.qq.com/cloudrun):  WeChat Cloud hosting is a Docker environment that requires a Dockerfile. The current repository provides one. If you need to deploy to WeChat Cloud, simply fork a repository, then modify cmd/nomo/.env to your own WeChat Cloud hosting environment configuration, and you can deploy directly.

- [Ansible](https://www.ansible.com/): Ansible deployment is more universal and can be deployed in any cloud host or physical environment. The deploy directory implements an Ansible deployment method, hosting the process through systemd to ensure that the machine automatically starts the service after a restart.

For those who do not want to maintain it themselves, you can also use a service I deployed on Tencent Cloud: https://nomo.openhex.cn/api/v1. For usage, please refer to: [How to use a Feishu robot to create a private Flomo?](https://blog.openhex.cn/posts/35d22c04-5518-4871-9812-832af9e8d5fa)

## RoadMap
- [x] Support for different tenants' Lark robots
- [x] Support for multiple Notion page themes, such as flat types and database types
- [x] Support for WeChat Cloud deployment
- [x] Support for Tencent Cloud deployment
- [x] Support for Feishu Doc storage Memo
- [x] Support for WeChat Subscription Number sending Memo

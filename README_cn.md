# Go OSS Server
迷你的OSS服务器，本地测试用

## 环境准备
* 编译时用的Go版本: `Golang 1.13`
* `OSS_CONFIG` OSS服务器配置文件路径

## 用法 (用 `testenv` 目录来实践一下)
* 复制模板配置文件
~~~
> cp env/ testenv/
~~~
* 编译可执行文件，然后进入目录 `testenv`
~~~
> PROJECT_PATH='github.com/dormao/go-oss-server/srv'
> OUTPUT_FILE='testenv/oss_server'
> go build -o $OUTPUT_FILE $PROJECT_PATH
> cd testenv/
~~~
* 编辑OSS配置文件
~~~
> vi config_template.yaml
~~~
* 设置环境变量
~~~
> cp env.template env.sh
> chmod a+x env.sh
> source env.sh
~~~
* 直接运行
~~~
> chmod a+x oss_server
> ./oss_server
~~~

# 用HTTP请求，与OSS服务器交互
## 用表单上传资源 ( 模拟地址: `localhost:8080` )
必须根据以下表单格式去创建，否则会返回 400(Bad Request) 或者 401(Unauthorized)

|  表单字段 | 描述            |
|-------------|-----------------|
| bucket      | 桶名 |
| object      | 对象名  |
| file        | 二进制文件    |
| accesskey   | Access key      |
| secret      | Access secret   |

然后服务器会返回以下Json

| Json 字段     | 类型                                              |
|----------------|---------------------------------------------------|
| code           | 返回码\(基本上与HTTP状态码一样\) |
| msg            | 错误信息\(不出错的时候返回空字符串\)                    |
| result         | 上传结果                                 |
| result\.object | 上传的对象名                               |
| reuslt\.url    | 资源的URL                               |

## 从OSS下载资源 ( 模拟地址: `localhost:8080` )
GET: localhost:8080/`桶`/`对象名`

## 对象名
你可以用类似这些的对象名来模拟出一些看起来有目录的文件

`my-object-file.yml`

`electron/my-object-file.yml`

`avatars/dormao.png`

## 许可证
[MIT License](https://www.mit-license.org/)

## 其他
[English](README.md)

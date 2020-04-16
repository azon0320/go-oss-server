# Go OSS Server
Tiny oss server for local usage

## Environment Preparation
* Programing Version: `Golang 1.13`
* `OSS_CONFIG` The OSS server config path

## Usage (`testenv` as runtime path)
* Copy the template path
~~~
> cp env/ testenv/
~~~
* Build program and go to `testenv`
~~~
> PROJECT_PATH='github.com/dormao/go-oss-server/srv'
> OUTPUT_FILE='testenv/oss_server'
> go build -o $OUTPUT_FILE $PROJECT_PATH
> cd testenv/
~~~
* Edit Template Config
~~~
> vi config_template.yaml
~~~
* Set Environment Variables
~~~
> cp env.template env.sh
> chmod a+x env.sh
> source env.sh
~~~
* Tiny run
~~~
> chmod a+x oss_server
> ./oss_server
~~~

# Interact with OSS Server by http request
## Upload File by form ( mock address: `localhost:8080` )
You must POST form data with these fields

|  Form Field | Type            |
|-------------|-----------------|
| bucket      | The Bucket Name |
| object      | The object key  |
| file        | Binary file     |
| accesskey   | Access key      |
| secret      | Access secret   |

Then the server will returns JSON

| Json Field     | Type                                              |
|----------------|---------------------------------------------------|
| code           | return code\(default sync with http status code\) |
| msg            | error message\(default blank\)                    |
| result         | the result object                                 |
| result\.object | uploaded object key                               |
| reuslt\.url    | public resource url                               |

## Download file ( mock address: `localhost:8080` )
GET: localhost:8080/`BucketName`/`objectkey`

## Object Keys
You can use object keys like these as it seems like the file in directory

`my-object-file.yml`

`electron/my-object-file.yml`

`avatars/dormao.png`

## License
[MIT License](https://www.mit-license.org/)

## Others
[中文](README_cn.md)

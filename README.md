# License

GoLaravel 框架是根据 [MIT 许可证](https://opensource.org/licenses/MIT) 授权的开源软件。

# 1.运行配置

- ## 拉取代码库并安依赖包
    ````
    git clone "代码库地址" && go mod tidy
    ````
- ## 复制env配置文件 并修改相应的配置（配置参考config目录下的文件）
    ````
    cp .env.example .env
    ````
- ## 运行
    ````
    go run main.go
    ````
# 2.指令介绍

- ## 开启api服务的指令（也为默认指令）
    ````
    go run main.go

    或者

    go run main.go api

    ````
- ## 数据库迁移指令 (迁根目录下的database/migrations下的所有文件)

  - ### 迁移生成对应的表
    ```
    go run main.go migrate up
    ```
  - ### 迁移回滚表上一个操作
    ```
    go run main.go migrate down
    ```
  - ### 回溯所有迁移
    ```
    go run main.go migrate reset
    ```   
  - ### reset命令 后重新之前 migrate up 命令
    ```
    go run main.go migrate fresh
    ```
- ## 自动生成指令

  - ### 根据数据库的表自动生成模型文件指令
    ```
    go run main.go generate model [参数为表名称 可有零个或者多个用空格分开 不传参数代表库中所有的表]
    ```

  - ### 自动生成应用密钥指令
    ```
    go run main.go generate key
    ```

- ## 文件生成指令

  - ### 生成迁移文件
    ```
    go run main.go make migration [需要传一个参数为文件名]
    ```
  - ### 生成命令行（记得生成的指令要加入到app/cmd/kernel.go 下的InintCmd里面才能生效）
    ```
    go run main.go make cmd [需要传一个参数为文件名]
    ```
  - ### 生成模型文件
    ```
    go run main.go make model [需要传一个参数为文件名]
    ```
- ## 生成接口文档指令 (访问链接为 http://host:port/swagger/index.html)
    ````
    swag init 
    ````
# 3.异步任务模块（基于redis实现）

- ## 生成异步任务文件
    ````
    go run main.go make job [需要传一个参数为文件名]
    ````

- ## 异步队列的调用
    ````
    #示例代码 NewExampletJob为 app/job/ 下的自定义任务
    asynq.Delivery(job.NewExampletJob(job.ExamplePayload{
        UserId: 11111,
    }))
    ````
- ## 异步队列处理命令
    ````
    go run main.go queue server
    ````

# 4.服务启动、重启、关闭指令（linux）

- ## 开启服务
    ````
    go build -o <编译打包的文件服务名>
    nohup ./<编译打包的文件服务名> &
    ````

- ## 关闭服务
    ````
    ps -ef | grep <编译打包的文件服务名> | awk -F' ' '{print$2}' | xargs kill -9
    ````
- ## 重启服务
    ````
    go build -o <编译打包的文件服务名> > restart.txt 2>&1
    nohup ps -ef | grep ./<编译打包的文件服务名> | awk -F' ' '{print$2}' | xargs kill -1  > restart.txt 2>&1 &
    ````

# Command

## skyflow

skyflow 命令行工具是用来启动 skyflow 服务的最主要的方式。

1. 启动skyflow api & dispatcher 服务

```shell
skyflow start api dispatcher --port 8080 --config ./trpc.yaml
skyflow syncschema --config ./trpc.yaml
```

* api : 启动skyflow http api 网关服务
* dispatcher: 启动skyflow 后端调度执行服务


## skyflowcli

skyflowcli 是skyflow 命令行工具的客户端，用来和skyflow api 服务进行交互。
包括:

1. 工作流模板相关
    1. 创建模板、更新模板、删除模板
    2. 查询模板列表
    3. 查询模板详情
2. 任务相关
    1. 创建任务
    2. 查询任务列表
    3. 查询任务详情
    4. 查询任务日志
    5. 查询任务状态
    6. 查询任务结果
    7. 查询任务执行历史
    8. 查询任务执行历史详情
    9. 查询任务执行历史日志


```Shell
   # 创建新任务
    skyflowcli startexecution --statemachineuri="statemachine:unitest/deploy_cluster" --title="" --description="" --input="{}" --executionid="1234567890"
    # 查询新任务
    skyflowcli describeexecution --executionid="1234567890"

    # 查询任务列表
    skyflowcli listexecution --statemachineuri="statemachine:unitest/deploy_cluster" --status="running" --limit=10 --offset=0

```



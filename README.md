# skyflow_backbend
skyflow is a workflow based on AWS  Amazon States Language (ASL)



# repository code struct
代码结构
```
+ proto             // api proto file defintion
+ gen               // protoc generate code
    |- pb           // generated protobuf files. *.pb.go
    |- apidoc       // generated api doc files . swagger.json / openapi.json
+ docs              // docs for this repository
+ schemas           // schemas for statemachine language
+ cmd               // build command defintion
    |- skyflow      // skyflow server
    |- skyflowcli   // skyflow cli
+ config            // config struct defintion
+ workflow          // main service part for workflow
    |- parser       // statemachine parser
    |- executor     // execution executor
    |- tempate      // template management
+ server            // server defintion
    |- apiserver    // api server
    |- dispatcher   // dispatcher server
+ examples          // example test file
```

调用依赖图

```mermaid
flowchart TD
    A1[cmd.skyflow_server]
    A2[cmd.skyflow_cli]
    B1(server.apiserver)
    B2(server.dispatcher)
    A1 --> B1
    A1 --> B2
    C1(workflow.parser)
    C2(workflow.template)
    C3(workflow.executor)
    B1 --> C2
    B1 --> C1
    B1 --> C3
    B2 --> C1
    B2 --> C3
    D1(workflow.repository)
    C1 --> D1
    C2 --> D1
    C3 --> D1
    F1(pkg)
    C1 --> F1
    C2 --> F1
    C3 --> F1
    D1 --> F1
    P1(proto.gen)
    A2 --> P1
```
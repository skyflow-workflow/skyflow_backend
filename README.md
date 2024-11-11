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

调用流程图

```mermaid
gantt
    section Section
    Completed :done,    des1, 2014-01-06,2014-01-08
    Active        :active,  des2, 2014-01-07, 3d
    Parallel 1   :         des3, after des1, 1d
    Parallel 2   :         des4, after des1, 1d
    Parallel 3   :         des5, after des3, 1d
    Parallel 4   :         des6, after des4, 1d
```
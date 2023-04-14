# cicdgo

## Build
```
go1.15.15 mod tidy

go1.15.15 build -o bin/cicd-server  cmd/cicd-server/main.go
go1.15.15 build -o bin/cicd-console cmd/cicd-console/main.go
go1.15.15 build -o bin/cicd-agent   cmd/cicd-agent/main.go
go1.15.15 build -o bin/cicd-proxy   cmd/cicd-proxy/main.go
```

## Architecture flow
```
                                                                              ws         +----------------+  ftp/s3/http..
                                                                       +-----------------+    AgentCI     +-----------------+
                                                                       |                 +----------------+                 |
                                                                       |                                                +---+-----------+
                                                                       |                                                |               |
                                         ws             +--------------+--+     ws                                  +---+ ftp/s3/nas..  |
                               +------------------------+                 +-------------------+                     |   |               |
                               |                        |    proxy        |                   |                     |   |               |
                               |                        |                 +----------+   +----+-----------+         |   +---------------+
                   +-----------+-----+                  |                 |          |   |   AgentCD      +---------+
                   |                 |                  +---------------+-+          |   +----------------+ s3/http..
      +----------- |     server      |                                  |            |
      |            |                 +--------+                         |            |
      |            +-----------------+        |                         |            |ws
      |                    ^                  |                         |            |
      v                    |                  |                         |            |
+------------+             |                  |                         |            |
|            |             | http             |                         |            |
|    mysql   |             |                  |                         |ws          |   +-----------------+  ftp/s3/http..
|            |             |                ws|                         |            +---+   AgentCD       +-------------------+
+------------+             |                  |                         |                +-----------------+                   |
      ^                                       |                         |                                                      |
      |            +------------------+       |                         +---------+                                            |
      |            |                  |       |                                   |                                            |
      +----------- |     console      |       |         +-----------------+       |                                            |
                   |                  |       |         |                 |       |                                      +-----+--------+
                   +------------------+       +---------+    proxy        |       |                                      |              |
                                                        |                 |       |                                +-----+ ftp/s3/nas.. |
                                                        |                 |       |                                |     |              |
                                                        +-----------------+       |      +------------------+      |     |              |
                                                                                  |      |                  +------+     +--------------+
                                                                                  +------+    AgentCD       |    ftp/http..
                                                                                         |                  |
                                                                                         +------------------+
                                                                                         
```

## Pipeline config sample
```json
{
    "stageCI":{
        "envs":{
            "scmrepo_name1": "git|github|tangr/webconsole",
            "scmrepo_name2": "git|github|tangr/webconsole2",
            "aaa":[
                "a",
                "b",
                "c"
            ],
            "bbb":[
                "a",
                "b",
                "c"
            ]
        },
        "script":"ci_script_name1",
        "args":""
    },
    "stageCD":{
        "script":"cd_script_name1",
        "envs":{
            "env1":"env1",
            "env2":"env2"
        },
        "args":""
    }
}
```

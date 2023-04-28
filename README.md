# cicdgo

## Build

```
# https://go.dev/doc/install

go1.15.15 mod tidy

go build -o bin/cicd-server  cmd/cicd-server/main.go
go build -o bin/cicd-console cmd/cicd-console/main.go
go build -o bin/cicd-agent   cmd/cicd-agent/main.go
go build -o bin/cicd-proxy   cmd/cicd-proxy/main.go
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
    "stageCD": {
        "args": "",
        "envs": {
            "env1": [
                "env1"
            ],
            "env2": [
                "env2"
            ],
            "PKGRDM": ""
        },
        "script": "nginx-conf-cd"
    },
    "stageCI": {
        "args": "1111args",
        "envs": {
            "aaa": [
                "a",
                "b",
                "c"
            ],
            "bbb": [
                "a",
                "b",
                "c"
            ],
            "REPO_nginx_conf": "git|gitlab|144",
            "REPO_nginx_conf_ssl": "git|gitlab|144"
        },
        "script": "nginx-conf-ci"
    }
}
```

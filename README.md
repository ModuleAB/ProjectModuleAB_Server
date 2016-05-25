Module AB Server
=====

An easy way to backup and archive @ aliyun

Requirements
----
1. bower >= 1.7.6
2. go >= 1.6
3. beego framework >= 1.6.1
4. redis >= 3.0.0
5. mysql >= 5.6

Build
----

```bash
mkdir -p project/src
cd project/src
git clone --recursive https://github.com/ProjectModuleAngelaBaby/ProjectModuleAB_Server moduleab_server
cd ProjectModuleAB_Server
go get -v
go get github.com/beego/bee
export PATH="$GOPATH/bin:$PATH"
make
```

Then use the `moduleab_server.tar.gz` to deploy anywhere you want.

Configuration
----

```ini
appname = moduleab_server
httpport = 7001
# run mode has following options:
# dev: development mode
# deb: debug mode, log will be HUGE!
# initdb: create data in database, DONT USE if you already have data in database.
# proc: production mode.
runmode = dev
autorender = false
copyrequestbody = true
EnableDocs = false

EnableAdmin = false
AdminHttpAddr = "localhost"
AdminHttpPort = 8088

sessionon = true

loginkey = 61oETzKXQAGaYdkL5gEmGeJJFuYh7EQnp2XdTP1o

logFile = "logs/moduleab_server.log"
pidFile = "logs/moduleab_server.pid"

[database]
mysqluser = "ModulesAB"
mysqlpass = "ModulesAB"
mysqlurl = "127.0.0.1:3306"
mysqldb   = "ModuleAB"
mysqlprefex = ""

[aliapi]
apikey= "TestAbcd" # Ali api key
secret="TestAAA"   # Ali api secret
oasport=80
oasusessl=false

[redis]
host = "127.0.0.1:6379"
password = ""
key = "ModuleAB"

[websocket]
timeout=10
pingperiod=5

# policyrun use cron-like syntax: "s m h dom mon dow"
[misc]
checkoasjobperiod=10
policyrun="0 * * * * 1"
```

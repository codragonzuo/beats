## Go Package Example

```
git clone https://github.com/codragonzuo/beats.git



git clone https://github.com/codragonzuo/go-example.git

export http_proxy=http://127.0.0.1:10808/

export https_proxy=http://127.0.0.1:10808/

export http_proxy=

export https_proxy=


```

### 修改 ./git/config,这样git才能提交
```
[core]
        repositoryformatversion = 0
        filemode = true
        bare = false
        logallrefupdates = true
[remote "origin"]
        fetch = +refs/heads/*:refs/remotes/origin/*
#       url = https://github.com/codragonzuo/go-example.git
        url = ssh://git@github.com/codragonzuo/go-example.git

[branch "master"]
        remote = origin
        merge = refs/heads/master
```

### go run main.go运行hello/main.go


### Go Enviroment

```
[root@node1 hello]# go env
GO111MODULE="auto"
GOARCH="amd64"
GOBIN=""
GOCACHE="/root/.cache/go-build"
GOENV="/root/.config/go/env"
GOEXE=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GONOPROXY=""
GONOSUMDB=""
GOOS="linux"
GOPATH="/root/go-example/hello"
GOPRIVATE=""
GOPROXY="https://goproxy.cn,direct"
GOROOT="/usr/local/go"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/usr/local/go/pkg/tool/linux_amd64"
GCCGO="gccgo"
AR="ar"
CC="gcc"
CXX="g++"
CGO_ENABLED="1"
GOMOD="/root/go-example/go.mod"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build372689315=/tmp/go-build"
```

### 查找和替换
```
sed -i "s/github.com\/elastic\/beats\/v7/github.com\/codragonzuo\/beats/g" `grep -rl "github.com" ./`
```



```
import yaml
import os
file = open('filebeat.yml', 'r', encoding="utf-8")
file_data = file.read()
file.close()
data=yaml.load(file_data)
data['filebeat.config.inputs']
{'enabled': True, 'path': 'myconfig.yml', 'reload.enabled': True, 'reload.period
': '10s'}
data['filebeat.config.inputs']
{'enabled': True, 'path': 'myconfig.yml', 'reload.enabled': True, 'reload.period
': '10s'}
data['filebeat.config.inputs']['path']
```



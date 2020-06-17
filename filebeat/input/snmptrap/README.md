## Snmp Trap
 
```
import socket

address = ('127.0.0.1', 31500)
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
s.bind(address)

while True:
    data, addr = s.recvfrom(2048)
    if not data:
        print "client has exist"
        break
    print "received:", data, "from", addr

s.close()


import socket

address = ('127.0.0.1', 31500)
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

while True:
    msg = raw_input()
    if not msg:
        break
    s.sendto(msg, address)

s.close()

```

## send udp message
```
import socket
address = ('127.0.0.1', 9000)
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
msg = 'hello'
s.sendto(msg, address)
```

## send kafka message
```
from kafka import KafkaProducer
producer = KafkaProducer(bootstrap_servers='192.168.20.45:6667')
msg = "Hello World".encode('utf-8')  # 发送内容,必须是bytes类型
producer.send('snmp', msg)  # 发送的topic为test
```




## go环境变量

```
yum install  golang-go

[root@ecs-centos-7 ~]# go env -w GO111MODULE=on
[root@ecs-centos-7 ~]# go env -w GOPROXY=https://goproxy.cn,direct
[root@ecs-centos-7 ~]# ls
8nodes.png  beats  beats-back-0608  classification_pr.png  go  index.html  master.zip  mydask.png  neo4j-4.0.2-1.noarch.rpm  phantomjs-2.1.1-linux-x86_64.tar.bz2  sigma  sigma.wiki  v2ray-linux-64

```

```
[root@ecs-centos-7 ~]# go env
GO111MODULE="on"
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
GOPATH="/root/go"
GOPRIVATE=""
GOPROXY="https://goproxy.cn,direct"
GOROOT="/usr/lib/golang"
GOSUMDB="off"
GOTMPDIR=""
GOTOOLDIR="/usr/lib/golang/pkg/tool/linux_amd64"
GCCGO="gccgo"
AR="ar"
CC="gcc"
CXX="g++"
CGO_ENABLED="1"
GOMOD="/dev/null"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build616420941=/tmp/go-build -gno-record-gcc-switches"
[root@ecs-centos-7 ~]# ls
```
## 代码下载
```
git clone https://github.com/codragonzuo/beats.git

git clone https://github.com/codragonzuo/beats.git

git clone https://github.com/codragonzuo/go-example.git

export http_proxy=http://127.0.0.1:10808/

export https_proxy=http://127.0.0.1:10808/

export http_proxy=

export https_proxy=
```




------------------------------------------------------------------------------------------------
## python3 ssl 设置

首先要安装openssl，因为python3需要引用openssl模块，但是centos需要的openssl版本最低为1.0.2，但是centos 默认的为1.0.1，所以需要重新更新openssl 。
```shell
yum install -y zlib zlib-dev openssl-devel sqlite-devel bzip2-devel libffi libffi-devel gcc gcc-c++
wget http://www.openssl.org/source/openssl-1.1.1.tar.gz
tar -zxvf openssl-1.1.1.tar.gz
cd openssl-1.1.1
./config --prefix=$HOME/openssl shared zlib
make && make install
echo "export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$HOME/openssl/lib" >> $HOME/.bash_profile
source $HOME/.bash_profile

yum install zlib-devel bzip2-devel openssl-devel ncurses-devel sqlite-devel readline-devel tk-devel gcc make libffi-devel -y
yum install python-pip -y

wget https://www.python.org/ftp/python/3.7.0/Python-3.7.0.tgz
tar -zxvf Python-3.7.0.tgz
./configure prefix=/usr/local/python3  --with-openssl=$HOME/openssl
make && make install
ln -s /usr/local/python3/bin/python3.7 /usr/bin/python3.7 
ln -s /usr/local/python3/bin/pip3.7 /usr/bin/pip3.7

libcrypto.so -> libcrypto.so.1.0.1e

libcrypto.so.1.1 -> /usr/local/openssl/lib/libcrypto.so.1.1

```


```
cd Python-3.6.2
./configure --with-ssl
make
sudo make install
```
------------------------------------------------------------------------------------------------
## ssh设置
```shell
[root@node1 filebeat]# eval "$(ssh-agent -s)"
Agent pid 8296
[root@node1 filebeat]# cd /root/.ssh
[root@node1 .ssh]# ls
authorized_keys  ida_rsa_codragonzuo  ida_rsa_codragonzuo.pub  id_rsa  id_rsa.pub  known_hosts
[root@node1 .ssh]# ssh-add  ./id_rsa_codragonzuo
[root@node1 .ssh]# ssh-keygen -t rsa -b 4096 -C "codragon@163.com"
```

## git环境变量设置
```shell
git config --global user.name codragonzuo
git config --global user.email codragon@163.com


git push  -u origin master
```

## go环境变量设置

```shell
[root@node1 filebeat]# go env
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
GOPATH="/root/beats/filebeat"
GOPRIVATE=""
GOPROXY="https://proxy.golang.org,direct"
GOROOT="/usr/local/go"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/usr/local/go/pkg/tool/linux_amd64"
GCCGO="gccgo"
AR="ar"
CC="gcc"
CXX="g++"
CGO_ENABLED="1"
GOMOD="/root/beats/go.mod"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build033069755=/tmp/go-build"
[root@node1 filebeat]# go env -w GOPROXY="https://goproxy.cn,direct"
```

```
[root@ecs-centos-7 thrift-0.13.0]# gcc  --version
gcc (GCC) 4.8.5 20150623 (Red Hat 4.8.5-39)
Copyright (C) 2015 Free Software Foundation, Inc.
This is free software; see the source for copying conditions.  There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

[root@ecs-centos-7 thrift-0.13.0]# sudo yum install gcc-c++
Loaded plugins: fastestmirror
Loading mirror speeds from cached hostfile
epel/x86_64/metalink                                                                             | 7.8 kB  00:00:00
 * base: mirrors.tuna.tsinghua.edu.cn
 * epel: mirrors.tuna.tsinghua.edu.cn
 * extras: mirrors.tuna.tsinghua.edu.cn
 * updates: mirrors.tuna.tsinghua.edu.cn
```

注意该过程需要你github的用户名和密码认证一下，该用户名也可以缓存在本地的配置文件中（[credential]部分），如果你用户名输入有误，则需要清除这个，不然后续认证会有问题。

上面提到的git remote add 和git config设置的参数是加上都是配置在git的配置文件中的。其中git remote添加的仓库Url链接信息是保存在项目目录.git/config文件中，这是项目级的配置文件。而git config设置的global信息保存在用户级的配置文件中，该文件的位置为用户目录~/.gitconfig。实际上git还有一个系统级的配置文件，该文件位于/etc/gitconfig。配置文件生效生效优先级是先项目（.git/config），后用户(~/.gitconfig),最后系统（/etc/gitconfig）。由于git的使用都是以用户级的，所以我们日常用的最多都是通过用户级别的配置，本文主要是用户级别配置文件的介绍。


## 代码字符批量替换
```shell
[root@node1 beats-source]# cd beats
[root@node1 beats]# sed -i "s/github.com\/elastic\/beats\/v7/github.com\/codragonzuo\/beats/g" `grep -rl "github.com/elastic/beats/v7" ./`
```

------------------------------------------------------------------------------------------------




```
python3 pysnmptrap.py 192.168.20.45:9000 -c public .1.3.6.1.4.1.2.3.1.2.1.2 a 192.168.0.250 0 0

snmpbulkwalk.py

snmptrap.py -v3 -e 8000000001020304 -l authPriv -u usr-sha-aes -A authkey1 -X privkey1 -a SHA -x AES demo.snmplabs.com 12345   1.3.6.1.4.1.20408.4.1.1.2 SNMPv2-MIB::sysName.0 = "my system"

snmptrap.py -v1 -c public demo.snmplabs.com 1.3.6.1.4.1.20408.4.1.1.2 0.0.0.0 1 0 0 1.3.6.1.2.1.1.1.0 s "my system"

snmptrap.py -v1 -c public 192.168.20.45:9000 1.3.6.1.4.1.20408.4.1.1.2 127.0.0.1 6 432 12345 1.3.6.1.2.1.1.1.0 s "my system"
```
## SNMP报文发送测试
```python
from pysnmp.hlapi import *

g= sendNotification(SnmpEngine(),CommunityData('public', mpModel=0),UdpTransportTarget(('127.0.0.1', 9000)), ContextData(), 'trap',        NotificationType( ObjectIdentity('1.3.6.1.4.1.20408.4.1.1.2.0.432'),).addVarBinds(('1.3.6.1.2.1.1.1.0', OctetString('my system')) ))

next(g)

##g= sendNotification(SnmpEngine(),CommunityData('public', mpModel=0),UdpTransportTarget(('192.168.20.45', 9000)), ContextData(), 'trap',        NotificationType( ObjectIdentity('1.3.6.1.4.1.20408.4.1.1.2.0.432'),).addVarBinds( ('1.3.6.1.2.1.1.3.0', 12345), ('1.3.6.1.6.3.18.1.3.0', '127.0.0.1'),    ('1.3.6.1.6.3.1.1.4.3.0', '1.3.6.1.4.1.20408.4.1.1.2'), ('1.3.6.1.2.1.1.1.0', OctetString('my system')) ))

```

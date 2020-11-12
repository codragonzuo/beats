
# Beats - The Lightweight Shippers of the Elastic Stack

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


msg = b'{"event":{"plugin_id":10,"plugin_sid":1000,"src_ip":"192.168.1.200","dst_ip":"192.168.100.9","src_port":100,"dst_port":101}}'
producer.send('snmp', msg)
```

## filebeat编译
```
cd beats/filetbeat
make
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
## 修改Linux系统时间
时间不对无法https上网，无法git下载代码
```
[root@node1]# date
[root@node1]# date -s 20200618
[root@node1]# date -s 19:32:30
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

## 修改kafka版本

```

[root@node1 beats]# vi ./libbeat/common/kafka/version.go

里列出支持的kafka版本

[root@node1 beats]# vi  libbeat/outputs/kafka/config.go

Version:          kafka.Version("1.0.0"),
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
[root@node1 .ssh]# ssh-add  ./id_rsa_dragon
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

## libpcap安装
```
yum install  libpcap-devel

##也可以手工编译libpcap 1.7.4

yum -y install flex
yum -y install bison



[root@node1 libpcap]# find / -name "libpcap.*"
/usr/local/lib/libpcap.so.1
/usr/local/lib/libpcap.so.1.7.4
/usr/local/lib/libpcap.so
/usr/local/lib/libpcap.a
/root/libpcap_master/libpcap.pc.in
/root/libpcap/libpcap-libpcap-1.7.4/Win32/Prj/libpcap.dsw
/root/libpcap/libpcap-libpcap-1.7.4/Win32/Prj/libpcap.dsp
/root/libpcap/libpcap-libpcap-1.7.4/libpcap.so.1.7.4
/root/libpcap/libpcap-libpcap-1.7.4/libpcap.a
[root@node1 libpcap]# cd ..
```


## 包冲突解决方法

### （1）
```
# github.com/codragonzuo/beats/packetbeat/protos/dhcpv4
protos/dhcpv4/dhcpv4.go:103:7: v4.Opcode undefined (type *dhcpv4.DHCPv4 has no field or method Opcode, but does have OpCode)
protos/dhcpv4/dhcpv4.go:124:39: v4.OpcodeToString undefined (type *dhcpv4.DHCPv4 has no field or method OpcodeToString)
protos/dhcpv4/dhcpv4.go:125:23: v4.HwTypeToString undefined (type *dhcpv4.DHCPv4 has no field or method HwTypeToString)
protos/dhcpv4/dhcpv4.go:126:32: cannot call non-function v4.HopCount (type uint8)
protos/dhcpv4/dhcpv4.go:127:59: cannot call non-function v4.TransactionID (type dhcpv4.TransactionID)
protos/dhcpv4/dhcpv4.go:128:34: cannot call non-function v4.NumSeconds (type uint16)
protos/dhcpv4/dhcpv4.go:130:23: v4.ClientHwAddrToString undefined (type *dhcpv4.DHCPv4 has no field or method ClientHwAddrToString)
protos/dhcpv4/dhcpv4.go:134:21: cannot call non-function v4.ClientIPAddr (type net.IP)
protos/dhcpv4/dhcpv4.go:135:44: cannot call non-function v4.ClientIPAddr (type net.IP)
protos/dhcpv4/dhcpv4.go:137:19: cannot call non-function v4.YourIPAddr (type net.IP)
protos/dhcpv4/dhcpv4.go:137:19: too many errors
```
解决办法：修改beats/go.mod里的require的dhcpv4版本号 
```
github.com/insomniacslk/dhcp v0.0.0-20180716144031-256240854619
```
### （2）thrift库冲突
解决办法：修改beats/go.mod里的require的thrift版本号 
```
github.com/samuel/go-thrift v0.0.0-20140522043831-2187045faa54
```
```
# github.com/codragonzuo/beats/packetbeat/protos/thrift
protos/thrift/thrift_idl.go:47:11: field.Id undefined (type *parser.Field has no field or method Id, but does have ID)
protos/thrift/thrift_idl.go:48:15: field.Id undefined (type *parser.Field has no field or method Id, but does have ID)
protos/thrift/thrift_idl.go:56:16: field.Id undefined (type *parser.Field has no field or method Id, but does have ID)

```
###  （3） C.TPACKET_V3未定义
```
pkg/mod/github.com/tsg/gopacket@v0.0.0-20190320122513-dd3d0e41124a/pcap/pcap_poll_linux.go:22:34: could not determine kind of name for C.TPACKET_V3
make: *** [packetbeat] 错误 2
```
官方解释： https://discuss.elastic.co/t/compile-packetbeat-on-centos6-8/190639/5

解决方法： 手工修改pcap_poll_linux.go里的 C.TPACKET_V3 为 2


## mage文件找不到

```
$ go get -u -v github.com/magefile/mage
github.com/magefile/mage (download)
github.com/magefile/mage/types
github.com/magefile/mage/mg
github.com/magefile/mage/parse
github.com/magefile/mage/sh
github.com/magefile/mage/mage
github.com/magefile/mage

$ cd ~/Go/src/github.com/magefile/mage
$ mage build
$ mage -version
Mage Build Tool v2
Build Date: <not set>
Commit: <not set>
```

metricbeat 编译
```
mage -v build
```

## go get
```
go get repo@branchname
go get repo@tag
go get repo@commithash


go get github.com/someone/some_module@master
go get github.com/someone/some_module@dev_branch

```



For example:

module /my/module

require (
...
github.com/someone/some_module latest
...
)

will become

module /my/module

require (
...
github.com/someone/some_module v2.0.39
...
)
after running go mod tidy

share  improve this answer  


```
Using “replace” in go.mod to point to your local module
Pam Selle   /   July 30, 2019   /  11 Comments
If you want to say, point to the local version of a dependency in Go rather than the one over the web, use the replace keyword.

The replace line goes above your require statements, like so:

module github.com/pselle/foo

replace github.com/pselle/bar => /Users/pselle/Projects/bar

require (
	github.com/pselle/bar v1.0.0
)
And now when you compile this module (go install), it will use your local code rather than the other dependency.
```

go list -m -versions github.com/docker/docker

## docker engine
```
replace (
	github.com/docker/docker v1.13.1 => github.com/docker/engine v0.0.0-20191113042239-ea84732a7725
)
```

## Thrift 安装

wget http://apache.mirrors.hoobly.com/thrift/0.13.0/thrift-0.13.0.tar.gz
tar -xvf thrift-0.13.0.tar.gz

cd thrift-0.13.0

./configure

./configure --without-nodejs
	

checking for bison version >= 2.5... no
configure: error: Bison version 2.5 or higher must be installed on the system!

[root@node1 thrift-0.13.0]# yum install  bison
已加载插件：fastestmirror, priorities, security
设置安装进程
Loading mirror speeds from cached hostfile
包 bison-2.4.1-5.el6.x86_64 已安装并且是最新版本

thrift -r --gen go tutorial.thrift



## gcc5.0安装
```

./contrib/download_prerequisites

https://stackoverflow.com/questions/9297933/cannot-configure-gcc-mpfr-not-found



安装依赖
./contrib/download_prerequisites

./configure

configure --disable-multilib   ？？？？

make


查看版本：
[root@node1 gcc-5.4.0]# /usr/local/bin/gcc --version
gcc (GCC) 5.4.0
Copyright © 2015 Free Software Foundation, Inc.
[root@node1 gcc-5.4.0]#

[root@node1 gcc-5.4.0]# /usr/local/bin/g++ --version
g++ (GCC) 5.4.0
Copyright © 2015 Free Software Foundation, Inc.




rm -rf /usr/bin/gcc
rm -rf /usr/bin/g++
ln -s /usr/local/gcc5/bin/gcc /usr/bin/gcc
ln -s /usr/local/gcc5/bin/g++ /usr/bin/g++



运行程序时可能会出现/lib64/libstdc++.so.6: version  `GLIBCXX_3.4.20' not found，是因为升级安装了gcc，生成的动态库没有替换老版本的gcc动态库导致的。

 查看包含最新的动态链接库的位置

find / -name "libstdc++.so*"


找到在/usr/local/gcc5/lib64/文件夹下

cp /usr/local/gcc5/lib64/libstdc++.so.6.0.21 /usr/lib64/libstdc++.so.6.0.21
rm -f /usr/lib64/libstdc++.so.6
ln /usr/lib64/libstdc++.so.6.0.21 /usr/lib64/libstdc++.so.6

[root@node1 gcc-5.4.0]# ls /usr/local/lib64/libstdc++* -all -t
-rw-r--r-- 1 root root     2397 7月   6 08:28 /usr/local/lib64/libstdc++.so.6.0.21-gdb.py
-rw-r--r-- 1 root root 28717220 7月   6 08:28 /usr/local/lib64/libstdc++.a
-rwxr-xr-x 1 root root      965 7月   6 08:28 /usr/local/lib64/libstdc++.la
lrwxrwxrwx 1 root root       19 7月   6 08:28 /usr/local/lib64/libstdc++.so -> libstdc++.so.6.0.21
lrwxrwxrwx 1 root root       19 7月   6 08:28 /usr/local/lib64/libstdc++.so.6 -> libstdc++.so.6.0.21
-rwxr-xr-x 1 root root 11169643 7月   6 08:28 /usr/local/lib64/libstdc++.so.6.0.21
-rw-r--r-- 1 root root 10838586 7月   6 08:28 /usr/local/lib64/libstdc++fs.a
-rwxr-xr-x 1 root root      905 7月   6 08:28 /usr/local/lib64/libstdc++fs.la
```

## librdkafka
```
[root@node1 correlation]# find  /  -name 'librdkafka*'
/usr/local/include/librdkafka
/usr/local/share/doc/librdkafka
//usr/local/lib/librdkafka++.so
/usr/local/lib/librdkafka.a
/usr/local/lib/librdkafka++.so.1
/usr/local/lib/librdkafka++.a
/usr/local/lib/librdkafka.so
/usr/local/lib/librdkafka.so.1
/root/librdkafka-1.5.0

[root@node1 correlation]# find  /  -name 'rdkafkacpp.h'
/usr/local/include/librdkafka/rdkafkacpp.h


[root@node1 librdkafka]# find  /  -name 'rdkafka*.h'
/usr/local/include/librdkafka/rdkafkacpp.h
/usr/local/include/librdkafka/rdkafka_mock.h
/usr/local/include/librdkafka/rdkafka.h


共享库找不到解决方案：
增加  export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH
source bash_profile

[root@node1 correlation]# cat ../.bash_profile
# .bash_profile

# Get the aliases and functions
if [ -f ~/.bashrc ]; then
        . ~/.bashrc
fi

# User specific environment and startup programs

PATH=$PATH:$HOME/bin

export PATH



export LD_LIBRARY_PATH=/usr/local/lib:$LD_LIBRARY_PATH

```

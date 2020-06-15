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

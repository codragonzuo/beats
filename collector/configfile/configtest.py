#coding:utf-8
import yaml
#import os
import sys

fr = open('autotest.yml', 'r', encoding='utf-8')
c = fr.read()
y = yaml.load(c, Loader=yaml.FullLoader)



data = y['USECASE']

inputlist = data['input']
monitor = data['monitoring']




clist = list()
cinput = dict()
plist = list()


for input in inputlist:
    print(input)
    print(input['type'])
    if    input['type'] == 'file':
        path = list()
        cinput['type'] = 'log'
        cinput['enabled'] = 'true'
        cinput['name'] = input['name']
        cinput['label'] = input['label']
        path.append(input['path'])
        cinput['paths'] = path
        clist.append(cinput)
    elif  input['type'] == 'snmptrap':
        cinput = dict()
        cinput['type'] = 'snmptrap'
        cinput['name'] = input['name']
        cinput['label'] = input['label']
        cinput['enabled'] = 'true'
        protocol = dict()
        protocol['host'] = input['host']
        cinput['protocol.udp'] = protocol
        clist.append(cinput)
    elif  input['type'] == 'netflow':
        cinput = dict()
        cinput['type'] = 'netflow'
        cinput['name'] = input['name']
        cinput['label'] = input['label']
        cinput['enabled'] = 'true'
        cinput['host'] = input['host']
        cinput['protocols'] = [ 'v5', 'v9', 'ipfix' ]
        cinput['max_message_size'] = '10KiB'
        cinput['queue_size'] =  8192
        cinput['expiration_timeout'] = '30m'
        clist.append(cinput)
    elif  input['type'] == 'dblog':
        cinput = dict()
        cinput['type'] = 'dblog'
        cinput['name'] = input['name']
        cinput['label'] = input['label']
        cinput['host'] = input['host']
        cinput['username'] = input['username']
        cinput['password'] = input['password']
        cinput['dbname'] = input['dbname']
        cinput['dbtype'] = input['dbtype']
        cinput['id_name'] = input['id_name']
        cinput['id_start'] = input['id_start']
        cinput['query_sql'] = input['query_sql']
        cinput['max_message_num'] = input['max_message_num']
        clist.append(cinput)
    elif  input['type'] == 'monitor':
        cinput = dict()
        cinput['type'] = 'monitor'
        cinput['name'] = input['name']
        cinput['label'] = input['label']
        cinput['enabled'] = 'true'
        clist.append(cinput)


if 'monitoring' in data.keys():
    print(data['monitoring'])
    monitoring=dict()
    monitoring['enabled'] = data['monitoring']['enable']
    monitoring['cluster_uuid'] = 'metron'
    es = dict()
    es['hosts'] = ["https://192.168.20.45:9200"]
    es['username'] = 'metron'
    es['password'] = 'metron'
    monitoring['elasticsearch'] = es


if 'outputs' in data.keys():
    print(data['outputs'])
    outputs = dict()
    hostlist = list()
    host = data['outputs']['host']
    hostlist.append(host)
    outputs['hosts'] = hostlist
    outputs['topic'] = data['outputs']['topic']
    outputs['compression'] = 'gzip'
    outputs['max_message_bytes'] = 1000000
    outputs['required_acks'] = 1


cdata = dict()
cdata['filebeat.inputs'] = clist
cdata['output.kafka']=outputs
cdata['monitoring']=monitoring
print(cdata)

fw = open('filebeat.yml', 'w', encoding='utf-8')
yaml.dump(cdata,fw,allow_unicode=True)

##########################################################
if 'packet' in data.keys():
    packet = dict()
    packet['interface'] = data['packet']['interface']
    packet['filter'] = data['packet']['filter']
    pdata = dict()
    pdata['output.kafka']=outputs
    pdata['monitoring']=monitoring
    pdata['packetbeat.interfaces.device'] = packet['interface']
    pdata['packetbeat.interfaces.bpf_filter'] = packet['filter']
    flow=dict()
    flow = dict()
    flow['timeout'] = '30s'
    flow['period'] = '10s'
    pdata['packetbeat.flows'] = flow
    fw = open('packetbeat.yml', 'w', encoding='utf-8')
    yaml.dump(pdata,fw,allow_unicode=True)


"""
# 方法1
print '遍历列表方法1：'
for i in list:
    print ("序号：%s   值：%s" % (list.index(i) + 1, i))

print '\n遍历列表方法2：'
# 方法2
for i in range(len(list)):
    print ("序号：%s   值：%s" % (i + 1, list[i]))

# 方法3
print '\n遍历列表方法3：'
for i, val in enumerate(list):
    print ("序号：%s   值：%s" % (i + 1, val))

# 方法3
print '\n遍历列表方法3 （设置遍历开始初始位置，只改变了起始序号）：'
for i, val in enumerate(list, 2):
    print ("序号：%s   值：%s" % (i + 1, val))


"""
#coding:utf-8
import yaml
#import os
import sys

#reload(sys)
#sys.setdefaultencoding('utf-8')

# 追加写入，w,覆盖写入
yamlPath = 'out.yml'

if sys.version < '3':
    fw = open(yamlPath,'w')
else:
    fw = open(yamlPath,'w',encoding='utf-8')

# 构建数据
data = {"cookie1":{'domain': '.yiyao.cc', 'expiry': 1521558688.480118, 'httpOnly': False, 'name': '_ui_', 'path': '/', 'secure': False, 'value': 'HSX9fJjjCIImOJoPUkv/QA=='}}
# 装载数据
yaml.dump(data,fw)
# 读取数据，获取文件
if sys.version < '3':
    f = open(yamlPath,'r')
else:
    f = open(yamlPath,'r',encoding='utf-8')

# 读取文件
cont = f.read()
# 加载数据
x = yaml.load(cont, Loader=yaml.FullLoader)
# 打印数据
print(x)
# 打印读取写入的数据
print(x.get("cookie1"))

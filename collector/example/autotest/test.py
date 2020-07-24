import yaml
f = open('autotest.yml', 'r', encoding='utf-8')
x = f.read()
y = yaml.load(x, Loader=yaml.FullLoader)
fw = open('autotest-tmp.yml', 'w', encoding='utf-8')
yaml.dump(y,fw,allow_unicode=True)


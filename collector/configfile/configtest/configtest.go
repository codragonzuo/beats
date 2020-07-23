package main

import (
	"fmt"
	"log"
	//"github.com/ghodss/yaml"
	yaml2 "gopkg.in/yaml.v2"
	"io/ioutil"
)



func main() {
	fmt.Printf("config test \n")
	t := make(map[string] interface{})
	buffer, err := ioutil.ReadFile("configtest.yml")
	err = yaml2.Unmarshal(buffer, &t)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("%v\n\n\n\n",t)

	for country := range t {
		fmt.Println(country, t [country])
	}



        testBBB()
}


func Yamltojsontest() {
    t := make(map[string] interface{})
    buffer, err := ioutil.ReadFile("/root/beats/filebeat/filebeat.yml")
    if err != nil {
        fmt.Printf("fileread error !\n")
    }
    err = yaml2.Unmarshal(buffer, &t)
    if err != nil {
        log.Fatalf(err.Error())
    }
    fmt.Printf("%v\n\n\n\n",t)

    for country := range t {
        fmt.Println(country, t [country])
    }

}



func testttt() {
    mapNum := make(map[string]int)
    mapNum["key1"] = 1
    mapNum["key2"] = 2
    mapNum["key3"] = 3
    mapNum["key4"] = 4
    //输出map集合key和value
    for k, v := range mapNum {
        fmt.Printf("key:%s value:%d \n", k, v)
    }
    //只输出value
    for _, v := range mapNum {
        fmt.Printf(" value:%d \n", v)
    }
    //只输出key
    for k := range mapNum {
        fmt.Printf("key:%s \n", k)
    }
}




var data = `
blog: xiaorui.cc
best_authors: ["fengyun","lee","park"]
desc:
  counter: 521
  plist: [3, 4]
`


type T struct {
    Blog    string
    Authors []string `yaml:"best_authors,flow"`
    Desc    struct {
        Counter int   `yaml:"Counter"`
        Plist   []int `yaml:",flow"`
    }
}





func testBBB() {
    t := T{}
    //把yaml形式的字符串解析成struct类型
    err := yaml2.Unmarshal([]byte(data), &t)
    //修改struct里面的记录
    t.Blog = "this is Blog"
    t.Authors = append(t.Authors, "myself")
    t.Desc.Counter = 99
    fmt.Printf("--- t:\n%v\n\n", t)
    //转换成yaml字符串类型
    d, err := yaml2.Marshal(&t)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
    fmt.Printf("--- t dump:\n%s\n\n", string(d))
}



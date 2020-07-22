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


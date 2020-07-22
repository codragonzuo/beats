package configfile

import (
    "fmt"
    "log"
    "github.com/ghodss/yaml"
    yaml2 "gopkg.in/yaml.v2"
    "io/ioutil"
)



//  var { filebeatCfg map[string]interface{}{}  }



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





func Jsontoyamltest() {

    j := []byte(`{"name": "John", "age": 30}`)
    y, err := yaml.JSONToYAML(j)
    if err != nil {
        fmt.Printf("err: %v\n", err)
        return
    }
    fmt.Println(string(y))
    fmt.Printf("\n")
    /* Output:
    name: John
    age: 30
    */
    j2, err := yaml.YAMLToJSON(y)
    if err != nil {
        fmt.Printf("err: %v\n", err)
        return
    }
    fmt.Println(string(j2))
    /* Output:
    {"age":30,"name":"John"}
    */

    fmt.Printf("\n\n\n")
}





package main

import (
    "fmt"
    "log"
    "runtime"
    "os"
    "os/signal"
    "os/exec"
    "time"
    "github.com/codragonzuo/beats/collector/meal"
    "github.com/codragonzuo/beats/collector/life"
    "github.com/codragonzuo/beats/collector/rpcserver"
    _ "github.com/codragonzuo/beats/collector/configfile"
    _ "github.com/codragonzuo/beats/collector/commu/commu"
    _ "github.com/codragonzuo/beats/collector/commu/sharedlib"
    "github.com/apache/thrift/lib/go/thrift"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"

    "github.com/bitly/go-simplejson"
    "encoding/json"
)



/*
 * Tag... - a very simple struct
 */
type Tag struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}



func main() {
    meal.Getlunch()
    
    life.PrintWork()

    fmt.Println("Hello World !")

    //go say ("world")
    //say ("hello")
    //sqlquery()

    //cmdtest()

    //configfile.Jsontoyamltest()
    //configfile.Yamltojsontest()
    
    runserver()

    jsontest()
}

func sayhello(){
    say ("hello")
    sqlquery()

}

func say(s string){
    //for  {
    for i:=0; i<5; i++ {
        runtime.Gosched()
        fmt.Println(s)
    }
}




func cmdtest() { 
    cmd := exec.Cmd{ 
        Path: "/root/beats/filebeat/filebeat", 
//      Args: []string{"-u", "-l", "8888"}, 
        Dir:  "/root/beats/filebeat", 
    } 
 
    if err := cmd.Start(); err != nil { 
        log.Panic(err) 
    } 
 
    fmt.Println("Start child process with pid", cmd.Process.Pid) 
 
    // Wait releases any resources associated with the Cmd 
    go func() { 
        if err := cmd.Wait(); err != nil { 
            fmt.Printf("Child process %d exit with err: %v\n", cmd.Process.Pid, err) 
        } 
    }() 
 
    // After five second, kill cmd's process 
    time.Sleep(5 * time.Second) 
    cmd.Process.Kill() 
    fmt.Printf("Wait for kill success\n")
    time.Sleep(30 * time.Second) 

    c := make(chan os.Signal, 1) 
    signal.Notify(c, os.Interrupt) 
    s := <-c 
    fmt.Println("Got signal:", s) 
}


func runserver(){
       var protocolFactory thrift.TProtocolFactory
       var transportFactory thrift.TTransportFactory

	//compact
	protocolFactory = thrift.NewTCompactProtocolFactory()
	//simplejson
	//protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	//case "json":
	protocolFactory = thrift.NewTJSONProtocolFactory()
	//case "binary", "":
	//protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	//buffered
	transportFactory = thrift.NewTBufferedTransportFactory(8192)
	//} else {
	//transportFactory = thrift.NewTTransportFactory()
	
	//framed
        //transportFactory = thrift.NewTFramedTransportFactory(transportFactory)

        if err := rpcserver.RunServer(transportFactory, protocolFactory, "192.168.20.45:8091", false); err != nil {
			fmt.Println("error running server:", err)
        }

}




func sqlquery(){
    
    db, err := sql.Open("mysql", "root:qwer1234@tcp(127.0.0.1:3306)/ambari")

    // if there is an error opening the connection, handle it
    if err != nil {
        log.Print(err.Error())
        fmt.Printf("connect failed ! \n")
        return
    }
    defer db.Close()

    // Execute the query
    results, err := db.Query("select alert_id ,alert_label from alert_history where alert_id > 8840")
    if err != nil {
        fmt.Printf("db.query error \n")
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    for results.Next() {
        var tag Tag
        // for each row, scan the result into our tag composite object
        err = results.Scan(&tag.ID, &tag.Name)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }
                // and then print out the tag's Name attribute
        log.Printf(tag.Name)
        fmt.Printf("%d  %s\n", tag.ID, tag.Name)
    }

}


func jsontest() {
       //拼凑json   body为map数组
	var rbody []map[string]interface{}
	t := make(map[string]interface{})
	t["device_id"] = "dddddd"
	t["device_hid"] = "ddddddd"
 
	rbody = append(rbody, t)
	t1 := make(map[string]interface{})
	t1["device_id"] = "aaaaa"
	t1["device_hid"] = "aaaaa"
 
	rbody = append(rbody, t1)
 
	cnnJson := make(map[string]interface{})
	cnnJson["code"] = 0
	cnnJson["request_id"] = 123
	cnnJson["code_msg"] = ""
	cnnJson["body"] = rbody
	cnnJson["page"] = 0
	cnnJson["page_size"] = 0
 
	b, _ := json.Marshal(cnnJson)
	cnnn := string(b)
	fmt.Println("cnnn:%s", cnnn)
	cn_json, _ := simplejson.NewJson([]byte(cnnn))
	cn_body, _ := cn_json.Get("body").Array()
 
	for _, di := range cn_body {
                //就在这里对di进行类型判断
		newdi, _ := di.(map[string]interface{})
		device_id := newdi["device_id"]
		device_hid := newdi["device_hid"]
		fmt.Println(device_hid, device_id)
	}
 
}


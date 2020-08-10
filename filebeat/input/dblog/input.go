// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package dblog

import (
	_ "strings"
        _ "strconv"
	"sync"
	"time"
        "fmt"
	_ "github.com/pkg/errors"
    //"github.com/bitly/go-simplejson"
    _ "encoding/json"
	"github.com/codragonzuo/beats/filebeat/channel"
	"github.com/codragonzuo/beats/filebeat/harvester"
	"github.com/codragonzuo/beats/filebeat/input"
	_ "github.com/codragonzuo/beats/filebeat/inputsource"
	"github.com/codragonzuo/beats/libbeat/beat"
	"github.com/codragonzuo/beats/libbeat/common"
	"github.com/codragonzuo/beats/libbeat/common/cfgwarn"
	"github.com/codragonzuo/beats/libbeat/logp"
        "database/sql"
        _ "github.com/go-sql-driver/mysql"

)

// Parser is generated from a ragel state machine using the following command:
//go:generate ragel -Z -G2 parser.rl -o parser.go
//go:generate goimports -l -w parser.go

// Severity and Facility are derived from the priority, theses are the human readable terms
// defined in https://tools.ietf.org/html/rfc3164#section-4.1.1.
//
// Example:
// 2 => "Critical"
type mapper []string


var (
    Monitorfowwarder * harvester.Forwarder
)


/*
 * Tag... - a very simple struct
 */
type Tag struct {
    ID   int32    `json:"id"`
    Name string `json:"name"`
}





func init() {
    fmt.Printf("dblog  input Register\n")
	err := input.Register("dblog", NewInput)
	if err != nil {
		panic(err)
	}
}

// Input define a snmptrap input
type Input struct {
	sync.Mutex
	started bool
	outlet  channel.Outleter
	config  *config
	log     *logp.Logger
	done    chan interface{}
	myfowwarder * harvester.Forwarder
	id_start int32
        IdName string
}

// NewInput creates a new syslog input
func NewInput(
	cfg *common.Config,
	outlet channel.Connector,
	context input.Context,
) (input.Input, error) {
	cfgwarn.Experimental("dblog input type is used")

	log := logp.NewLogger("dblog")

	out, err := outlet.ConnectWith(cfg, beat.ClientConfig{
		Processing: beat.ProcessingConfig{
			DynamicFields: context.DynamicFields,
		},
	})
	if err != nil {
		return nil, err
	}
	config := defaultConfig
    fmt.Printf("dblog NewInput config=%v\n", config)
	if err = cfg.Unpack(&config); err != nil {
		return nil, err
	}
	fmt.Printf("dblog NewInput config=%v\n", config)

	forwarder := harvester.NewForwarder(out)

	fmt.Printf("input dblog NewInput forwarder=%v\n", forwarder)

	if err != nil {
		return nil, err
	}

	return &Input{
		outlet:  out,
		started: false,
		config:  &config,
		log:     log,
		myfowwarder: forwarder,
		id_start: config.IdStart,
                IdName: config.IdName,
	}, nil
}

// Run starts listening for Syslog events over the network.
func (p *Input) Run() {
	p.Lock()
	defer p.Unlock()

	if !p.started {
		p.started = true
		p.done =  make(chan interface{})
		go func (){
			for {
				select {
					case <-p.done: 
						return
					default:
				}
				time.Sleep(5*time.Second)
				p.sqlquery()
			}
		}()
	}
}

// Stop stops the syslog input.
func (p *Input) Stop() {
	defer p.outlet.Close()
	p.Lock()
	defer p.Unlock()

	if !p.started {
		return
	}

	p.log.Info("Stopping Syslog input")
	p.done <- 1
	p.started = false
}

// Wait stops the syslog input.
func (p *Input) Wait() {
	p.Stop()
}




func (p *Input) sqlquery(){
    
//    db, err := sql.Open("mysql", "root:qwer1234@tcp(127.0.0.1:3306)/ambari")
     db_conn := fmt.Sprintf("%s:%s@tcp(%s)/%s", p.config.Username, p.config.Password, p.config.Host, p.config.DBName)
	 fmt.Printf("DBType=%s, db_conn=%s\n", p.config.DBType, db_conn)
	 db, err := sql.Open(p.config.DBType, db_conn)
    // if there is an error opening the connection, handle it
    if err != nil {
        //log.Print(err.Error())
        fmt.Printf("connect failed ! \n")
        return
    }
    defer db.Close()

    // Execute the query
    // results, err := db.Query("select alert_id ,alert_label from alert_history where alert_id > 8960") 
	
	sql_query := fmt.Sprintf("%s where %s > %d", p.config.QueryString, p.config.IdName, p.id_start) //p.config.IdStart)
	fmt.Println(sql_query)
	results, err := db.Query(sql_query)
    if err != nil {
        fmt.Printf("db.query error \n")
        return
    }
	
	p.DoRowsMapper(results, p.myfowwarder)
    
    for results.Next() {
        var tag Tag
        // for each row, scan the result into our tag composite object
        err = results.Scan(&tag.ID, &tag.Name)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }

        fmt.Printf("%d  %d\n", tag.ID, p.id_start)
		if p.id_start < tag.ID {
			p.id_start = tag.ID
		}
    }

}

func (p *Input) DoRowsMapper(rows *sql.Rows, forwarder * harvester.Forwarder) () { 
 
  // 获取列名 
  columns, err := rows.Columns() 
  if err != nil { 
    panic(err.Error()) // proper error handling instead of panic in your app 
  } 
 
  // Make a slice for the values 
  values := make([]sql.RawBytes, len(columns)) 
 
  scanArgs := make([]interface{}, len(values)) 
  for i := range values { 
    scanArgs[i] = &values[i] 
  } 
  var rbody []map[string] string//interface{}
	for rows.Next() {
		err = rows.Scan(scanArgs...) 
		if err != nil { 
			panic(err.Error())
		} 

                //var id_start int
                //id_start = 0 
		rowMap := make(map[string]string) 
		var value string 
		for i, col := range values { 
			if col != nil { 
				value = string(col) 
				rowMap[columns[i]] = value 
			        if columns[i] == p.IdName {
                                     fmt.Printf("idname= %s, value=%s\n", columns[i], value)
                                     //id_start, _ := strconv.Atoi(value)
                                 }
                        }
		}
                
		
                //if p.id_start < int32(id_start) {
                //    p.id_start = int32(id_start)
                //}
                //fmt.Printf(" %v\n", rowMap)
		rbody = append(rbody, rowMap)
	}
        //fmt.Printf("id_start=%d\n", p.id_start)
	event := beat.Event{Fields: common.MapStr{}}
	event.PutValue("@dblog", "dblog-dragon")
	event.PutValue("dblog", rbody)
	forwarder.Send(event)
}


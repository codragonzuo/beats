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
	"strings"
	"sync"
	"time"
        "fmt"
	"github.com/pkg/errors"

	"github.com/codragonzuo/beats/filebeat/channel"
	"github.com/codragonzuo/beats/filebeat/harvester"
	"github.com/codragonzuo/beats/filebeat/input"
	"github.com/codragonzuo/beats/filebeat/inputsource"
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
    ID   int    `json:"id"`
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
//	server  inputsource.Network
	config  *config
	log     *logp.Logger
        done    chan interface{}
        Monitorfowwarder * harvester.Forwarder
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
	if err = cfg.Unpack(&config); err != nil {
		return nil, err
	}

	forwarder := harvester.NewForwarder(out)
	//callback := func(data []byte, metadata inputsource.NetworkMetadata) {
	//	ev := parseAndCreateEvent(data, metadata, time.Local, log)
	//	forwarder.Send(ev)
	//}
        fmt.Printf("input dblog NewInput forwarder=%v\n", forwarder)
        //server, err := factory(callback, config.Protocol)
	if err != nil {
		return nil, err
	}

	return &Input{
		outlet:  out,
		started: false,
//		server:  server,
		config:  &config,
		log:     log,
                Monitorfowwarder: forwarder,
	}, nil
}

// Run starts listening for Syslog events over the network.
func (p *Input) Run() {
	p.Lock()
	defer p.Unlock()

	if !p.started {
		//p.log.Infow("Starting Syslog input", "protocol", p.config.Protocol.Name())
		//err := p.server.Start()
		//if err != nil {
		//	p.log.Error("Error starting the server", "error", err)
		//	return
		//}
		p.started = true
                p.done =  make(chan interface{})
                go func (){
                 //defer  
                 for {
                    select {
                        case <-p.done: 
                            return
                       default:
                    }
                    time.Sleep(5*time.Second)
                    sqlquery()
                    //Monitorfowwarder
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
	//p.server.Stop()
        //close(p.done)
        //p.done.Close()
        p.done <- 1
	p.started = false

}

// Wait stops the syslog input.
func (p *Input) Wait() {
	p.Stop()
}

func createEvent(ev *event, metadata inputsource.NetworkMetadata, timezone *time.Location, log *logp.Logger) beat.Event {
	f := common.MapStr{
		"message": strings.TrimRight(ev.Message(), "\n"),
	}
	return newBeatEvent(ev.Timestamp(timezone), metadata, f)
}

func parseAndCreateEvent(data []byte, metadata inputsource.NetworkMetadata, timezone *time.Location, log *logp.Logger) beat.Event {
	ev := newEvent()
	//Parse(data, ev)
	//if !ev.IsValid() {
	//	log.Errorw("can't parse event as syslog rfc3164", "message", string(data))
	//	return newBeatEvent(time.Now(), metadata, common.MapStr{
	//		"message": string(data),
	//	})
	//}
	return createEvent(ev, metadata, time.Local, log)
}

func newBeatEvent(timestamp time.Time, metadata inputsource.NetworkMetadata, fields common.MapStr) beat.Event {
	event := beat.Event{
		Timestamp: timestamp,
		Meta: common.MapStr{
			"truncated": metadata.Truncated,
		},
		Fields: fields,
	}
	if metadata.RemoteAddr != nil {
		event.Fields.Put("log.source.address", metadata.RemoteAddr.String())
	}
	return event
}

func mapValueToName(v int, m mapper) (string, error) {
	if v < 0 || v >= len(m) {
		return "", errors.Errorf("value out of bound: %d", v)
	}
	return m[v], nil
}


func sqlquery(){
    
    db, err := sql.Open("mysql", "root:qwer1234@tcp(127.0.0.1:3306)/ambari")

    // if there is an error opening the connection, handle it
    if err != nil {
        //log.Print(err.Error())
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
        //log.Printf(tag.Name)
        fmt.Printf("%d  %s\n", tag.ID, tag.Name)
    }

}


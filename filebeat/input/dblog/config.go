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
	"fmt"
	"time"

	"github.com/codragonzuo/beats/libbeat/common"
	"github.com/codragonzuo/beats/libbeat/common/transport/kerberos"
	"github.com/codragonzuo/beats/libbeat/common/transport/tlscommon"
)

type config struct {
	Topics                    common.ConfigNamespace `config:"topics"`
	ClientID                 string            `config:"client_id"`
    Username                 string            `config:"username" validate:"required"`
    Password                 string            `config:"password" validate:"required"`
    Host                     string            `config:"host" validate:"required"`
    DBType                   string            `config:"dbtype" validate:"required"`
    DBName                   string            `config:"dbname" validate:"required"`
	IdName                   string            `config:"id_name"  validate:"required"`
	IdStart                  int32             `config:"id_start"  validate:"required"`
	MaxMessageNum            int32             `config:"max_message_num"  validate:"required"`
	QueryString              string            `config:"query_sql"  validate:"required"`
}

var defaultConfig = config{
	ClientID : "client_xxx",
}





func getconfig(
	config common.ConfigNamespace,
) (error) {
	n, cfg := config.Name(), config.Config()
        
    fmt.Printf("filebeat input dblog config=%s\n", n)
	dbconfig := defaultConfig
	if err := cfg.Unpack(&dbconfig); err != nil {
		return err
	}
    return nil
}



type kafkaInputConfig struct {
        Hosts                    []string          `config:"hosts" validate:"required"`
        Topics                   []string          `config:"topics" validate:"required"`
        ClientID                 string            `config:"client_id"`
        ConnectBackoff           time.Duration     `config:"connect_backoff" validate:"min=0"`
        TLS                      *tlscommon.Config `config:"ssl"`
        Kerberos                 *kerberos.Config  `config:"kerberos"`
        Username                 string            `config:"username"`
        Password                 string            `config:"password"`
		DBType                 string            `config:"password"`
}

type ConfigFetch struct {
        Min     int32 `config:"min" validate:"min=1"`
        Default int32 `config:"default" validate:"min=1"`
        Max     int32 `config:"max" validate:"min=0"`
}



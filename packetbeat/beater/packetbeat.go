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

package beater

import (
	"errors"
        "flag"
	"fmt"
	"sync"
	"time"
        "encoding/base64"

	"github.com/tsg/gopacket/layers"

	"github.com/codragonzuo/beats/libbeat/beat"
	"github.com/codragonzuo/beats/libbeat/common"
	"github.com/codragonzuo/beats/libbeat/logp"
	"github.com/codragonzuo/beats/libbeat/processors"
	"github.com/codragonzuo/beats/libbeat/service"
        //mypipeline  "github.com/codragonzuo/beats/libbeat/publisher/pipeline"
        //PB "github.com/codragonzuo/beats/packetbeat/pb"


	"github.com/codragonzuo/beats/packetbeat/config"
	"github.com/codragonzuo/beats/packetbeat/decoder"
	"github.com/codragonzuo/beats/packetbeat/flows"
	"github.com/codragonzuo/beats/packetbeat/procs"
	"github.com/codragonzuo/beats/packetbeat/protos"
	"github.com/codragonzuo/beats/packetbeat/protos/icmp"
	"github.com/codragonzuo/beats/packetbeat/protos/tcp"
	"github.com/codragonzuo/beats/packetbeat/protos/udp"
	"github.com/codragonzuo/beats/packetbeat/publish"
	"github.com/codragonzuo/beats/packetbeat/sniffer"

	"github.com/tsg/gopacket"
	//"github.com/tsg/gopacket/layers"

	// Add packetbeat default processors
	_ "github.com/codragonzuo/beats/packetbeat/processor/add_kubernetes_metadata"
)

// Beater object. Contains all objects needed to run the beat
type packetbeat struct {
	config      config.Config
	cmdLineArgs flags
	sniff       *sniffer.Sniffer

	// publisher/pipeline
	pipeline beat.Pipeline
	transPub *publish.TransactionPublisher
	flows    *flows.Flows
        //MsgPublisher * mypipeline.Pipeline
        
}

type flags struct {
	file       *string
	loop       *int
	oneAtAtime *bool
	topSpeed   *bool
	dumpfile   *string
}

var cmdLineArgs flags

func init() {
	cmdLineArgs = flags{
		file:       flag.String("P", "", "Read packet data from specified file"),
		loop:       flag.Int("pl", 1,  "Loop file. 0 - loop forever"),
		oneAtAtime: flag.Bool("pO", false, "Read packets one at a time (press Enter)"),
		topSpeed:   flag.Bool("pt", false, "Read packets as fast as possible, without sleeping"),
		dumpfile:   flag.String("pdump", "", "Write all captured packets to this libpcap file"),
	}
}

func New(b *beat.Beat, rawConfig *common.Config) (beat.Beater, error) {
	config := config.Config{
		Interfaces: config.InterfacesConfig{
			File:       *cmdLineArgs.file,
			Loop:       *cmdLineArgs.loop,
			TopSpeed:   *cmdLineArgs.topSpeed,
			OneAtATime: *cmdLineArgs.oneAtAtime,
			Dumpfile:   *cmdLineArgs.dumpfile,
		},
	}
	err := rawConfig.Unpack(&config)
	if err != nil {
		logp.Err("fails to read the beat config: %v, %v", err, config)
		return nil, err
	}



        fmt.Printf("Newbeat\n")
	pb := &packetbeat{
		config:      config,
		cmdLineArgs: cmdLineArgs,
	}
	err = pb.init(b)
        //pb.MsgPublisher = pb.pipeline
	if err != nil {
		return nil, err
	}

	return pb, nil
}

// init packetbeat components
func (pb *packetbeat) init(b *beat.Beat) error {
	var err error
	cfg := &pb.config
	// Enable the process watcher only if capturing live traffic
	if cfg.Interfaces.File == "" {
		err = procs.ProcWatcher.Init(cfg.Procs)
		if err != nil {
			logp.Critical(err.Error())
			return err
		}
	} else {
		logp.Info("Process watcher disabled when file input is used")
	}

	pb.pipeline = b.Publisher
        fmt.Printf("packetbeat pb.pipeline=%v\n", pb.pipeline)
	pb.transPub, err = publish.NewTransactionPublisher(
		b.Info.Name,
		b.Publisher,
		pb.config.IgnoreOutgoing,
		pb.config.Interfaces.File == "",
	)
	if err != nil {
		return err
	}

	logp.Debug("main", "Initializing protocol plugins")
	err = protos.Protos.Init(false, pb.transPub, cfg.Protocols, cfg.ProtocolsList)
	if err != nil {
		return fmt.Errorf("Initializing protocol analyzers failed: %v", err)
	}

	if err := pb.setupFlows(); err != nil {
		return err
	}

	return pb.setupSniffer()
}

func (pb *packetbeat) setupSniffer() error {
	config := &pb.config

	icmp, err := pb.icmpConfig()
	if err != nil {
		return err
	}

	withVlans := config.Interfaces.WithVlans
	withICMP := icmp.Enabled()

	filter := config.Interfaces.BpfFilter
	if filter == "" && !config.Flows.IsEnabled() {
		filter = protos.Protos.BpfFilter(withVlans, withICMP)
	}

	pb.sniff, err = sniffer.New(false, filter, pb.createWorker,  pb.createMyWorker ,config.Interfaces)
	return err
}

func (pb *packetbeat) setupFlows() error {
	config := &pb.config
	if !config.Flows.IsEnabled() {
		return nil
	}

	processors, err := processors.New(config.Flows.Processors)
	if err != nil {
		return err
	}

	client, err := pb.pipeline.ConnectWith(beat.ClientConfig{
		Processing: beat.ProcessingConfig{
			EventMetadata: config.Flows.EventMetadata,
			Processor:     processors,
			KeepNull:      config.Flows.KeepNull,
		},
	})
	if err != nil {
		return err
	}

	pb.flows, err = flows.NewFlows(client.PublishAll, config.Flows)
	if err != nil {
		return err
	}

	return nil
}

func (pb *packetbeat) Run(b *beat.Beat) error {
	defer func() {
		if service.ProfileEnabled() {
			logp.Debug("main", "Waiting for streams and transactions to expire...")
			time.Sleep(time.Duration(float64(protos.DefaultTransactionExpiration) * 1.2))
			logp.Debug("main", "Streams and transactions should all be expired now.")
		}
	}()

	defer pb.transPub.Stop()

	timeout := pb.config.ShutdownTimeout
	if timeout > 0 {
		defer time.Sleep(timeout)
	}

	if pb.flows != nil {
		pb.flows.Start()
		defer pb.flows.Stop()
	}

	var wg sync.WaitGroup
	errC := make(chan error, 1)

	// Run the sniffer in background
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := pb.sniff.Run()
		if err != nil {
			errC <- fmt.Errorf("Sniffer main loop failed: %v", err)
		}
	}()

	logp.Debug("main", "Waiting for the sniffer to finish")
	wg.Wait()
	select {
	default:
	case err := <-errC:
		return err
	}

	return nil
}

// Called by the Beat stop function
func (pb *packetbeat) Stop() {
	logp.Info("Packetbeat send stop signal")
	pb.sniff.Stop()
}

func (pb *packetbeat) createWorker(dl layers.LinkType) (sniffer.Worker, error) {
	var icmp4 icmp.ICMPv4Processor
	var icmp6 icmp.ICMPv6Processor
	cfg, err := pb.icmpConfig()
	if err != nil {
		return nil, err
	}
	if cfg.Enabled() {
		reporter, err := pb.transPub.CreateReporter(cfg)
		if err != nil {
			return nil, err
		}

		icmp, err := icmp.New(false, reporter, cfg)
		if err != nil {
			return nil, err
		}

		icmp4 = icmp
		icmp6 = icmp
	}

	tcp, err := tcp.NewTCP(&protos.Protos)
	if err != nil {
		return nil, err
	}

	udp, err := udp.NewUDP(&protos.Protos)
	if err != nil {
		return nil, err
	}


        config := &pb.config
        if !config.Flows.IsEnabled() {
                return nil,nil
        }

        processors, err := processors.New(config.Flows.Processors)
        if err != nil {
                return nil, err
        }

        client, err := pb.pipeline.ConnectWith(beat.ClientConfig{
                Processing: beat.ProcessingConfig{
                        EventMetadata: config.Flows.EventMetadata,
                        Processor:     processors,
                        KeepNull:      config.Flows.KeepNull,
                },
        })


	worker, err := decoder.New(pb.flows, dl, icmp4, icmp6, tcp, udp, client)
	if err != nil {
		return nil, err
	}

	return worker, nil
}



type MyCoder struct {
        ss string	
        client beat.Client
}


func (pb *packetbeat) Myworker(ch chan beat.Event, client beat.Client) {
	//for {
	//	select {
	//	//case <-pb.done:
	//	//	return
	//	case event := <-ch:
	//		//pub, _ := p.processor.Run(&event)
	//		fields, err := pb.GetFields(event.Fields)
        //              if fields != nil {
	//			client.Publish(*field)
	//		}
	//	}
	//}
}


func (pb *packetbeat)NewMyCoder(
) (*MyCoder, error) {
        //
	config := &pb.config
	if !config.Flows.IsEnabled() {
		return nil,nil
	}

	processors, err := processors.New(config.Flows.Processors)
	if err != nil {
		return nil, err
	}

	client, err := pb.pipeline.ConnectWith(beat.ClientConfig{
		Processing: beat.ProcessingConfig{
			EventMetadata: config.Flows.EventMetadata,
			Processor:     processors,
			KeepNull:      config.Flows.KeepNull,
		},
	})
	if err != nil {
		return nil, err
	}

        d := MyCoder{
            ss: "MyDecoder",
            client: client}

        return &d, nil
}


func (d *MyCoder) OnPacket(data []byte, ci *gopacket.CaptureInfo) {

    packet_time := fmt.Sprintf("%s", ci.Timestamp.Format("20060102150405"))
    packet_len  := fmt.Sprintf("%d", ci.Length) 
    event := beat.Event{Fields: common.MapStr{}}
    event.PutValue("@mycount", "mycount") 
     
    encodeString := base64.StdEncoding.EncodeToString(data)

    event.PutValue("@base64Packet", encodeString)

    event.PutValue("mylabel", "LABEL#1")
    
    event.PutValue("timestamp", packet_time)
    event.PutValue("packetlen", packet_len)
    event.PutValue("@timestamp", time.Now())
    d.client.Publish(event) 
//    fmt.Printf("MyCoder: -------------------\n%X\n", data)        
    fmt.Printf("packetbeat beater packetbeat.go MyCoder-OnPacket full  packet Success !\n")
}


func (pb *packetbeat) createMyWorker(dl layers.LinkType) (sniffer.Worker, error) {

        
        return  pb.NewMyCoder()
}




func (pb *packetbeat) icmpConfig() (*common.Config, error) {
	var icmp *common.Config
	if pb.config.Protocols["icmp"].Enabled() {
		icmp = pb.config.Protocols["icmp"]
	}

	for _, cfg := range pb.config.ProtocolsList {
		info := struct {
			Type string `config:"type" validate:"required"`
		}{}

		if err := cfg.Unpack(&info); err != nil {
			return nil, err
		}

		if info.Type != "icmp" {
			continue
		}

		if icmp != nil {
			return nil, errors.New("More then one icmp configurations found")
		}

		icmp = cfg
	}

	return icmp, nil
}

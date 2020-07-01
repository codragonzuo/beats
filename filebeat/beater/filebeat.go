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
	"flag"
	"fmt"
	"strings"
	"encoding/base64"

	"github.com/pkg/errors"

	"github.com/codragonzuo/beats/filebeat/channel"
	"github.com/codragonzuo/beats/filebeat/config"
	"github.com/codragonzuo/beats/filebeat/fileset"
	_ "github.com/codragonzuo/beats/filebeat/include"
	"github.com/codragonzuo/beats/filebeat/input"
	"github.com/codragonzuo/beats/filebeat/registrar"
	"github.com/codragonzuo/beats/libbeat/autodiscover"
	"github.com/codragonzuo/beats/libbeat/beat"
	"github.com/codragonzuo/beats/libbeat/cfgfile"
	"github.com/codragonzuo/beats/libbeat/common"
	"github.com/codragonzuo/beats/libbeat/common/cfgwarn"
	"github.com/codragonzuo/beats/libbeat/common/reload"
	"github.com/codragonzuo/beats/libbeat/esleg/eslegclient"
	"github.com/codragonzuo/beats/libbeat/logp"
	"github.com/codragonzuo/beats/libbeat/management"
	"github.com/codragonzuo/beats/libbeat/monitoring"
	"github.com/codragonzuo/beats/libbeat/outputs/elasticsearch"
	"github.com/codragonzuo/beats/libbeat/publisher/pipetool"

	_ "github.com/codragonzuo/beats/filebeat/include"

	// Add filebeat level processors
	_ "github.com/codragonzuo/beats/filebeat/processor/add_kubernetes_metadata"
	_ "github.com/codragonzuo/beats/libbeat/processors/decode_csv_fields"

	// include all filebeat specific builders
	_ "github.com/codragonzuo/beats/filebeat/autodiscover/builder/hints"
	
	
	//"errors"
	// "sync"
	// "time"

    "github.com/tsg/gopacket"
	"github.com/tsg/gopacket/layers"

    
	//"github.com/elastic/beats/libbeat/beat"
	//"github.com/elastic/beats/libbeat/common"
	//"github.com/elastic/beats/libbeat/logp"
	"github.com/codragonzuo/beats/libbeat/processors"
	//"github.com/codragonzuo/beats/libbeat/service"

	//"github.com/elastic/beats/packetbeat/config"
	 "github.com/codragonzuo/beats/packetbeat/decoder"
	 "github.com/codragonzuo/beats/packetbeat/flows"
	 "github.com/codragonzuo/beats/packetbeat/procs"
	"github.com/codragonzuo/beats/packetbeat/protos"
	 "github.com/codragonzuo/beats/packetbeat/protos/icmp"
	 "github.com/codragonzuo/beats/packetbeat/protos/tcp"
	 "github.com/codragonzuo/beats/packetbeat/protos/udp"
	 "github.com/codragonzuo/beats/packetbeat/publish"
	"github.com/codragonzuo/beats/packetbeat/sniffer"

	// Add packetbeat default processors
	_ "github.com/codragonzuo/beats/packetbeat/processor/add_kubernetes_metadata"
)

const pipelinesWarning = "Filebeat is unable to load the Ingest Node pipelines for the configured" +
	" modules because the Elasticsearch output is not configured/enabled. If you have" +
	" already loaded the Ingest Node pipelines or are using Logstash pipelines, you" +
	" can ignore this warning."

var (
	once = flag.Bool("once", false, "Run filebeat only once until all harvesters reach EOF")
)

// Filebeat is a beater object. Contains all objects needed to run the beat
type Filebeat struct {
	config         *config.Config
	moduleRegistry *fileset.ModuleRegistry
	done           chan struct{}
	pipeline       beat.PipelineConnector
	
	cmdLineArgs  config.Flags
	sniff       *sniffer.Sniffer
	// publisher/pipeline
	pipeline2 beat.Pipeline
	transPub *publish.TransactionPublisher
	flows    *flows.Flows
}

// New creates a new Filebeat pointer instance.
func New(b *beat.Beat, rawConfig *common.Config) (beat.Beater, error) {
	fbconfig := config.DefaultConfig
	if err := rawConfig.Unpack(&fbconfig); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
    fmt.Printf("fliebeat beater filebeat.go New()  interfaces=%v\n", fbconfig.Interfaces)
	fmt.Printf("fliebeat beater filebeat.go New()  interfaces=%v\n", fbconfig.Flows)

	if err := cfgwarn.CheckRemoved6xSettings(
		rawConfig,
		"prospectors",
		"config.prospectors",
		"registry_file",
		"registry_file_permissions",
		"registry_flush",
	); err != nil {
		return nil, err
	}

	moduleRegistry, err := fileset.NewModuleRegistry(fbconfig.Modules, b.Info, true)
	if err != nil {
		return nil, err
	}
	if !moduleRegistry.Empty() {
		logp.Info("Enabled modules/filesets: %s", moduleRegistry.InfoString())
	}

	moduleInputs, err := moduleRegistry.GetInputConfigs()
	if err != nil {
		return nil, err
	}

	if err := fbconfig.FetchConfigs(); err != nil {
		return nil, err
	}

	// Add inputs created by the modules
	fbconfig.Inputs = append(fbconfig.Inputs, moduleInputs...)

	enabledInputs := fbconfig.ListEnabledInputs()
	var haveEnabledInputs bool
	if len(enabledInputs) > 0 {
		haveEnabledInputs = true
	}

	if !fbconfig.ConfigInput.Enabled() && !fbconfig.ConfigModules.Enabled() && !haveEnabledInputs && fbconfig.Autodiscover == nil && !b.ConfigManager.Enabled() {
		if !b.InSetupCmd {
			return nil, errors.New("no modules or inputs enabled and configuration reloading disabled. What files do you want me to watch?")
		}

		// in the `setup` command, log this only as a warning
		logp.Warn("Setup called, but no modules enabled.")
	}

	if *once && fbconfig.ConfigInput.Enabled() && fbconfig.ConfigModules.Enabled() {
		return nil, errors.New("input configs and -once cannot be used together")
	}

	if fbconfig.IsInputEnabled("stdin") && len(enabledInputs) > 1 {
		return nil, fmt.Errorf("stdin requires to be run in exclusive mode, configured inputs: %s", strings.Join(enabledInputs, ", "))
	}

	fb := &Filebeat{
		done:           make(chan struct{}),
		config:         &fbconfig,
		moduleRegistry: moduleRegistry,
		cmdLineArgs:    config.CmdLineArgs,
	}

	err = fb.setupPipelineLoaderCallback(b)
	if err != nil {
		return nil, err
	}
	
	//err = fb.init(b)
	//if err != nil {
	//	return nil, err
	//}

	return fb, nil
}

// init packetbeat components
func (fb *Filebeat) init(b *beat.Beat) error {
	var err error
	//fbcfg := &fb.config
	// Enable the process watcher only if capturing live traffic
	if fb.config.Interfaces.File == "" {
		err = procs.ProcWatcher.Init(fb.config.Procs)
		if err != nil {
			logp.Critical(err.Error())
			return err
		}
	} else {
		logp.Info("Process watcher disabled when file input is used")
	}
    fmt.Printf("filebeat filebeat.go init call pipeline\n")
	fb.pipeline2 = b.Publisher
	fb.transPub, err = publish.NewTransactionPublisher(
		b.Info.Name,
		b.Publisher,
		fb.config.IgnoreOutgoing,
		fb.config.Interfaces.File == "",
	)
	if err != nil {
		return err
	}

	logp.Debug("main", "Initializing protocol plugins")
	err = protos.Protos.Init(false, fb.transPub, fb.config.Protocols, fb.config.ProtocolsList)
	if err != nil {
		return fmt.Errorf("Initializing protocol analyzers failed: %v", err)
	}
    fmt.Printf("filebeat filebeat.go init call setupFlows\n")
	if err := fb.setupFlows(); err != nil {
		return err
	}

	return fb.setupSniffer()
}


func (fb *Filebeat) mypatch() () {
        //config := &fb.config
        if !fb.config.Flows.IsEnabled() {
                return 
        }

        processors, err := processors.New(fb.config.Flows.Processors)
        if err != nil {
                return 
        }

        client, err := fb.pipeline2.ConnectWith(beat.ClientConfig{
                Processing: beat.ProcessingConfig{
                        EventMetadata: fb.config.Flows.EventMetadata,
                        Processor:     processors,
                        KeepNull:      fb.config.Flows.KeepNull,
                },
        })

     fmt.Printf("filebeat filebeat.go mypatch client= %v\n", client)

}

func (fb *Filebeat) createWorker(dl layers.LinkType) (sniffer.Worker, error) {
	var icmp4 icmp.ICMPv4Processor
	var icmp6 icmp.ICMPv6Processor
	cfg, err := fb.icmpConfig()
	if err != nil {
		return nil, err
	}
	if cfg.Enabled() {
		reporter, err := fb.transPub.CreateReporter(cfg)
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


        //config := &fb.config
        if !fb.config.Flows.IsEnabled() {
                return nil,nil
        }

        processors, err := processors.New(fb.config.Flows.Processors)
        if err != nil {
                return nil, err
        }

        client, err := fb.pipeline2.ConnectWith(beat.ClientConfig{
                Processing: beat.ProcessingConfig{
                        EventMetadata: fb.config.Flows.EventMetadata,
                        Processor:     processors,
                        KeepNull:      fb.config.Flows.KeepNull,
                },
        })

     fmt.Printf("filebeat filebeat.go createWorker client= %v\n", client)

	worker, err := decoder.New(fb.flows, dl, icmp4, icmp6, tcp, udp, client)
	if err != nil {
		return nil, err
	}

	return worker, nil
}

func (fb *Filebeat) setupSniffer() error {
	//config := &fb.config
    fmt.Printf("filebeat filebeat.go init setupSniffer start\n")
	icmp, err := fb.icmpConfig()
	if err != nil {
		return err
	}

	withVlans := fb.config.Interfaces.WithVlans
	withICMP := icmp.Enabled()
	
	fmt.Printf("filebeat filebeat.go withICMP=%v withVlans=%v\n", withVlans, withICMP)

	filter := fb.config.Interfaces.BpfFilter
	if filter == "" && !fb.config.Flows.IsEnabled() {
		filter = protos.Protos.BpfFilter(withVlans, withICMP)
	}
    fmt.Printf("filebeat filebeat.go init setupSniffer call sniffer New\n")
	//fb.sniff, err = sniffer.New(false, filter, fb.createWorker, fb.config.Interfaces)
	//fb.sniff, err = sniffer.New(false, filter, fb.createWorker,  fb.createMyWorker , fb.config.Interfaces)
	return err
}

func (fb *Filebeat) setupFlows() error {
	//config := &fb.config
	fmt.Printf("filebeat filebeat.go setupFlows start \n")
	fmt.Printf("filebeat filebeat.go setupFlows 1 config Flow=%v\n", fb.config.Flows) 
	if !fb.config.Flows.IsEnabled() {
	    fmt.Printf("SetupFlows return nil")
		return nil
	}
    fmt.Printf("filebeat filebeat.go setupFlows 2 config Flow=%d\n", fb.config.Flows) 
	processors, err := processors.New(fb.config.Flows.Processors)
	if err != nil {
		return err
	}
    fmt.Printf("filebeat filebeat.go  settupFlows call pipeline2 connect\n")
	client, err := fb.pipeline2.ConnectWith(beat.ClientConfig{
		Processing: beat.ProcessingConfig{
			EventMetadata: fb.config.Flows.EventMetadata,
			Processor:     processors,
			KeepNull:      fb.config.Flows.KeepNull,
		},
	})
	fmt.Printf("filebeat filebeat.go  settupFlows client=%v\n", client)
	if err != nil {
		return err
	}

	//fb.flows, err = flows.NewFlows(client.PublishAll, fb.config.Flows)
	if err != nil {
		return err
	}

	return nil
}


type MyCoder struct {
        ss string	
        client beat.Client
}


func (fb *Filebeat) Myworker(ch chan beat.Event, client beat.Client) {
}



func (fb *Filebeat)NewMyCoder(
) (*MyCoder, error) {
        //
	fmt.Printf("filebeat filebeat.go NewMyCoder start 0000000\n")
	//fbconfig := &fb.config
	if !fb.config.Flows.IsEnabled() {
		return nil,nil
	}
    fmt.Printf("1\n")
	processors, err := processors.New(fb.config.Flows.Processors)
	if err != nil {
		return nil, err
	}
    fmt.Printf("2\n")
	client, err := fb.pipeline2.ConnectWith(beat.ClientConfig{
		Processing: beat.ProcessingConfig{
			EventMetadata: fb.config.Flows.EventMetadata,
			Processor:     processors,
			KeepNull:      fb.config.Flows.KeepNull,
		},
	})
	if client == nil	{
	    fmt.Printf("client == nil\n")
	}
	if err != nil {
		return nil, err
	}
    
    d := MyCoder{
            ss: "MyDecoder",
            client: client}
    fmt.Printf("filebeat filebeat.go NewMyCoder end\n")
    
	return &d, nil
}


func (d *MyCoder) OnPacket(data []byte, ci *gopacket.CaptureInfo) {
     
     event := beat.Event{Fields: common.MapStr{}}
     event.PutValue("@mycount", "mycount") 
     
     encodeString := base64.StdEncoding.EncodeToString(data)

     event.PutValue("@base64Packet", encodeString)

     event.PutValue("mylabel", "LABEL#1")
	 //fmt.Printf("Mycoder = %v\n", d)
	 //fmt.Printf("client = %v\n", d.client)
	 //fmt.Printf("client = %v\n", data)
     d.client.Publish(event) 
     //fmt.Printf("MyCoder: %X\n", data)        
}


func (fb *Filebeat) createMyWorker(dl layers.LinkType) (sniffer.Worker, error) {
        Coder, err := fb.NewMyCoder()
		fmt.Printf("Coder=%v", Coder)
        return  Coder, err//fb.NewMyCoder()
}

func (fb *Filebeat) icmpConfig() (*common.Config, error) {
	var icmp *common.Config
	if fb.config.Protocols["icmp"].Enabled() {
		icmp = fb.config.Protocols["icmp"]
	}

	for _, cfg := range fb.config.ProtocolsList {
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

    fmt.Printf("filebeat beater filebeat.go icmp=%v\n", icmp)
	return icmp, nil
}


// setupPipelineLoaderCallback sets the callback function for loading pipelines during setup.
func (fb *Filebeat) setupPipelineLoaderCallback(b *beat.Beat) error {
	if b.Config.Output.Name() != "elasticsearch" {
		logp.Warn(pipelinesWarning)
		return nil
	}

	overwritePipelines := true
	b.OverwritePipelinesCallback = func(esConfig *common.Config) error {
		esClient, err := eslegclient.NewConnectedClient(esConfig)
		if err != nil {
			return err
		}

		// When running the subcommand setup, configuration from modules.d directories
		// have to be loaded using cfg.Reloader. Otherwise those configurations are skipped.
		pipelineLoaderFactory := newPipelineLoaderFactory(b.Config.Output.Config())
		modulesFactory := fileset.NewSetupFactory(b.Info, pipelineLoaderFactory)
		if fb.config.ConfigModules.Enabled() {
			modulesLoader := cfgfile.NewReloader(fb.pipeline, fb.config.ConfigModules)
			modulesLoader.Load(modulesFactory)
		}

		return fb.moduleRegistry.LoadPipelines(esClient, overwritePipelines)
	}
	return nil
}

// loadModulesPipelines is called when modules are configured to do the initial
// setup.
func (fb *Filebeat) loadModulesPipelines(b *beat.Beat) error {
	if b.Config.Output.Name() != "elasticsearch" {
		logp.Warn(pipelinesWarning)
		return nil
	}

	overwritePipelines := fb.config.OverwritePipelines
	if b.InSetupCmd {
		overwritePipelines = true
	}

	// register pipeline loading to happen every time a new ES connection is
	// established
	callback := func(esClient *eslegclient.Connection) error {
		return fb.moduleRegistry.LoadPipelines(esClient, overwritePipelines)
	}
	_, err := elasticsearch.RegisterConnectCallback(callback)

	return err
}

// Run allows the beater to be run as a beat.
func (fb *Filebeat) Run(b *beat.Beat) error {
	var err error
	config := fb.config

    fmt.Printf("filebeat filebeat.go Run start\n")

    

	if !fb.moduleRegistry.Empty() {
		err = fb.loadModulesPipelines(b)
		if err != nil {
			return err
		}
	}

	waitFinished := newSignalWait()
	waitEvents := newSignalWait()

	// count active events for waiting on shutdown
	wgEvents := &eventCounter{
		count: monitoring.NewInt(nil, "filebeat.events.active"),
		added: monitoring.NewUint(nil, "filebeat.events.added"),
		done:  monitoring.NewUint(nil, "filebeat.events.done"),
	}
	finishedLogger := newFinishedLogger(wgEvents)

	// Setup registrar to persist state
	
	
	 
	fmt.Printf("filebeat filebeat.go Run config.Registry=%v\n", config.Registry)
	registrar, err := registrar.New(config.Registry, finishedLogger)
	if err != nil {
		logp.Err("Could not init registrar: %v", err)
		return err
	}

	// Make sure all events that were published in
	registrarChannel := newRegistrarLogger(registrar)
    fmt.Printf("filebeat filebeat.go Run registrarChannel=%v\n", registrarChannel)
    fmt.Printf("filebeat filebeat.go Run  call SetACKHandler begin\n")

	err = b.Publisher.SetACKHandler(beat.PipelineACKHandler{
		ACKEvents: newEventACKer(finishedLogger, registrarChannel).ackEvents,
	})
	fmt.Printf("filebeat filebeat.go Run  call SetACKHandler over\n")
	if err != nil {
	    fmt.Printf("filebeat filebeat.go Run  register error!\n")
		logp.Err("Failed to install the registry with the publisher pipeline: %v", err)
		return err
	}
    fmt.Printf("filebeat filebeat.go Run  call pipeline begin\n")
	
	
	fb.pipeline = pipetool.WithDefaultGuarantees(b.Publisher, beat.GuaranteedSend)
	fb.pipeline = withPipelineEventCounter(fb.pipeline, wgEvents)
    fmt.Printf("filebeat filebeat.go Run  call pipeline end\n")
	
	outDone := make(chan struct{}) // outDone closes down all active pipeline connections
	pipelineConnector := channel.NewOutletFactory(outDone).Create

	// Create a ES connection factory for dynamic modules pipeline loading
	var pipelineLoaderFactory fileset.PipelineLoaderFactory
	if b.Config.Output.Name() == "elasticsearch" {
		pipelineLoaderFactory = newPipelineLoaderFactory(b.Config.Output.Config())
	} else {
		logp.Warn(pipelinesWarning)
	}

	inputLoader := channel.RunnerFactoryWithCommonInputSettings(b.Info,
		input.NewRunnerFactory(pipelineConnector, registrar, fb.done))
	moduleLoader := fileset.NewFactory(inputLoader, b.Info, pipelineLoaderFactory, config.OverwritePipelines)

	crawler, err := newCrawler(inputLoader, moduleLoader, config.Inputs, fb.done, *once)
	if err != nil {
		logp.Err("Could not init crawler: %v", err)
		return err
	}

	// The order of starting and stopping is important. Stopping is inverted to the starting order.
	// The current order is: registrar, publisher, spooler, crawler
	// That means, crawler is stopped first.

	// Start the registrar
	err = registrar.Start()
	if err != nil {
		return fmt.Errorf("Could not start registrar: %v", err)
	}

	
	//begin
	//err = fb.init(b)
	//if err != nil {
	//	return err
	//}
	//fb.mypatch()
	//defer func() {
	//	if service.ProfileEnabled() {
	//		logp.Debug("main", "Waiting for streams and transactions to expire...")
	//		time.Sleep(time.Duration(float64(protos.DefaultTransactionExpiration) * 1.2))
	//		logp.Debug("main", "Streams and transactions should all be expired now.")
	//	}
	//}()

	//defer fb.transPub.Stop()

	//timeout := fb.config.ShutdownTimeout
	//if timeout > 0 {
	//	defer time.Sleep(timeout)
	//}

    fmt.Printf("filebeat beater filebeat.go Run call fb.flows.start()\n")

	//if fb.flows != nil {
	//	fb.flows.Start()
	//	defer fb.flows.Stop()
	//}

	//var wg sync.WaitGroup
	//errC := make(chan error, 1)

	// Run the sniffer in background
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
    //
	//	err := fb.sniff.Run()
	//	if err != nil {
	//		errC <- fmt.Errorf("Sniffer main loop failed: %v", err)
	//	}
	//}()
	

    //end




	// Stopping registrar will write last state
	defer registrar.Stop()

	// Stopping publisher (might potentially drop items)
	defer func() {
		// Closes first the registrar logger to make sure not more events arrive at the registrar
		// registrarChannel must be closed first to potentially unblock (pretty unlikely) the publisher
		registrarChannel.Close()
		close(outDone) // finally close all active connections to publisher pipeline
	}()
	

    

	// Wait for all events to be processed or timeout
	defer waitEvents.Wait()




	if config.OverwritePipelines {
		logp.Debug("modules", "Existing Ingest pipelines will be updated")
	}

	err = crawler.Start(fb.pipeline, config.ConfigInput, config.ConfigModules)
	if err != nil {
		crawler.Stop()
		return fmt.Errorf("Failed to start crawler: %+v", err)
	}

	// If run once, add crawler completion check as alternative to done signal
	if *once {
		runOnce := func() {
			logp.Info("Running filebeat once. Waiting for completion ...")
			crawler.WaitForCompletion()
			logp.Info("All data collection completed. Shutting down.")
		}
		waitFinished.Add(runOnce)
	}

	// Register reloadable list of inputs and modules
	inputs := cfgfile.NewRunnerList(management.DebugK, inputLoader, fb.pipeline)
	reload.Register.MustRegisterList("filebeat.inputs", inputs)

	modules := cfgfile.NewRunnerList(management.DebugK, moduleLoader, fb.pipeline)
	reload.Register.MustRegisterList("filebeat.modules", modules)

	var adiscover *autodiscover.Autodiscover
	if fb.config.Autodiscover != nil {
		adiscover, err = autodiscover.NewAutodiscover(
			"filebeat",
			fb.pipeline,
			cfgfile.MultiplexedRunnerFactory(
				cfgfile.MatchHasField("module", moduleLoader),
				cfgfile.MatchDefault(inputLoader),
			),
			autodiscover.QueryConfig(),
			config.Autodiscover,
			b.Keystore,
		)
		if err != nil {
			return err
		}
	}
	adiscover.Start()

	// Add done channel to wait for shutdown signal
	waitFinished.AddChan(fb.done)
	waitFinished.Wait()

	// Stop reloadable lists, autodiscover -> Stop crawler -> stop inputs -> stop harvesters
	// Note: waiting for crawlers to stop here in order to install wgEvents.Wait
	//       after all events have been enqueued for publishing. Otherwise wgEvents.Wait
	//       or publisher might panic due to concurrent updates.
	inputs.Stop()
	modules.Stop()
	adiscover.Stop()
	crawler.Stop()


    //begin
	//logp.Debug("main", "Waiting for the sniffer to finish")
	//wg.Wait()
	//select {
	//default:
	//case err := <-errC:
	//	return err
	//}
	//end

	timeout := fb.config.ShutdownTimeout
	// Checks if on shutdown it should wait for all events to be published
	waitPublished := fb.config.ShutdownTimeout > 0 || *once
	if waitPublished {
		// Wait for registrar to finish writing registry
		waitEvents.Add(withLog(wgEvents.Wait,
			"Continue shutdown: All enqueued events being published."))
		// Wait for either timeout or all events having been ACKed by outputs.
		if fb.config.ShutdownTimeout > 0 {
			logp.Info("Shutdown output timer started. Waiting for max %v.", timeout)
			waitEvents.Add(withLog(waitDuration(timeout),
				"Continue shutdown: Time out waiting for events being published."))
		} else {
			waitEvents.AddChan(fb.done)
		}
	}

	return nil
}

// Stop is called on exit to stop the crawling, spooling and registration processes.
func (fb *Filebeat) Stop() {
	logp.Info("Stopping filebeat")
	//fb.sniff.Stop()
	// Stop Filebeat
	close(fb.done)
}

// Create a new pipeline loader (es client) factory
func newPipelineLoaderFactory(esConfig *common.Config) fileset.PipelineLoaderFactory {
	pipelineLoaderFactory := func() (fileset.PipelineLoader, error) {
		esClient, err := eslegclient.NewConnectedClient(esConfig)
		if err != nil {
			return nil, errors.Wrap(err, "Error creating Elasticsearch client")
		}
		return esClient, nil
	}
	return pipelineLoaderFactory
}

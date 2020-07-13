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

	"github.com/pkg/errors"

	"github.com/codragonzuo/beats/filebeat/channel"
	cfg "github.com/codragonzuo/beats/filebeat/config"
	"github.com/codragonzuo/beats/filebeat/fileset"
	_ "github.com/codragonzuo/beats/filebeat/include"
	"github.com/codragonzuo/beats/filebeat/input"
	"github.com/codragonzuo/beats/filebeat/registrar"
	"github.com/codragonzuo/beats/libbeat/autodiscover"
	"github.com/codragonzuo/beats/libbeat/beat"
	"github.com/codragonzuo/beats/libbeat/cfgfile"
	"github.com/codragonzuo/beats/libbeat/common"
	//"github.com/codragonzuo/beats/libbeat/common/config"
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
	_ "sync"
	_ "time"

	_ "github.com/tsg/gopacket/layers"

	//"github.com/elastic/beats/libbeat/beat"
	//"github.com/elastic/beats/libbeat/common"
	//"github.com/elastic/beats/libbeat/logp"
	_ "github.com/codragonzuo/beats/libbeat/processors"
	_ "github.com/codragonzuo/beats/libbeat/service"

	//"github.com/elastic/beats/packetbeat/config"
	_ "github.com/codragonzuo/beats/packetbeat/decoder"
	 "github.com/codragonzuo/beats/packetbeat/flows"
	_ "github.com/codragonzuo/beats/packetbeat/procs"
	_ "github.com/codragonzuo/beats/packetbeat/protos"
	_ "github.com/codragonzuo/beats/packetbeat/protos/icmp"
	_ "github.com/codragonzuo/beats/packetbeat/protos/tcp"
	_ "github.com/codragonzuo/beats/packetbeat/protos/udp"
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
	config         *cfg.Config

	cmdLineArgs cfg.Flags
	sniff       *sniffer.Sniffer

	// publisher/pipeline
	pipeline2 beat.Pipeline
	transPub *publish.TransactionPublisher
	flows    *flows.Flows

	moduleRegistry *fileset.ModuleRegistry
	done           chan struct{}
	pipeline       beat.PipelineConnector
}



// New creates a new Filebeat pointer instance.
func New(b *beat.Beat, rawConfig *common.Config) (beat.Beater, error) {
        fmt.Printf("filebeat beater filebeat.go New()  beat.Beat! dragon\n")
        	
        config := cfg.DefaultConfig
	
        fmt.Printf("fliebeat beater filebeat.go New()  interfaces=%v\n", config.Interfaces)
        if err := rawConfig.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
        
        fmt.Printf("get DefaultConfig Success! dragon\n")
        fmt.Printf("fliebeat beater filebeat.go New()  interfaces=%v\n", config.Interfaces)

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

	moduleRegistry, err := fileset.NewModuleRegistry(config.Modules, b.Info, true)
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

	if err := config.FetchConfigs(); err != nil {
		return nil, err
	}

	// Add inputs created by the modules
	config.Inputs = append(config.Inputs, moduleInputs...)

	enabledInputs := config.ListEnabledInputs()
	var haveEnabledInputs bool
	if len(enabledInputs) > 0 {
		haveEnabledInputs = true
	}

	if !config.ConfigInput.Enabled() && !config.ConfigModules.Enabled() && !haveEnabledInputs && config.Autodiscover == nil && !b.ConfigManager.Enabled() {
		if !b.InSetupCmd {
			return nil, errors.New("no modules or inputs enabled and configuration reloading disabled. What files do you want me to watch?")
		}

		// in the `setup` command, log this only as a warning
		logp.Warn("Setup called, but no modules enabled.")
	}

	if *once && config.ConfigInput.Enabled() && config.ConfigModules.Enabled() {
		return nil, errors.New("input configs and -once cannot be used together")
	}

	if config.IsInputEnabled("stdin") && len(enabledInputs) > 1 {
		return nil, fmt.Errorf("stdin requires to be run in exclusive mode, configured inputs: %s", strings.Join(enabledInputs, ", "))
	}

	fb := &Filebeat{
		done:           make(chan struct{}),
		config:         &config,
		moduleRegistry: moduleRegistry,
	}


        fmt.Printf("filebeat beater filebeat.go call setupPipelineLoaderCallback\n")
	err = fb.setupPipelineLoaderCallback(b)
	if err != nil {
		return nil, err
	}

        fmt.Printf("filebeat beater filebeat.go New()  beat.Beat Over ! dragon\n")


	return fb, nil
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
	fmt.Printf("\n------------------------------------------------\nfilebeat beater filebeat.go Run begin dragon\n")
        var err error
	config := fb.config

        
	if !fb.moduleRegistry.Empty() {
	        fmt.Printf("filebeat beater filebeat.go Run call loadModulesPipelines\n")
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
	fmt.Printf("filebeat beater filebeat.go Run call registrar.New\n")
        registrar, err := registrar.New(config.Registry, finishedLogger)
	if err != nil {
		logp.Err("Could not init registrar: %v", err)
		return err
	}

	// Make sure all events that were published in
	registrarChannel := newRegistrarLogger(registrar)

        fmt.Printf("filebeat beater filebeat.go Run call Publisher.SetAckHandler\n")
	err = b.Publisher.SetACKHandler(beat.PipelineACKHandler{
		ACKEvents: newEventACKer(finishedLogger, registrarChannel).ackEvents,
	})
	if err != nil {
		logp.Err("Failed to install the registry with the publisher pipeline: %v", err)
		return err
	}

	fb.pipeline = pipetool.WithDefaultGuarantees(b.Publisher, beat.GuaranteedSend)
        fmt.Printf("filebeat beater filebeat.go Run fb.pipeline=%v\n", fb.pipeline)

	fb.pipeline = withPipelineEventCounter(fb.pipeline, wgEvents)

        fmt.Printf("filebeat beater filebeat.go Run fb.pipeline=%v\n", fb.pipeline)

        fmt.Printf("filebeat beater filebeat.go Run create pipelineConnector\n")
	outDone := make(chan struct{}) // outDone closes down all active pipeline connections
	pipelineConnector := channel.NewOutletFactory(outDone).Create

        fmt.Printf("filebeat beater filebeat.go Run pipelineConnector=%v\n", pipelineConnector)
	//out, err := pipelineConnector.ConnectWith(fb.config, beat.ClientConfig{
	//	Processing: b.ProcessingConfig{
	//		DynamicFields: context.DynamicFields,
	//	},
	//})

        


	// Create a ES connection factory for dynamic modules pipeline loading
	var pipelineLoaderFactory fileset.PipelineLoaderFactory
	if b.Config.Output.Name() == "elasticsearch" {
		pipelineLoaderFactory = newPipelineLoaderFactory(b.Config.Output.Config())
	} else {
		logp.Warn(pipelinesWarning)
	}

        fmt.Printf("filebeat beater filebeat.go Run call RunnerFactoryWithCommonInputSettings------------\n")
	inputLoader := channel.RunnerFactoryWithCommonInputSettings(b.Info,
		input.NewRunnerFactory(pipelineConnector, registrar, fb.done))

        fmt.Printf("filebeat beater filebeat.go Run call fileset.NewFactory-------------\n")
	moduleLoader := fileset.NewFactory(inputLoader, b.Info, pipelineLoaderFactory, config.OverwritePipelines)


        fmt.Printf("filebeat beater filebeat.go Run call newCrawler--------------\n")
	crawler, err := newCrawler(inputLoader, moduleLoader, config.Inputs, fb.done, *once)
	if err != nil {
		logp.Err("Could not init crawler: %v", err)
		return err
	}

	// The order of starting and stopping is important. Stopping is inverted to the starting order.
	// The current order is: registrar, publisher, spooler, crawler
	// That means, crawler is stopped first.

	// Start the registrar
        fmt.Printf("filebeat  beater filebeat.go Run call registrar begin-------------------\n")
	err = registrar.Start()
	if err != nil {
		return fmt.Errorf("Could not start registrar: %v", err)
	}

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


        fmt.Printf("filebeat  beater filebeat.go Run call crawler.Start begin----------------\n")
	err = crawler.Start(fb.pipeline, config.ConfigInput, config.ConfigModules)
	if err != nil {
		crawler.Stop()
		return fmt.Errorf("Failed to start crawler: %+v", err)
	}
        fmt.Printf("filebeat  beater filebeat.go Run call crawler.Start end----------------\n")
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
        fmt.Printf("filebeat  beater filebeat.go Run call cfgfile.NewRunnerList begin------------------\n")
	inputs := cfgfile.NewRunnerList(management.DebugK, inputLoader, fb.pipeline)
	reload.Register.MustRegisterList("filebeat.inputs", inputs)
        fmt.Printf("filebeat  beater filebeat.go Run call cfgfile.NewRunnerList end----------------------\n")

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

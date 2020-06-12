// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package azureeventhub

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/codragonzuo/beats/libbeat/logp"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/stretchr/testify/assert"

	"github.com/codragonzuo/beats/filebeat/channel"
	"github.com/codragonzuo/beats/libbeat/beat"
	"github.com/codragonzuo/beats/libbeat/common"
)

var (
	config = azureInputConfig{
		SAKey:            "",
		SAName:           "",
		SAContainer:      ephContainerName,
		ConnectionString: "",
		ConsumerGroup:    "",
	}
)

func TestProcessEvents(t *testing.T) {
	// Stub outlet for receiving events generated by the input.
	o := &stubOutleter{}
	out, err := newStubOutlet(o)
	if err != nil {
		t.Fatal(err)
	}
	input := azureInput{
		config: config,
		outlet: out,
	}
	var sn int64 = 12
	now := time.Now()
	var off int64 = 1234
	var pID int16 = 1

	properties := eventhub.SystemProperties{
		SequenceNumber: &sn,
		EnqueuedTime:   &now,
		Offset:         &off,
		PartitionID:    &pID,
		PartitionKey:   nil,
	}
	single := "{\"test\":\"this is some message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}"
	msg := fmt.Sprintf("{\"records\":[%s]}", single)
	ev := eventhub.Event{
		Data:             []byte(msg),
		SystemProperties: &properties,
	}
	ok := input.processEvents(&ev, "0")
	if !ok {
		t.Fatal("OnEvent function returned false")
	}
	assert.Equal(t, len(o.Events), 1)
	message, err := o.Events[0].Fields.GetValue("message")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, message, single)
}

func TestParseMultipleMessages(t *testing.T) {
	// records object
	msg := "{\"records\":[{\"test\":\"this is some message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}," +
		"{\"test\":\"this is 2nd message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}," +
		"{\"test\":\"this is 3rd message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}]}"
	msgs := []string{
		fmt.Sprintf("{\"test\":\"this is some message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}"),
		fmt.Sprintf("{\"test\":\"this is 2nd message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}"),
		fmt.Sprintf("{\"test\":\"this is 3rd message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}")}
	input := azureInput{log: logp.NewLogger(fmt.Sprintf("%s test for input", inputName))}
	messages := input.parseMultipleMessages([]byte(msg))
	assert.NotNil(t, messages)
	assert.Equal(t, len(messages), 3)
	for _, ms := range messages {
		assert.Contains(t, msgs, ms)
	}

	// array of events
	msg1 := "[{\"test\":\"this is some message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}," +
		"{\"test\":\"this is 2nd message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}," +
		"{\"test\":\"this is 3rd message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}]"
	messages = input.parseMultipleMessages([]byte(msg1))
	assert.NotNil(t, messages)
	assert.Equal(t, len(messages), 3)
	for _, ms := range messages {
		assert.Contains(t, msgs, ms)
	}

	// one event only
	msg2 := "{\"test\":\"this is some message\",\"time\":\"2019-12-17T13:43:44.4946995Z\"}"
	messages = input.parseMultipleMessages([]byte(msg2))
	assert.NotNil(t, messages)
	assert.Equal(t, len(messages), 1)
	for _, ms := range messages {
		assert.Contains(t, msgs, ms)
	}
}

type stubOutleter struct {
	sync.Mutex
	cond   *sync.Cond
	done   bool
	Events []beat.Event
}

func newStubOutlet(stub *stubOutleter) (channel.Outleter, error) {
	stub.cond = sync.NewCond(stub)
	defer stub.Close()

	connector := channel.ConnectorFunc(func(_ *common.Config, _ beat.ClientConfig) (channel.Outleter, error) {
		return stub, nil
	})
	return connector.ConnectWith(nil, beat.ClientConfig{
		Processing: beat.ProcessingConfig{},
	})
}
func (o *stubOutleter) Close() error {
	o.Lock()
	defer o.Unlock()
	o.done = true
	return nil
}
func (o *stubOutleter) Done() <-chan struct{} { return nil }
func (o *stubOutleter) OnEvent(event beat.Event) bool {
	o.Lock()
	defer o.Unlock()
	o.Events = append(o.Events, event)
	o.cond.Broadcast()
	return o.done
}

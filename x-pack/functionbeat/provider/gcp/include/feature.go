// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package include

import (
	"github.com/codragonzuo/beats/libbeat/feature"
	"github.com/codragonzuo/beats/x-pack/functionbeat/function/provider"
	"github.com/codragonzuo/beats/x-pack/functionbeat/provider/gcp/gcp"
)

// Bundle exposes the trigger supported by the GCP provider.
var bundle = provider.MustCreate(
	"gcp",
	provider.NewDefaultProvider("gcp", provider.NewNullCli, provider.NewNullTemplateBuilder),
	feature.MakeDetails("Google Cloud Platform", "listen to events from Google Cloud Platform", feature.Stable),
).MustAddFunction("pubsub",
	gcp.NewPubSub,
	gcp.PubSubDetails(),
).MustAddFunction("storage",
	gcp.NewStorage,
	gcp.StorageDetails(),
).Bundle()

func init() {
	feature.MustRegisterBundle(bundle)
}

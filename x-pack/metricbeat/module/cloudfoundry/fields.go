// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package cloudfoundry

import (
	"github.com/codragonzuo/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "cloudfoundry", asset.ModuleFieldsPri, AssetCloudfoundry); err != nil {
		panic(err)
	}
}

// AssetCloudfoundry returns asset data.
// This is the base64 encoded gzipped contents of module/cloudfoundry.
func AssetCloudfoundry() string {
	return "eJzcVbFyGjEQ7fmKN25oYj7gihRxJjPuXDh1LKQ90KCTLtIKcn+fkeDwHQhDAqTIFcywkt57eqvdfcSKugrSuKhqF63y3QRgzYYqPDwNwg8TQFGQXresna3weQIAeQu+bfegcSoamgCeDIlAFebEYgLUmowKVT7yCCsaOiJNH3ctVVh4F9tdpEA5RhsipuP7YI+2om7jvBrEi5jb73VJ+RhcDVqTZdTeNeNbzkYnXlwIem4Ia2EiBWgrTVSEqXSWhbbkp5/Sn2iZ/BTCKkzz1unsSL9o2yP5QzMuEC/a1mgp0jJ4SWiIvZbQASIEJ7VgUthoXg5vcejnUJNWo/BpV8+I6wU+f03e8ljqsRV7964x5G2P8tYDhnfknTfhTIoPn/JFptnAwkr6oa2iX0UDjbOLP3LvOUH13vUEwxzPKWEGsJsVRck2zlrJRTW1ceJw5Uwyn16+IwaxILTkJVkWCyoTN9Q4383mHVO4jRf4krCSGTGQ6gk+Iv8ZHYv7SBBroY1ILeAjHUqH1R0tCOz8yQRk7n/iwJGM93rODfDKas4Yo1rOkTtWcvq9bQNMiH0Z7/SfyBqZkUxckapErHRdk6fUNebEG6LthDAiMFg3NFS0G35Oyug9qbJAdizMDQVmvO0cRe182aGePG+76jllhP/oMeX7lDMVrS73/b8mTogXEB9mCddMnNee64D4dwAAAP//IdTbYA=="
}

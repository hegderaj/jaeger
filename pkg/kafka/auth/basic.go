// Copyright (c) 2019 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	gotls "crypto/tls"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

// BasicConfig describes the configuration properties for Basic Auth Connections to the Kafka Brokers
type BasicConfig struct {
	CaPath   string
	Username string
	Password string
}

func setBasicConfiguration(config *BasicConfig, saramaConfig *sarama.Config) error {
	// If SSL is enabled with basic auth to avoid username and passwords as plain text
	// inf not provided then authentication happens with plain text
	if config.CaPath != "" {
		tlsConfig, err := config.getTLS()
		if err != nil {
			return errors.Wrap(err, "error loading tls config")
		}
		saramaConfig.Net.TLS.Enable = true
		saramaConfig.Net.TLS.Config = tlsConfig
	}
	// Set the SASL config
	saramaConfig.Net.SASL.User = config.Username
	saramaConfig.Net.SASL.Password = config.Password
	saramaConfig.Net.SASL.Enable = true
	return nil
}

func (basicConfig BasicConfig) getTLS() (*gotls.Config, error) {
	ca, err := loadCA(basicConfig.CaPath)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading ca")
	}
	return &gotls.Config{
		RootCAs: ca,
	}, nil
}

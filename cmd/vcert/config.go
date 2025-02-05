/*
 * Copyright 2018-2023 Venafi, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/Venafi/vcert/v5"
	"github.com/Venafi/vcert/v5/pkg/endpoint"
	"github.com/Venafi/vcert/v5/pkg/venafi"
)

func buildConfig(c *cli.Context, flags *commandFlags) (cfg vcert.Config, err error) {
	cfg.LogVerbose = flags.verbose

	if flags.config != "" {
		// Loading configuration from file
		cfg, err = vcert.LoadConfigFromFile(flags.config, flags.profile)
		if err != nil {
			return cfg, err
		}
	} else {
		// Loading configuration from CLI flags
		var connectorType endpoint.ConnectorType
		var baseURL string
		var auth = &endpoint.Authentication{}
		var identityProvider = &endpoint.OAuthProvider{}

		//case when access token can come from environment variable.
		tokenS := flags.token

		if tokenS == "" {
			tokenS = getPropertyFromEnvironment(vCertToken)
		}

		if flags.testMode {
			connectorType = endpoint.ConnectorTypeFake
			if flags.testModeDelay > 0 {
				logf("Running in -test-mode with emulating endpoint delay.")
				delay, _ := rand.Int(rand.Reader, big.NewInt(int64(flags.testModeDelay)))
				for i := int64(0); i < delay.Int64(); i++ {
					time.Sleep(1 * time.Second)
				}
			}
		} else if flags.platform == venafi.Firefly || (flags.userName != "" || tokenS != "" || flags.clientP12 != "" || c.Command.Name == "sshgetconfig") {

			if flags.platform == venafi.Firefly {
				connectorType = endpoint.ConnectorTypeFirefly
			} else {
				connectorType = endpoint.ConnectorTypeTPP
			}

			//add support for using environment variables begins
			baseURL = flags.url
			if baseURL == "" {
				baseURL = getPropertyFromEnvironment(vCertURL)
			}
			//add support for using environment variables ends

			if connectorType != endpoint.ConnectorTypeFirefly && tokenS == "" && flags.password == "" && flags.clientP12 == "" && c.Command.Name != "sshgetconfig" {
				return cfg, fmt.Errorf("A password is required to communicate with TPP")
			}

			if flags.token != "" {
				if c.Command.Name == commandGetCredName {
					auth.RefreshToken = flags.token
				} else {
					auth.AccessToken = flags.token
				}
			} else if flags.userName != "" && flags.password != "" {
				auth.User = flags.userName
				auth.Password = flags.password
			} else {
				tokenS := getPropertyFromEnvironment(vCertToken)
				if tokenS != "" {
					if c.Command.Name == commandGetCredName {
						auth.RefreshToken = tokenS
					} else {
						auth.AccessToken = tokenS
					}
				}
			}

			if flags.platform == venafi.Firefly && c.Command.Name == commandGetCredName {
				auth.ClientId = flags.clientId
				auth.ClientSecret = flags.clientSecret
				identityProvider.TokenURL = flags.tokenURL
				identityProvider.DeviceURL = flags.deviceURL
				identityProvider.Audience = flags.audience
				auth.IdentityProvider = identityProvider
				auth.Scope = flags.scope
			}
		} else {
			apiKey := flags.apiKey
			if apiKey == "" {
				apiKey = getPropertyFromEnvironment(vCertApiKey)
			}
			connectorType = endpoint.ConnectorTypeCloud
			baseURL = flags.url
			auth.APIKey = apiKey
			if flags.email != "" {
				auth.User = flags.email
				auth.Password = flags.password
			}
		}
		cfg.ConnectorType = connectorType
		cfg.Credentials = auth
		cfg.BaseUrl = baseURL
	}

	// trust bundle may be overridden by CLI flag
	if flags.trustBundle != "" {
		logf("Detected trust bundle flag at CLI.")
		if cfg.ConnectionTrust != "" {
			logf("Overriding trust bundle based on command line flag.")
		}
		data, err := os.ReadFile(flags.trustBundle)
		if err != nil {
			return cfg, fmt.Errorf("Failed to read trust bundle: %s", err)
		}
		cfg.ConnectionTrust = string(data)
	} else {
		trustBundleSrc := getPropertyFromEnvironment(vCertTrustBundle)
		if trustBundleSrc != "" {
			logf("Detected trust bundle in environment properties.")
			if cfg.ConnectionTrust != "" {
				logf("Overriding trust bundle based on environment property")
			}
			data, err := os.ReadFile(trustBundleSrc)
			if err != nil {
				return cfg, fmt.Errorf("Failed to read trust bundle: %s", err)
			}
			cfg.ConnectionTrust = string(data)
		}
	}

	// zone may be overridden by CLI flag
	if flags.zone != "" {
		if cfg.Zone != "" {
			logf("Overriding zone based on command line flag.")
		}
		cfg.Zone = flags.zone
	}

	zone := getPropertyFromEnvironment(vCertZone)
	if cfg.Zone == "" && zone != "" {
		cfg.Zone = zone
	}

	if c.Command.Name == commandEnrollName || c.Command.Name == commandPickupName {
		if cfg.Zone == "" && cfg.ConnectorType != endpoint.ConnectorTypeFake && !(flags.pickupID != "" || flags.pickupIDFile != "") {
			return cfg, fmt.Errorf("Zone cannot be empty. Use -z option")
		}
	}

	return cfg, nil
}

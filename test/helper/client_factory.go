/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/plugin"
	"testing"
)

// Creates a new in-memory DevLake server with default settings and returns a client to it
func StartDevLakeServer(t *testing.T, loadedGoPlugins []plugin.PluginMeta) *DevlakeClient {
	client := ConnectLocalServer(t, &LocalClientConfig{
		ServerPort:   8089,
		DbURL:        config.GetConfig().GetString("E2E_DB_URL"),
		CreateServer: true,
		DropDb:       false,
		TruncateDb:   true,
		Plugins:      loadedGoPlugins,
	})
	return client
}

// Connect to an existing DevLake server with default config. Tables are truncated. Useful for troubleshooting outside the IDE.
func ConnectDevLakeServer(t *testing.T) *DevlakeClient {
	client := ConnectRemoteServer(t, &RemoteClientConfig{
		Endpoint:   "http://localhost:8089",
		DbURL:      config.GetConfig().GetString("E2E_DB_URL"),
		TruncateDb: true,
	})
	return client
}

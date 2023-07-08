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

package api

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

const (
	TimeOut = 10 * time.Second
)

func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.GithubConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	apiClient, err := helper.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.Get(input.Params["path"], input.Query, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := errors.Convert01(io.ReadAll(resp.Body))
	if err != nil {
		return nil, err
	}
	// verify response body is json
	var tmp interface{}
	err = errors.Convert(json.Unmarshal(body, &tmp))
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Status: resp.StatusCode, Body: json.RawMessage(body)}, nil
}

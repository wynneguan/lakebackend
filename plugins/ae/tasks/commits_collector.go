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

package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"io"
	"net/http"
	"net/url"
)

const RAW_COMMITS_TABLE = "ae_commits"

func CollectCommits(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*AeTaskData)
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: AeApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_COMMITS_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    2000,
		UrlTemplate: "projects/{{ .Params.ProjectId }}/commits",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, errors.Default.Wrap(err, "error reading endpoint response by AE commit collector")
			}
			var results []json.RawMessage
			err = errors.Convert(json.Unmarshal(body, &results))
			return results, errors.Default.Wrap(err, "error deserializing endpoint response by AE commit collector")
		},
	})

	if err != nil {
		return err
	}
	return collector.Execute()
}

var CollectCommitsMeta = plugin.SubTaskMeta{
	Name:             "collectCommits",
	EntryPoint:       CollectCommits,
	EnabledByDefault: true,
	Description:      "Collect commit analysis data from AE api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

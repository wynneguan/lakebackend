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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
	"net/url"
)

const RAW_BUG_CUSTOM_FIELDS_TABLE = "tapd_api_bug_custom_fields"

var _ plugin.SubTaskEntryPoint = CollectBugCustomFields

func CollectBugCustomFields(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_CUSTOM_FIELDS_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect bug_custom_fields")
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		UrlTemplate:        "bugs/custom_fields_settings",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				BugCustomFields []json.RawMessage `json:"data"`
			}
			err := api.UnmarshalResponse(res, &data)
			return data.BugCustomFields, err
		},
	})
	if err != nil {
		logger.Error(err, "collect bug_custom_fields error")
		return err
	}
	return collector.Execute()
}

var CollectBugCustomFieldsMeta = plugin.SubTaskMeta{
	Name:             "collectBugCustomFields",
	EntryPoint:       CollectBugCustomFields,
	EnabledByDefault: true,
	Description:      "collect Tapd BugCustomFields",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

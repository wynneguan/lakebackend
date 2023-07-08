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

const RAW_STORY_CUSTOM_FIELDS_TABLE = "tapd_api_story_custom_fields"

var _ plugin.SubTaskEntryPoint = CollectStoryCustomFields

func CollectStoryCustomFields(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_CUSTOM_FIELDS_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect story_custom_fields")
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		UrlTemplate:        "stories/custom_fields_settings",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				StoryCustomFields []json.RawMessage `json:"data"`
			}
			err := api.UnmarshalResponse(res, &data)
			return data.StoryCustomFields, err
		},
	})
	if err != nil {
		logger.Error(err, "collect story_custom_fields error")
		return err
	}
	return collector.Execute()
}

var CollectStoryCustomFieldsMeta = plugin.SubTaskMeta{
	Name:             "collectStoryCustomFields",
	EntryPoint:       CollectStoryCustomFields,
	EnabledByDefault: true,
	Description:      "collect Tapd StoryCustomFields",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

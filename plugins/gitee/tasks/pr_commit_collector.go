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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"net/http"
	"net/url"
	"reflect"
)

const RAW_PULL_REQUEST_COMMIT_TABLE = "gitee_api_pull_request_commits"

var CollectApiPullRequestCommitsMeta = plugin.SubTaskMeta{
	Name:             "collectApiPullRequestCommits",
	EntryPoint:       CollectApiPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Collect PullRequestCommits data from Gitee api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

type SimplePr struct {
	Number  int
	GiteeId int
}

func CollectApiPullRequestCommits(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_COMMIT_TABLE)

	incremental := false

	cursor, err := db.Cursor(
		dal.Select("number, gitee_id"),
		dal.From(models.GiteePullRequest{}.TableName()),
		dal.Where("repo_id = ? and connection_id=?", data.Repo.GiteeId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimplePr{}))
	if err != nil {
		return err
	}
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        incremental,
		Input:              iterator,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/pulls/{{ .Input.Number }}/commits",

		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("state", "all")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var items []json.RawMessage
			err := api.UnmarshalResponse(res, &items)
			if err != nil {
				return nil, err
			}
			return items, nil
		},
		AfterResponse: ignoreHTTPStatus404,
	})

	if err != nil {
		return err
	}
	return collector.Execute()
}

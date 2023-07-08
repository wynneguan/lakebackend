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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"net/url"
	"strconv"
)

const RAW_COMMIT_TABLE = "gitee_api_commit"

var CollectCommitsMeta = plugin.SubTaskMeta{
	Name:             "collectApiCommits",
	EntryPoint:       CollectApiCommits,
	EnabledByDefault: true,
	Description:      "Collect commit data from gitee api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE, plugin.DOMAIN_TYPE_CROSS},
}

func CollectApiCommits(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)
	since := data.Since
	incremental := false
	if since == nil {
		latestUpdated := &models.GiteeCommit{}
		err := db.All(
			&latestUpdated,
			dal.Join("left join _tool_gitee_repo_commits on _tool_gitee_commits.sha = _tool_gitee_repo_commits.commit_sha"),
			dal.Join("left join _tool_gitee_repos on _tool_gitee_repo_commits.repo_id = _tool_gitee_repos.gitee_id"),
			dal.Where("_tool_gitee_repo_commits.repo_id = ? AND _tool_gitee_repo_commits.connection_id = ?", data.Repo.GiteeId, data.Repo.ConnectionId),
			dal.Orderby("committed_date DESC"),
			dal.Limit(1),
		)

		if err != nil {
			return errors.Default.Wrap(err, "failed to get latest gitee commit record")
		}
		if latestUpdated.Sha != "" {
			since = &latestUpdated.CommittedDate
			incremental = true
		}
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        incremental,
		UrlTemplate:        "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/commits",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			if since != nil {
				query.Set("since", since.String())
			}
			query.Set("page", strconv.Itoa(reqData.Pager.Page))
			query.Set("per_page", strconv.Itoa(reqData.Pager.Size))
			return query, nil
		},
		Concurrency:    20,
		ResponseParser: GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

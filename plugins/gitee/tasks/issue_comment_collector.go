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
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"net/url"
)

const RAW_COMMENTS_TABLE = "gitee_issue_comments"

var CollectApiIssueCommentsMeta = plugin.SubTaskMeta{
	Name:             "collectApiIssueComments",
	EntryPoint:       CollectApiIssueComments,
	EnabledByDefault: true,
	Description:      "Collect comments data from Gitee api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectApiIssueComments(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMENTS_TABLE)

	since := data.Since
	incremental := false
	// user didn't specify a time range to sync, try load from database
	// actually, for gitee pull, since doesn't make any sense, gitee pull api doesn't support it
	if since == nil {
		var latestUpdatedIssueComment models.GiteeIssueComment
		err := db.All(
			&latestUpdatedIssueComment,
			dal.Join("left join _tool_gitee_issues on _tool_gitee_issues.gitee_id = _tool_gitee_issue_comments.issue_id"),
			dal.Where(
				"_tool_gitee_issues.repo_id = ? AND _tool_gitee_issues.connection_id = ?", data.Repo.GiteeId, data.Repo.ConnectionId,
			),
			dal.Orderby("gitee_updated_at DESC"),
			dal.Limit(1),
		)
		if err != nil {
			return errors.Default.Wrap(err, "failed to get latest gitee issue record")
		}
		var latestUpdatedPrComt models.GiteePullRequestComment
		err = db.All(
			&latestUpdatedPrComt,
			dal.Join("left join _tool_gitee_pull_requests on _tool_gitee_pull_requests.gitee_id = _tool_gitee_pull_request_comments.pull_request_id"),
			dal.Where("_tool_gitee_pull_requests.repo_id = ? AND _tool_gitee_pull_requests.connection_id = ?", data.Repo.GiteeId, data.Repo.ConnectionId),
			dal.Orderby("gitee_updated_at DESC"),
			dal.Limit(1),
		)
		if err != nil {
			return errors.Default.Wrap(err, "failed to get latest gitee issue record")
		}
		if latestUpdatedIssueComment.GiteeId > 0 && latestUpdatedPrComt.GiteeId > 0 {
			if latestUpdatedIssueComment.GiteeUpdatedAt.Before(latestUpdatedPrComt.GiteeUpdatedAt) {
				since = &latestUpdatedPrComt.GiteeUpdatedAt
			} else {
				since = &latestUpdatedIssueComment.GiteeUpdatedAt
			}
			incremental = true
		} else if latestUpdatedIssueComment.GiteeId > 0 {
			since = &latestUpdatedIssueComment.GiteeUpdatedAt
			incremental = true
		} else if latestUpdatedPrComt.GiteeId > 0 {
			since = &latestUpdatedPrComt.GiteeUpdatedAt
			incremental = true
		}
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        incremental,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues/comments",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("state", "all")
			if since != nil {
				query.Set("since", since.String())
			}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		GetTotalPages:  GetTotalPagesFromResponse,
		ResponseParser: GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

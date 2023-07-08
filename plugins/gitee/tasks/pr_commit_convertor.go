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
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"reflect"
)

var ConvertPullRequestCommitsMeta = plugin.SubTaskMeta{
	Name:             "convertPullRequestCommits",
	EntryPoint:       ConvertPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitee_pull_request_commits into  domain layer table pull_request_commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

func ConvertPullRequestCommits(taskCtx plugin.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_COMMIT_TABLE)
	repoId := data.Repo.GiteeId

	pullIdGen := didgen.NewDomainIdGenerator(&models.GiteePullRequest{})

	cursor, err := db.Cursor(
		dal.From(&models.GiteePullRequestCommit{}),
		dal.Join(`left join _tool_gitee_pull_requests on _tool_gitee_pull_requests.gitee_id = _tool_gitee_pull_request_commits.pull_request_id`),
		dal.Where("_tool_gitee_pull_requests.repo_id = ? and _tool_gitee_pull_requests.connection_id = ?", repoId, data.Options.ConnectionId),
		dal.Orderby("pull_request_id ASC"),
	)

	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.GiteePullRequestCommit{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			giteePullRequestCommit := inputRow.(*models.GiteePullRequestCommit)
			domainPrCommit := &code.PullRequestCommit{
				CommitSha:          giteePullRequestCommit.CommitSha,
				PullRequestId:      pullIdGen.Generate(data.Options.ConnectionId, giteePullRequestCommit.PullRequestId),
				CommitAuthorName:   giteePullRequestCommit.CommitAuthorName,
				CommitAuthorEmail:  giteePullRequestCommit.CommitAuthorEmail,
				CommitAuthoredDate: giteePullRequestCommit.CommitAuthoredDate,
			}
			return []interface{}{
				domainPrCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

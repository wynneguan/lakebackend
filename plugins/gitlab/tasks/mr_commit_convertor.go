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
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"reflect"
)

func init() {
	RegisterSubtaskMeta(&ConvertApiMrCommitsMeta)
}

var ConvertApiMrCommitsMeta = plugin.SubTaskMeta{
	Name:             "convertApiMergeRequestsCommits",
	EntryPoint:       ConvertApiMergeRequestsCommits,
	EnabledByDefault: true,
	Description:      "Add domain layer PullRequestCommit according to GitlabMrCommit",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertMrCommentMeta},
}

func ConvertApiMergeRequestsCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_COMMITS_TABLE)
	db := taskCtx.GetDal()

	clauses := []dal.Clause{
		dal.From(&models.GitlabMrCommit{}),
		dal.Join(`left join _tool_gitlab_merge_requests
			on _tool_gitlab_merge_requests.gitlab_id =
			_tool_gitlab_mr_commits.merge_request_id`),
		dal.Where(`_tool_gitlab_merge_requests.project_id = ?
			and _tool_gitlab_merge_requests.connection_id = ?`,
			data.Options.ProjectId, data.Options.ConnectionId),
		dal.Orderby("merge_request_id ASC"),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}

	// TODO: adopt batch indate operation
	domainIdGenerator := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMrCommit{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			GitlabMrCommit := inputRow.(*models.GitlabMrCommit)
			domainPrcommit := &code.PullRequestCommit{
				CommitSha:          GitlabMrCommit.CommitSha,
				PullRequestId:      domainIdGenerator.Generate(data.Options.ConnectionId, GitlabMrCommit.MergeRequestId),
				CommitAuthorName:   GitlabMrCommit.CommitAuthorName,
				CommitAuthorEmail:  GitlabMrCommit.CommitAuthorEmail,
				CommitAuthoredDate: *GitlabMrCommit.CommitAuthoredDate,
			}
			return []interface{}{
				domainPrcommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

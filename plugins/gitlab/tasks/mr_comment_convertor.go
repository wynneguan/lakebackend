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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"reflect"
)

func init() {
	RegisterSubtaskMeta(&ConvertMrCommentMeta)
}

var ConvertMrCommentMeta = plugin.SubTaskMeta{
	Name:             "convertMergeRequestComment",
	EntryPoint:       ConvertMergeRequestComment,
	EnabledByDefault: true,
	Description:      "Add domain layer Comment according to GitlabMrComment",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertApiMergeRequestsMeta},
}

func ConvertMergeRequestComment(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.From(&models.GitlabMrComment{}),
		dal.Join(`left join _tool_gitlab_merge_requests on
			_tool_gitlab_merge_requests.gitlab_id =
			_tool_gitlab_mr_comments.merge_request_id`),
		dal.Where(`_tool_gitlab_merge_requests.project_id = ?
			and _tool_gitlab_mr_comments.connection_id = ?`,
			data.Options.ProjectId, data.Options.ConnectionId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	domainIdGeneratorComment := didgen.NewDomainIdGenerator(&models.GitlabMrComment{})
	prIdGen := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.GitlabAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMrComment{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			gitlabComments := inputRow.(*models.GitlabMrComment)

			domainComment := &code.PullRequestComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: domainIdGeneratorComment.Generate(data.Options.ConnectionId, gitlabComments.GitlabId),
				},
				PullRequestId: prIdGen.Generate(data.Options.ConnectionId, gitlabComments.MergeRequestId),
				Body:          gitlabComments.Body,
				AccountId:     accountIdGen.Generate(data.Options.ConnectionId, gitlabComments.AuthorUserId),
				CreatedDate:   gitlabComments.GitlabCreatedAt,
			}
			domainComment.Type = getStdCommentType(gitlabComments.Type)
			if domainComment.Body == "unapproved this merge request" {
				domainComment.Status = "CHANGES_REQUESTED"
			}
			if domainComment.Body == "approved this merge request" {
				domainComment.Status = "APPROVED"
			}
			return []interface{}{
				domainComment,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

func getStdCommentType(originType string) string {
	if originType == "DiffNote" {
		return code.DIFF_COMMENT
	}
	if originType == "REVIEW" {
		return code.REVIEW
	}
	return code.NORMAL_COMMENT
}

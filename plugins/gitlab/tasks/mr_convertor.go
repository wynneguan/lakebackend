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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertApiMergeRequestsMeta)
}

var ConvertApiMergeRequestsMeta = plugin.SubTaskMeta{
	Name:             "convertApiMergeRequests",
	EntryPoint:       ConvertApiMergeRequests,
	EnabledByDefault: true,
	Description:      "Add domain layer PullRequest according to GitlabMergeRequest",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertProjectMeta},
}

func ConvertApiMergeRequests(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_TABLE)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.From(&models.GitlabMergeRequest{}),
		dal.Where("project_id=? and connection_id = ?", data.Options.ProjectId, data.Options.ConnectionId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	domainMrIdGenerator := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})
	domainRepoIdGenerator := didgen.NewDomainIdGenerator(&models.GitlabProject{})
	domainUserIdGen := didgen.NewDomainIdGenerator(&models.GitlabAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMergeRequest{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			gitlabMr := inputRow.(*models.GitlabMergeRequest)

			domainPr := &code.PullRequest{
				DomainEntity: domainlayer.DomainEntity{
					Id: domainMrIdGenerator.Generate(data.Options.ConnectionId, gitlabMr.GitlabId),
				},
				HeadRepoId:     domainRepoIdGenerator.Generate(data.Options.ConnectionId, gitlabMr.SourceProjectId),
				BaseRepoId:     domainRepoIdGenerator.Generate(data.Options.ConnectionId, gitlabMr.TargetProjectId),
				OriginalStatus: gitlabMr.State,
				PullRequestKey: gitlabMr.Iid,
				Title:          gitlabMr.Title,
				Description:    gitlabMr.Description,
				Type:           gitlabMr.Type,
				Url:            gitlabMr.WebUrl,
				AuthorName:     gitlabMr.AuthorUsername,
				AuthorId:       domainUserIdGen.Generate(data.Options.ConnectionId, gitlabMr.AuthorUserId),
				CreatedDate:    gitlabMr.GitlabCreatedAt,
				MergedDate:     gitlabMr.MergedAt,
				ClosedDate:     gitlabMr.ClosedAt,
				MergeCommitSha: retrieveMrSha(gitlabMr.MergeCommitSha, gitlabMr.SquashCommitSha, gitlabMr.DiffHeadSha),
				HeadRef:        gitlabMr.SourceBranch,
				BaseRef:        gitlabMr.TargetBranch,
				Component:      gitlabMr.Component,
			}
			switch gitlabMr.State {
			case "opened":
				domainPr.Status = code.OPEN
			case "merged":
				domainPr.Status = code.MERGED
			case "closed", "locked":
				domainPr.Status = code.CLOSED
			default:
				domainPr.Status = gitlabMr.State
			}

			return []interface{}{
				domainPr,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

func retrieveMrSha(mergeCommitSha string, squashCommitSha string, sha string) string {
	if mergeCommitSha != "" {
		return mergeCommitSha
	}
	if squashCommitSha != "" {
		return squashCommitSha
	}
	return sha
}

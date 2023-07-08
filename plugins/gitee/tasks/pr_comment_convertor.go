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
	"github.com/apache/incubator-devlake/plugins/gitee/models"
	"reflect"
)

var ConvertPullRequestCommentsMeta = plugin.SubTaskMeta{
	Name:             "convertPullRequestComments",
	EntryPoint:       ConvertPullRequestComments,
	EnabledByDefault: true,
	Description:      "ConvertPullRequestComments data from Gitee api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

func ConvertPullRequestComments(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMENTS_TABLE)
	repoId := data.Repo.GiteeId

	cursor, err := db.Cursor(
		dal.From(&models.GiteePullRequestComment{}),
		dal.Join("left join _tool_gitee_pull_requests "+
			"on _tool_gitee_pull_requests.gitee_id = _tool_gitee_pull_request_comments.pull_request_id"),
		dal.Where("repo_id = ? and _tool_gitee_pull_requests.connection_id = ?", repoId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	prIdGen := didgen.NewDomainIdGenerator(&models.GiteePullRequest{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.GiteeAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.GiteePullRequestComment{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			giteePullRequestComment := inputRow.(*models.GiteePullRequestComment)
			domainPrComment := &code.PullRequestComment{
				DomainEntity: domainlayer.DomainEntity{
					Id: prIdGen.Generate(data.Options.ConnectionId, giteePullRequestComment.GiteeId),
				},
				PullRequestId: prIdGen.Generate(data.Options.ConnectionId, giteePullRequestComment.PullRequestId),
				Body:          giteePullRequestComment.Body,
				AccountId:     accountIdGen.Generate(data.Options.ConnectionId, giteePullRequestComment.AuthorUserId),
				CreatedDate:   giteePullRequestComment.GiteeCreatedAt,
				CommitSha:     "",
				Position:      0,
			}
			return []interface{}{
				domainPrComment,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

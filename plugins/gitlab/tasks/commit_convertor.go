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
	RegisterSubtaskMeta(&ConvertCommitsMeta)
}

var ConvertCommitsMeta = plugin.SubTaskMeta{
	Name:             "convertApiCommits",
	EntryPoint:       ConvertApiCommits,
	EnabledByDefault: false,
	Description:      "Update domain layer commit according to GitlabCommit",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertMrLabelsMeta},
}

func ConvertApiCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)
	db := taskCtx.GetDal()

	// select all commits belongs to the project
	clauses := []dal.Clause{
		dal.Select("gc.*"),
		dal.From("_tool_gitlab_commits gc"),
		dal.Join(`left join _tool_gitlab_project_commits gpc on (
			gpc.commit_sha = gc.sha
		)`),
		dal.Where("gpc.gitlab_project_id = ? and gpc.connection_id = ? ",
			data.Options.ProjectId, data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	// TODO: adopt batch indate operation
	//userDidGen := didgen.NewDomainIdGenerator(&models.GitlabAccount{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabCommit{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			gitlabCommit := inputRow.(*models.GitlabCommit)

			// convert commit
			commit := &code.Commit{}
			commit.Sha = gitlabCommit.Sha
			commit.Message = gitlabCommit.Message
			commit.Additions = gitlabCommit.Additions
			commit.Deletions = gitlabCommit.Deletions
			commit.AuthorId = gitlabCommit.AuthorEmail
			commit.AuthorName = gitlabCommit.AuthorName
			commit.AuthorEmail = gitlabCommit.AuthorEmail
			commit.AuthoredDate = gitlabCommit.AuthoredDate
			commit.CommitterName = gitlabCommit.CommitterName
			commit.CommitterEmail = gitlabCommit.CommitterEmail
			commit.CommittedDate = gitlabCommit.CommittedDate
			commit.CommitterId = gitlabCommit.CommitterEmail

			// convert repo / commits relationship
			repoCommit := &code.RepoCommit{
				RepoId:    didgen.NewDomainIdGenerator(&models.GitlabProject{}).Generate(data.Options.ConnectionId, data.Options.ProjectId),
				CommitSha: gitlabCommit.Sha,
			}

			return []interface{}{
				commit,
				repoCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

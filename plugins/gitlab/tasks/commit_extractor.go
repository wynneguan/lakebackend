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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiCommitsMeta)
}

var ExtractApiCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractApiCommits",
	EntryPoint:       ExtractApiCommits,
	EnabledByDefault: false,
	Description:      "Extract raw commit data into tool layer table GitlabCommit,GitlabAccount and GitlabProjectCommit",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	Dependencies:     []*plugin.SubTaskMeta{&CollectApiCommitsMeta},
}

func ExtractApiCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			// need to extract 3 kinds of entities here
			results := make([]interface{}, 0, 3)

			// create gitlab commit
			gitlabApiCommit := &GitlabApiCommit{}
			err := errors.Convert(json.Unmarshal(row.Data, gitlabApiCommit))
			if err != nil {
				return nil, err
			}
			gitlabCommit, err := ConvertCommit(gitlabApiCommit)
			if err != nil {
				return nil, err
			}

			// create project/commits relationship
			gitlabProjectCommit := &models.GitlabProjectCommit{GitlabProjectId: data.Options.ProjectId}
			gitlabProjectCommit.CommitSha = gitlabCommit.Sha

			// create gitlab user
			GitlabAccountAuthor := &models.GitlabAccount{}
			GitlabAccountAuthor.Email = gitlabCommit.AuthorEmail
			GitlabAccountAuthor.Name = gitlabCommit.AuthorName

			gitlabProjectCommit.ConnectionId = data.Options.ConnectionId
			GitlabAccountAuthor.ConnectionId = data.Options.ConnectionId
			results = append(results, gitlabCommit)
			results = append(results, gitlabProjectCommit)
			results = append(results, GitlabAccountAuthor)

			// For Commiter Email is not same as AuthorEmail
			if gitlabCommit.CommitterEmail != GitlabAccountAuthor.Email {
				gitlabAccountCommitter := &models.GitlabAccount{}
				gitlabAccountCommitter.Email = gitlabCommit.CommitterEmail
				gitlabAccountCommitter.Name = gitlabCommit.CommitterName
				gitlabAccountCommitter.ConnectionId = data.Options.ConnectionId
				results = append(results, gitlabAccountCommitter)
			}

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

// Convert the API response to our DB model instance
func ConvertCommit(commit *GitlabApiCommit) (*models.GitlabCommit, errors.Error) {
	gitlabCommit := &models.GitlabCommit{
		Sha:            commit.GitlabId,
		Title:          commit.Title,
		Message:        commit.Message,
		ShortId:        commit.ShortId,
		AuthorName:     commit.AuthorName,
		AuthorEmail:    commit.AuthorEmail,
		AuthoredDate:   commit.AuthoredDate.ToTime(),
		CommitterName:  commit.CommitterName,
		CommitterEmail: commit.CommitterEmail,
		CommittedDate:  commit.CommittedDate.ToTime(),
		WebUrl:         commit.WebUrl,
		Additions:      commit.Stats.Additions,
		Deletions:      commit.Stats.Deletions,
		Total:          commit.Stats.Total,
	}
	return gitlabCommit, nil
}

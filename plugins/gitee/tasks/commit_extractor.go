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
	"github.com/apache/incubator-devlake/plugins/gitee/models"
)

var ExtractCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractApiCommits",
	EntryPoint:       ExtractApiCommits,
	EnabledByDefault: true,
	Description:      "Extract raw commit data into tool layer table GiteeCommit,GiteeAccount and GiteeRepoCommit",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE, plugin.DOMAIN_TYPE_CROSS},
}

type GiteeCommit struct {
	Author struct {
		Date  api.Iso8601Time `json:"date"`
		Email string          `json:"email"`
		Name  string          `json:"name"`
	}
	Committer struct {
		Date  api.Iso8601Time `json:"date"`
		Email string          `json:"email"`
		Name  string          `json:"name"`
	}
	Message string `json:"message"`
}

type GiteeApiCommitResponse struct {
	Author      *models.GiteeAccount `json:"author"`
	CommentsUrl string               `json:"comments_url"`
	Commit      GiteeCommit          `json:"commit"`
	Committer   *models.GiteeAccount `json:"committer"`
	HtmlUrl     string               `json:"html_url"`
	Sha         string               `json:"sha"`
	Url         string               `json:"url"`
}

func ExtractApiCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_COMMIT_TABLE)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			results := make([]interface{}, 0, 4)

			commit := &GiteeApiCommitResponse{}

			err := errors.Convert(json.Unmarshal(row.Data, commit))

			if err != nil {
				return nil, err
			}

			if commit.Sha == "" {
				return nil, nil
			}

			giteeCommit, err := ConvertCommit(commit)

			if err != nil {
				return nil, err
			}

			if commit.Author != nil {
				giteeCommit.AuthorId = commit.Author.Id
				results = append(results, commit.Author)
			}
			if commit.Committer != nil {
				giteeCommit.CommitterId = commit.Committer.Id
				results = append(results, commit.Committer)
			}

			giteeRepoCommit := &models.GiteeRepoCommit{
				ConnectionId: data.Options.ConnectionId,
				RepoId:       data.Repo.GiteeId,
				CommitSha:    commit.Sha,
			}
			results = append(results, giteeCommit)
			results = append(results, giteeRepoCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

// ConvertCommit Convert the API response to our DB model instance
func ConvertCommit(commit *GiteeApiCommitResponse) (*models.GiteeCommit, errors.Error) {
	giteeCommit := &models.GiteeCommit{
		Sha:            commit.Sha,
		Message:        commit.Commit.Message,
		AuthorName:     commit.Commit.Author.Name,
		AuthorEmail:    commit.Commit.Author.Email,
		AuthoredDate:   commit.Commit.Author.Date.ToTime(),
		CommitterName:  commit.Commit.Author.Name,
		CommitterEmail: commit.Commit.Author.Email,
		CommittedDate:  commit.Commit.Author.Date.ToTime(),
		WebUrl:         commit.Url,
	}
	return giteeCommit, nil
}

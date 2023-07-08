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
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

var _ plugin.SubTaskEntryPoint = ExtractBugCommits

var ExtractBugCommitMeta = plugin.SubTaskMeta{
	Name:             "extractBugCommits",
	EntryPoint:       ExtractBugCommits,
	EnabledByDefault: true,
	Description:      "Extract raw BugCommits data into tool layer table _tool_tapd_issue_commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func ExtractBugCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_COMMIT_TABLE)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var issueCommitBody models.TapdBugCommit
			err := errors.Convert(json.Unmarshal(row.Data, &issueCommitBody))
			if err != nil {
				return nil, err
			}
			toolL := issueCommitBody
			toolL.ConnectionId = data.Options.ConnectionId
			issue := models.Input{}
			err = errors.Convert(json.Unmarshal(row.Input, &issue))
			if err != nil {
				return nil, err
			}
			toolL.BugId = issue.IssueId
			toolL.IssueUpdated = issue.UpdateTime
			toolL.WorkspaceId = data.Options.WorkspaceId
			results := make([]interface{}, 0, 1)
			results = append(results, &toolL)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

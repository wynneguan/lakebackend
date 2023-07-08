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
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
)

var EnrichPullRequestIssuesMeta = plugin.SubTaskMeta{
	Name:             "enrichPullRequestIssues",
	EntryPoint:       EnrichPullRequestIssues,
	EnabledByDefault: true,
	Description:      "Create tool layer table gitee_pull_request_issues from gitee_pull_reqeusts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func EnrichPullRequestIssues(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PULL_REQUEST_TABLE)
	repoId := data.Repo.GiteeId

	var prBodyCloseRegex *regexp.Regexp
	var err errors.Error
	prBodyClosePattern := data.Options.ScopeConfig.PrBodyClosePattern
	//the pattern before the issue number, sometimes, the issue number is #1098, sometimes it is https://xxx/#1098
	prBodyClosePattern = strings.Replace(prBodyClosePattern, "%s", data.Options.Owner, 1)
	prBodyClosePattern = strings.Replace(prBodyClosePattern, "%s", data.Options.Repo, 1)
	if len(prBodyClosePattern) > 0 {
		prBodyCloseRegex, err = errors.Convert01(regexp.Compile(prBodyClosePattern))
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile prBodyClosePattern failed")
		}
	}
	charPattern := regexp.MustCompile(`[a-zA-Z\s,]+`)
	cursor, err := db.Cursor(dal.From(&models.GiteePullRequest{}), dal.Where("repo_id = ?", repoId))
	if err != nil {
		return err
	}
	defer cursor.Close()
	// iterate all rows

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.GiteePullRequest{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			giteePullRequst := inputRow.(*models.GiteePullRequest)
			results := make([]interface{}, 0, 1)

			//find the matched string in body
			issueNumberListStr := ""

			if prBodyCloseRegex != nil {
				issueNumberListStr = prBodyCloseRegex.FindString(giteePullRequst.Body)
			}

			if issueNumberListStr == "" {
				return nil, nil
			}

			issueNumberListStr = charPattern.ReplaceAllString(issueNumberListStr, "#")
			//split the string by '#'
			issueNumberList := strings.Split(issueNumberListStr, "#")

			for _, issueNumberStr := range issueNumberList {
				issue := &models.GiteeIssue{}
				issueNumberStr = strings.TrimSpace(issueNumberStr)
				//change the issueNumberStr to int, if cannot be changed, just continue
				issueNumber, numFormatErr := strconv.Atoi(issueNumberStr)
				if numFormatErr != nil {
					continue
				}
				err = db.All(
					issue,
					dal.Where("number = ? and repo_id = ? and connection_id = ?", issueNumber, repoId, data.Options.ConnectionId),
					dal.Limit(1),
				)
				if err != nil {
					return nil, err
				}
				if issue.Number == "" {
					continue
				}
				giteePullRequstIssue := &models.GiteePullRequestIssue{
					ConnectionId:      data.Options.ConnectionId,
					PullRequestId:     giteePullRequst.GiteeId,
					IssueId:           issue.GiteeId,
					PullRequestNumber: giteePullRequst.Number,
					IssueNumber:       issue.Number,
				}
				results = append(results, giteePullRequstIssue)
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

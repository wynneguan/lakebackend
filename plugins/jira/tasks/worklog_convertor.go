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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"reflect"
)

var ConvertWorklogsMeta = plugin.SubTaskMeta{
	Name:             "convertWorklogs",
	EntryPoint:       ConvertWorklogs,
	EnabledByDefault: true,
	Description:      "convert Jira work logs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertWorklogs(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("convert worklog")
	// select all worklogs belongs to the board
	clauses := []dal.Clause{
		dal.From(&models.JiraWorklog{}),
		dal.Select("_tool_jira_worklogs.*"),
		dal.Join(`LEFT JOIN _tool_jira_board_issues
              ON _tool_jira_board_issues.connection_id = _tool_jira_worklogs.connection_id
                   AND _tool_jira_board_issues.issue_id = _tool_jira_worklogs.issue_id`),
		dal.Where("_tool_jira_board_issues.connection_id = ? AND _tool_jira_board_issues.board_id = ?", connectionId, boardId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		logger.Error(err, "convert worklog error")
		return err
	}
	defer cursor.Close()

	worklogIdGen := didgen.NewDomainIdGenerator(&models.JiraWorklog{})
	accountIdGen := didgen.NewDomainIdGenerator(&models.JiraAccount{})
	issueIdGen := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_WORKLOGS_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraWorklog{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			jiraWorklog := inputRow.(*models.JiraWorklog)
			worklog := &ticket.IssueWorklog{
				DomainEntity:     domainlayer.DomainEntity{Id: worklogIdGen.Generate(jiraWorklog.ConnectionId, jiraWorklog.IssueId, jiraWorklog.WorklogId)},
				IssueId:          issueIdGen.Generate(jiraWorklog.ConnectionId, jiraWorklog.IssueId),
				TimeSpentMinutes: jiraWorklog.TimeSpentSeconds / 60,
				StartedDate:      &jiraWorklog.Started,
				LoggedDate:       &jiraWorklog.Updated,
			}
			if jiraWorklog.AuthorId != "" {
				worklog.AuthorId = accountIdGen.Generate(connectionId, jiraWorklog.AuthorId)
			}
			return []interface{}{worklog}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

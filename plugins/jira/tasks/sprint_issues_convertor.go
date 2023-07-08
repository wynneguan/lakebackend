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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"reflect"
)

var ConvertSprintIssuesMeta = plugin.SubTaskMeta{
	Name:             "convertSprintIssues",
	EntryPoint:       ConvertSprintIssues,
	EnabledByDefault: true,
	Description:      "convert Jira sprint_issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertSprintIssues(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JiraTaskData)

	jiraSprintIssue := &models.JiraSprintIssue{}
	// select all issues belongs to the board
	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From(jiraSprintIssue),
		dal.Where("_tool_jira_sprint_issues.connection_id = ? ", data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	sprintIdGen := didgen.NewDomainIdGenerator(&models.JiraSprint{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.JiraSprintIssue{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_ISSUE_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			jiraSprintIssue := inputRow.(*models.JiraSprintIssue)
			sprintIssue := &ticket.SprintIssue{
				SprintId: sprintIdGen.Generate(data.Options.ConnectionId, jiraSprintIssue.SprintId),
				IssueId:  issueIdGen.Generate(data.Options.ConnectionId, jiraSprintIssue.IssueId),
			}
			return []interface{}{sprintIssue}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

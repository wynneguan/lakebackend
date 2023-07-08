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

package e2e

import (
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/jira/impl"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
)

func TestLabelDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	dataflowTester.FlushTabler(&ticket.IssueLabel{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_board_issues_for_changelog.csv", &models.JiraBoardIssue{})
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_jira_issue_labels_for_convertor.csv", &models.JiraIssueLabel{})
	dataflowTester.Subtask(tasks.ConvertIssueLabelsMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueLabel{},
		"./snapshot_tables/issue_labels.csv",
		e2ehelper.ColumnWithRawData(
			"issue_id",
			"label_name",
		),
	)
}

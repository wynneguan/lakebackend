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
	"github.com/apache/incubator-devlake/plugins/tapd/impl"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
)

func TestTapdTaskChangelogDataFlow(t *testing.T) {

	var tapd impl.Tapd
	dataflowTester := e2ehelper.NewDataFlowTester(t, "tapd", tapd)

	taskData := &tasks.TapdTaskData{
		Options: &tasks.TapdOptions{
			ConnectionId: 1,
			WorkspaceId:  991,
			ScopeConfig: &tasks.TapdScopeConfig{
				TypeMappings: tasks.TypeMappings{
					"Techstory": "REQUIREMENT",
					"技术债":       "REQUIREMENT",
					"需求":        "REQUIREMENT",
				},
				StatusMappings: tasks.StatusMappings{
					"已关闭":                   "DONE",
					"接受/处理":                 "IN_PROGRESS",
					"开发中":                   "IN_PROGRESS",
					"developing":            "IN_PROGRESS",
					"test-11test-11test-12": "IN_PROGRESS",
					"新":                     "TODO",
					"planning":              "TODO",
					"test-11test-11test-11": "TODO",
				},
			},
		},
	}
	// iteration
	dataflowTester.ImportCsvIntoTabler("./snapshot_tables/_tool_tapd_iterations.csv", &models.TapdIteration{})

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_task_changelogs.csv",
		"_raw_tapd_api_task_changelogs")

	// verify extraction
	dataflowTester.FlushTabler(&models.TapdTaskChangelog{})
	dataflowTester.FlushTabler(&models.TapdTaskChangelogItem{})
	dataflowTester.Subtask(tasks.ExtractTaskChangelogMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdTaskChangelog{},
		"./snapshot_tables/_tool_tapd_task_changelogs.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"id",
			"workspace_id",
			"workitem_type_id",
			"creator",
			"created",
			"change_summary",
			"comment",
			"entity_type",
			"change_type",
			"task_id",
		),
	)

	dataflowTester.VerifyTable(
		models.TapdTaskChangelogItem{},
		"./snapshot_tables/_tool_tapd_task_changelog_items.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"changelog_id",
			"field",
			"value_before_parsed",
			"value_after_parsed",
			"iteration_id_from",
			"iteration_id_to",
		),
	)

	dataflowTester.FlushTabler(&ticket.IssueChangelogs{})
	dataflowTester.Subtask(tasks.ConvertTaskChangelogMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.IssueChangelogs{},
		"./snapshot_tables/issue_changelogs_task.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"issue_id",
			"author_id",
			"author_name",
			"field_id",
			"field_name",
			"from_value",
			"to_value",
			"created_date",
			"original_from_value",
			"original_to_value",
		),
	)
}

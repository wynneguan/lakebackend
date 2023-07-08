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
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/jira/impl"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"testing"
)

func TestProjectDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_projects.csv", "_raw_jira_api_projects")

	// verify issue extraction
	dataflowTester.FlushTabler(&models.JiraProject{})
	dataflowTester.Subtask(tasks.ExtractProjectsMeta, taskData)
	dataflowTester.VerifyTable(
		models.JiraProject{},
		"./snapshot_tables/_tool_jira_projects.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"id",
			"project_key",
			"name",
		),
	)
}

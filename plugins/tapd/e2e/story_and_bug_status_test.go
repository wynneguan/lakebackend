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
	"github.com/apache/incubator-devlake/plugins/tapd/impl"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"github.com/apache/incubator-devlake/plugins/tapd/tasks"
	"testing"
)

func TestTapdStoryAndBugStatusDataFlow(t *testing.T) {

	var tapd impl.Tapd
	dataflowTester := e2ehelper.NewDataFlowTester(t, "tapd", tapd)

	taskData := &tasks.TapdTaskData{
		Options: &tasks.TapdOptions{
			ConnectionId: 1,
			WorkspaceId:  991,
		},
	}
	// story status
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_story_status.csv",
		"_raw_tapd_api_story_status")
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_story_status_last_steps.csv",
		"_raw_tapd_api_story_status_last_steps")
	// verify extraction
	dataflowTester.FlushTabler(&models.TapdStoryStatus{})
	dataflowTester.Subtask(tasks.ExtractStoryStatusMeta, taskData)
	dataflowTester.Subtask(tasks.EnrichStoryStatusLastStepMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdStoryStatus{},
		"./snapshot_tables/_tool_tapd_story_statuses.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"workspace_id",
			"english_name",
			"chinese_name",
			"is_last_step",
		),
	)

	// bug status
	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_bug_status.csv",
		"_raw_tapd_api_bug_status")
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_tapd_api_bug_status_last_steps.csv",
		"_raw_tapd_api_bug_status_last_steps")
	// verify extraction
	dataflowTester.FlushTabler(&models.TapdBugStatus{})
	dataflowTester.Subtask(tasks.ExtractBugStatusMeta, taskData)
	dataflowTester.Subtask(tasks.EnrichBugStatusLastStepMeta, taskData)
	dataflowTester.VerifyTable(
		models.TapdBugStatus{},
		"./snapshot_tables/_tool_tapd_bug_statuses.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"workspace_id",
			"english_name",
			"chinese_name",
			"is_last_step",
		),
	)
}

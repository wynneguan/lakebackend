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

	"github.com/apache/incubator-devlake/plugins/github/models"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestEventDataFlow(t *testing.T) {
	var plugin impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "github", plugin)

	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Name:         "panjf2000/ants",
			GithubId:     134018330,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_events.csv", "_raw_github_api_events")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubIssueEvent{})
	dataflowTester.FlushTabler(&models.GithubRepoAccount{})
	dataflowTester.Subtask(tasks.ExtractApiEventsMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubIssueEvent{},
		"./snapshot_tables/_tool_github_issue_events.csv",
		[]string{
			"connection_id",
			"github_id",
			"issue_id",
			"type",
			"author_username",
			"github_created_at",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
	dataflowTester.VerifyTable(
		models.GithubRepoAccount{},
		"./snapshot_tables/_tool_github_accounts_in_event.csv",
		[]string{
			"connection_id",
			"account_id",
			"repo_github_id",
			"login",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)
}

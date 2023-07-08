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
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/jira/impl"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"testing"
)

func TestJiraAccountDataFlow(t *testing.T) {
	var plugin impl.Jira
	dataflowTester := e2ehelper.NewDataFlowTester(t, "jira", plugin)

	taskData := &tasks.JiraTaskData{
		Options: &tasks.JiraOptions{
			ConnectionId: 2,
			BoardId:      8,
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_jira_api_users.csv", "_raw_jira_api_users")

	// verify issue extraction
	dataflowTester.FlushTabler(&models.JiraAccount{})
	dataflowTester.Subtask(tasks.ExtractAccountsMeta, taskData)
	dataflowTester.VerifyTable(
		models.JiraAccount{},
		"./snapshot_tables/_tool_jira_accounts.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"account_id",
			"account_type",
			"name",
			"email",
			"avatar_url",
			"timezone",
		),
	)

	// verify board conversion
	dataflowTester.FlushTabler(&crossdomain.Account{})
	dataflowTester.Subtask(tasks.ConvertAccountsMeta, taskData)
	dataflowTester.VerifyTable(
		crossdomain.Account{},
		"./snapshot_tables/accounts.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"email",
			"full_name",
			"user_name",
			"avatar_url",
			"organization",
			"created_date",
			"status",
		),
	)

}

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

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/bitbucket/impl"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/bitbucket/tasks"
)

func TestIssueDataFlow(t *testing.T) {
	var plugin impl.Bitbucket
	dataflowTester := e2ehelper.NewDataFlowTester(t, "bitbucket", plugin)

	taskData := &tasks.BitbucketTaskData{
		Options: &tasks.BitbucketOptions{
			ConnectionId: 1,
			FullName:     "likyh/likyhphp",
			BitbucketScopeConfig: &models.BitbucketScopeConfig{
				IssueStatusTodo:       "new,open,wantfix",
				IssueStatusInProgress: "on hold,resolved",
				IssueStatusDone:       "closed,invalid",
			},
		},
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_bitbucket_api_issues.csv", "_raw_bitbucket_api_issues")

	// verify issue extraction
	dataflowTester.FlushTabler(&models.BitbucketIssue{})
	dataflowTester.FlushTabler(&models.BitbucketAccount{})
	dataflowTester.Subtask(tasks.ExtractApiIssuesMeta, taskData)
	dataflowTester.VerifyTable(
		models.BitbucketIssue{},
		"./snapshot_tables/_tool_bitbucket_issues.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"bitbucket_id",
			"repo_id",
			"number",
			"std_state",
			"state",
			"title",
			"body",
			"priority",
			"type",
			"author_id",
			"author_name",
			"assignee_id",
			"assignee_name",
			"milestone_id",
			"lead_time_minutes",
			"url",
			"closed_at",
			"severity",
			"component",
		),
	)

	dataflowTester.VerifyTable(
		models.BitbucketAccount{},
		"./snapshot_tables/_tool_bitbucket_accounts_in_issue.csv",
		e2ehelper.ColumnWithRawData(
			"connection_id",
			"user_name",
			"account_id",
			"account_status",
			"display_name",
			"avatar_url",
			"html_url",
			"uuid",
			"has2_fa_enabled",
		),
	)

	// verify issue conversion
	dataflowTester.FlushTabler(&ticket.Issue{})
	dataflowTester.FlushTabler(&ticket.BoardIssue{})
	dataflowTester.FlushTabler(&ticket.IssueAssignee{})
	dataflowTester.Subtask(tasks.ConvertIssuesMeta, taskData)
	dataflowTester.VerifyTable(
		ticket.Issue{},
		"./snapshot_tables/issues.csv",
		e2ehelper.ColumnWithRawData(
			"id",
			"url",
			"icon_url",
			"issue_key",
			"title",
			"description",
			"epic_key",
			"type",
			"status",
			"original_status",
			"story_point",
			"resolution_date",
			"lead_time_minutes",
			"parent_issue_id",
			"priority",
			"original_estimate_minutes",
			"time_spent_minutes",
			"time_remaining_minutes",
			"creator_id",
			"creator_name",
			"assignee_id",
			"assignee_name",
			"severity",
			"component",
		),
	)
	dataflowTester.VerifyTable(
		ticket.BoardIssue{},
		"./snapshot_tables/board_issues.csv",
		e2ehelper.ColumnWithRawData(
			"board_id",
			"issue_id",
		),
	)
	dataflowTester.VerifyTableWithOptions(ticket.IssueAssignee{}, e2ehelper.TableOptions{
		CSVRelPath:  "./snapshot_tables/issue_assignees.csv",
		IgnoreTypes: []interface{}{common.NoPKModel{}},
	})
}

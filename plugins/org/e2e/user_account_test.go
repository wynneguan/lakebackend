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
	"github.com/apache/incubator-devlake/plugins/org/impl"
	"github.com/apache/incubator-devlake/plugins/org/tasks"
	"testing"
)

func TestUserAccountDataFlow(t *testing.T) {
	var plugin impl.Org
	dataflowTester := e2ehelper.NewDataFlowTester(t, "org", plugin)

	taskData := &tasks.TaskData{
		Options: &tasks.Options{
			ConnectionId: 2,
		},
	}

	// import raw data table
	dataflowTester.FlushTabler(&crossdomain.User{})
	dataflowTester.FlushTabler(&crossdomain.Account{})
	dataflowTester.FlushTabler(&crossdomain.UserAccount{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/users.csv", &crossdomain.User{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/accounts.csv", &crossdomain.Account{})
	dataflowTester.ImportCsvIntoTabler("./raw_tables/user_accounts.csv", &crossdomain.UserAccount{})

	dataflowTester.Subtask(tasks.ConnectUserAccountsExactMeta, taskData)
	dataflowTester.VerifyTable(
		crossdomain.UserAccount{},
		"./snapshot_tables/user_accounts.csv",
		e2ehelper.ColumnWithRawData(
			"user_id",
			"account_id",
		),
	)
}

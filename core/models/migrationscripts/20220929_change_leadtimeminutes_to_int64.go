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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*changeLeadTimeMinutesToInt64)(nil)

type changeLeadTimeMinutesToInt64 struct{}

type issues20220929 struct {
	LeadTimeMinutes int64
}

func (issues20220929) TableName() string {
	return "issues"
}

func (script *changeLeadTimeMinutesToInt64) Up(basicRes context.BasicRes) errors.Error {
	// Yes, issues.lead_time_minutes might be negative, we ought to change the type
	// for the column from `uint` to `int64`
	// related issue: https://github.com/apache/incubator-devlake/issues/3224
	db := basicRes.GetDal()
	return migrationhelper.ChangeColumnsType[issues20220929](
		basicRes,
		script,
		issues20220929{}.TableName(),
		[]string{"lead_time_minutes"},
		func(tmpColumnParams []interface{}) errors.Error {
			return db.UpdateColumn(
				&issues20220929{},
				"lead_time_minutes",
				dal.DalClause{Expr: " ? ", Params: tmpColumnParams},
				dal.Where("? != 0", tmpColumnParams...),
			)
		},
	)
}

func (*changeLeadTimeMinutesToInt64) Version() uint64 {
	return 20220929145125
}

func (*changeLeadTimeMinutesToInt64) Name() string {
	return "modify lead_time_minutes"
}

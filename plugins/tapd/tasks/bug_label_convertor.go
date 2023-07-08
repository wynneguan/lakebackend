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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"
)

var ConvertBugLabelsMeta = plugin.SubTaskMeta{
	Name:             "convertBugLabels",
	EntryPoint:       ConvertBugLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table tapd_issue_labels into  domain layer table issue_labels",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertBugLabels(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_TABLE)
	clauses := []dal.Clause{
		dal.From(&models.TapdBugLabel{}),
		dal.Join("left join _tool_tapd_workspace_bugs on _tool_tapd_workspace_bugs.bug_id = _tool_tapd_bug_labels.bug_id"),
		dal.Where("_tool_tapd_workspace_bugs.workspace_id = ? and _tool_tapd_workspace_bugs.connection_id = ?",
			data.Options.WorkspaceId, data.Options.ConnectionId),
		dal.Orderby("bug_id ASC"),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	bugIdGen := didgen.NewDomainIdGenerator(&models.TapdBug{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TapdBugLabel{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			issueLabel := inputRow.(*models.TapdBugLabel)
			domainBugLabel := &ticket.IssueLabel{
				IssueId:   bugIdGen.Generate(issueLabel.BugId),
				LabelName: issueLabel.LabelName,
			}
			return []interface{}{
				domainBugLabel,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

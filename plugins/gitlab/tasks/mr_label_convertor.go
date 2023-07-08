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
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"reflect"
)

func init() {
	RegisterSubtaskMeta(&ConvertMrLabelsMeta)
}

var ConvertMrLabelsMeta = plugin.SubTaskMeta{
	Name:             "convertMrLabels",
	EntryPoint:       ConvertMrLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_mr_labels into  domain layer table pull_request_labels",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertIssueLabelsMeta},
}

func ConvertMrLabels(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_TABLE)
	projectId := data.Options.ProjectId
	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From(&models.GitlabMrLabel{}),
		dal.Join(`left join _tool_gitlab_merge_requests on
			_tool_gitlab_merge_requests.gitlab_id = _tool_gitlab_mr_labels.mr_id`),
		dal.Where(`_tool_gitlab_merge_requests.project_id = ?
			and _tool_gitlab_merge_requests.connection_id = ?`,
			projectId, data.Options.ConnectionId),
		dal.Orderby("mr_id ASC"),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	mrIdGen := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMrLabel{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			mrLabel := inputRow.(*models.GitlabMrLabel)
			domainIssueLabel := &code.PullRequestLabel{
				PullRequestId: mrIdGen.Generate(data.Options.ConnectionId, mrLabel.MrId),
				LabelName:     mrLabel.LabelName,
			}
			return []interface{}{
				domainIssueLabel,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

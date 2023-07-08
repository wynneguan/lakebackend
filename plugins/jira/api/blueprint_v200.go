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

package api

import (
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

func MakeDataSourcePipelinePlanV200(subtaskMetas []plugin.SubTaskMeta, connectionId uint64, bpScopes []*plugin.BlueprintScopeV200, syncPolicy *plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	plan := make(plugin.PipelinePlan, len(bpScopes))
	plan, err := makeDataSourcePipelinePlanV200(subtaskMetas, plan, bpScopes, connectionId, syncPolicy)
	if err != nil {
		return nil, nil, err
	}
	scopes, err := makeScopesV200(bpScopes, connectionId)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	plan plugin.PipelinePlan,
	bpScopes []*plugin.BlueprintScopeV200,
	connectionId uint64,
	syncPolicy *plugin.BlueprintSyncPolicy,
) (plugin.PipelinePlan, errors.Error) {
	for i, bpScope := range bpScopes {
		stage := plan[i]
		if stage == nil {
			stage = plugin.PipelineStage{}
		}
		// construct task options for Jira
		options := make(map[string]interface{})
		options["scopeId"] = bpScope.Id
		options["connectionId"] = connectionId
		if syncPolicy.TimeAfter != nil {
			options["timeAfter"] = syncPolicy.TimeAfter.Format(time.RFC3339)
		}

		// get scope config from db
		_, scopeConfig, err := scopeHelper.DbHelper().GetScopeAndConfig(connectionId, bpScope.Id)
		if err != nil {
			return nil, err
		}

		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, scopeConfig.Entities)
		if err != nil {
			return nil, err
		}
		stage = append(stage, &plugin.PipelineTask{
			Plugin:   "jira",
			Subtasks: subtasks,
			Options:  options,
		})
		plan[i] = stage
	}

	return plan, nil
}

func makeScopesV200(
	bpScopes []*plugin.BlueprintScopeV200,
	connectionId uint64,
) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0)
	for _, bpScope := range bpScopes {
		// get board and scope config from db
		jiraBoard, scopeConfig, err := scopeHelper.DbHelper().GetScopeAndConfig(connectionId, bpScope.Id)
		if err != nil {
			return nil, err
		}
		// add board to scopes
		if utils.StringsContains(scopeConfig.Entities, plugin.DOMAIN_TYPE_TICKET) {
			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.JiraBoard{}).Generate(jiraBoard.ConnectionId, jiraBoard.BoardId),
				},
				Name: jiraBoard.Name,
				Url:  jiraBoard.Self,
				Type: jiraBoard.Type,
			}
			scopes = append(scopes, domainBoard)
		}
	}
	return scopes, nil
}

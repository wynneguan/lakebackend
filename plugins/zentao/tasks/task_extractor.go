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
	"encoding/json"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ExtractTask

var ExtractTaskMeta = plugin.SubTaskMeta{
	Name:             "extractTask",
	EntryPoint:       ExtractTask,
	EnabledByDefault: true,
	Description:      "extract Zentao task",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractTask(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)

	// this collect only work for project
	if data.Options.ProjectId == 0 {
		return nil
	}

	statusMappings := getTaskStatusMapping(data)
	stdTypeMappings := getStdTypeMappings(data)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ScopeParams(
				data.Options.ConnectionId,
				data.Options.ProjectId,
				data.Options.ProductId,
			),
			Table: RAW_TASK_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			res := &models.ZentaoTaskRes{}
			err := json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}

			// set storyList and FromBugList
			data.StoryList[res.StoryID] = res.Story
			data.FromBugList[res.FromBug] = true

			task := &models.ZentaoTask{
				ConnectionId:       data.Options.ConnectionId,
				ID:                 res.Id,
				Project:            res.Project,
				Parent:             res.Parent,
				Execution:          res.Execution,
				Module:             res.Module,
				Design:             res.Design,
				Story:              res.Story,
				StoryVersion:       res.StoryVersion,
				DesignVersion:      res.DesignVersion,
				FromBug:            res.FromBug,
				Feedback:           res.Feedback,
				FromIssue:          res.FromIssue,
				Name:               res.Name,
				Type:               res.Type,
				Mode:               res.Mode,
				Pri:                res.Pri,
				Estimate:           res.Estimate,
				Consumed:           res.Consumed,
				Left:               res.Left,
				Deadline:           res.Deadline,
				Status:             res.Status,
				SubStatus:          res.SubStatus,
				Color:              res.Color,
				Description:        res.Description,
				Version:            res.Version,
				OpenedById:         getAccountId(res.OpenedBy),
				OpenedByName:       getAccountName(res.OpenedBy),
				OpenedDate:         res.OpenedDate,
				AssignedToId:       getAccountId(res.AssignedTo),
				AssignedToName:     getAccountName(res.AssignedTo),
				AssignedDate:       res.AssignedDate,
				EstStarted:         res.EstStarted,
				RealStarted:        res.RealStarted,
				FinishedId:         getAccountId(res.FinishedBy),
				FinishedDate:       res.FinishedDate,
				FinishedList:       res.FinishedList,
				CanceledId:         getAccountId(res.CanceledBy),
				CanceledDate:       res.CanceledDate,
				ClosedById:         getAccountId(res.ClosedBy),
				ClosedDate:         res.ClosedDate,
				PlanDuration:       res.PlanDuration,
				RealDuration:       res.RealDuration,
				ClosedReason:       res.ClosedReason,
				LastEditedId:       getAccountId(res.LastEditedBy),
				LastEditedDate:     res.LastEditedDate,
				ActivatedDate:      res.ActivatedDate,
				OrderIn:            res.OrderIn,
				Repo:               res.Repo,
				Mr:                 res.Mr,
				Entry:              res.Entry,
				NumOfLine:          res.NumOfLine,
				V1:                 res.V1,
				V2:                 res.V2,
				Deleted:            res.Deleted,
				Vision:             res.Vision,
				StoryID:            res.Story,
				StoryTitle:         res.StoryTitle,
				LatestStoryVersion: 0,
				//Product:            getAccountId(res.Product),
				//Branch:             res.Branch,
				//LatestStoryVersion: res.LatestStoryVersion,
				//StoryStatus:        res.StoryStatus,
				AssignedToRealName: res.AssignedToRealName,
				PriOrder:           res.PriOrder,
				NeedConfirm:        res.NeedConfirm,
				//ProductType:        res.ProductType,
				Progress: res.Progress,
				Url:      row.Url,
			}

			task.StdType = stdTypeMappings[task.Type]
			if task.StdType == "" {
				task.StdType = ticket.TASK
			}

			if len(statusMappings) != 0 {
				task.StdStatus = statusMappings[task.Status]
			} else {
				task.StdStatus = ticket.GetStatus(&ticket.StatusRule{
					Done:    []string{"done", "closed", "cancel"},
					Todo:    []string{"wait"},
					Default: ticket.IN_PROGRESS,
				}, task.Status)
			}

			results := make([]interface{}, 0)
			results = append(results, task)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

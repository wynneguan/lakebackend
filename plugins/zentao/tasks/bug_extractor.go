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

var _ plugin.SubTaskEntryPoint = ExtractBug

var ExtractBugMeta = plugin.SubTaskMeta{
	Name:             "extractBug",
	EntryPoint:       ExtractBug,
	EnabledByDefault: true,
	Description:      "extract Zentao bug",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractBug(taskCtx plugin.SubTaskContext) errors.Error {
	return RangeProductOneByOne(taskCtx, ExtractBugForOneProduct)
}

func ExtractBugForOneProduct(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)

	// this Extract only work for product
	if data.Options.ProductId == 0 {
		return nil
	}

	statusMappings := getBugStatusMapping(data)
	stdTypeMappings := getStdTypeMappings(data)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ScopeParams(
				data.Options.ConnectionId,
				data.Options.ProjectId,
				data.Options.ProductId,
			),
			Table: RAW_BUG_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			res := &models.ZentaoBugRes{}
			err := json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}

			// project scope need filter
			if data.Options.ProjectId != 0 {
				if init, ok := data.FromBugList[int(res.ID)]; !ok || !init {
					return nil, nil
				}
			}

			bug := &models.ZentaoBug{
				ConnectionId:   data.Options.ConnectionId,
				ID:             res.ID,
				Project:        res.Project,
				Product:        res.Product,
				Injection:      res.Injection,
				Identify:       res.Identify,
				Branch:         res.Branch,
				Module:         res.Module,
				Execution:      res.Execution,
				Plan:           res.Plan,
				Story:          res.Story,
				StoryVersion:   res.StoryVersion,
				Task:           res.Task,
				ToTask:         res.ToTask,
				ToStory:        res.ToStory,
				Title:          res.Title,
				Keywords:       res.Keywords,
				Severity:       res.Severity,
				Pri:            res.Pri,
				Type:           res.Type,
				Os:             res.Os,
				Browser:        res.Browser,
				Hardware:       res.Hardware,
				Found:          res.Found,
				Steps:          res.Steps,
				Status:         res.Status,
				SubStatus:      res.SubStatus,
				Color:          res.Color,
				Confirmed:      res.Confirmed,
				ActivatedCount: res.ActivatedCount,
				ActivatedDate:  res.ActivatedDate,
				FeedbackBy:     res.FeedbackBy,
				NotifyEmail:    res.NotifyEmail,
				OpenedById:     getAccountId(res.OpenedBy),
				OpenedByName:   getAccountName(res.OpenedBy),
				OpenedDate:     res.OpenedDate,
				OpenedBuild:    res.OpenedBuild,
				AssignedToId:   getAccountId(res.AssignedTo),
				AssignedToName: getAccountName(res.AssignedTo),
				AssignedDate:   res.AssignedDate,
				Deadline:       res.Deadline,
				ResolvedById:   getAccountId(res.ResolvedBy),
				Resolution:     res.Resolution,
				ResolvedBuild:  res.ResolvedBuild,
				ResolvedDate:   res.ResolvedDate,
				ClosedById:     getAccountId(res.ClosedBy),
				ClosedDate:     res.ClosedDate,
				DuplicateBug:   res.DuplicateBug,
				LinkBug:        res.LinkBug,
				Feedback:       res.Feedback,
				Result:         res.Result,
				Repo:           res.Repo,
				Mr:             res.Mr,
				Entry:          res.Entry,
				NumOfLine:      res.NumOfLine,
				V1:             res.V1,
				V2:             res.V2,
				RepoType:       res.RepoType,
				IssueKey:       res.IssueKey,
				Testtask:       res.Testtask,
				LastEditedById: getAccountId(res.LastEditedBy),
				LastEditedDate: res.LastEditedDate,
				Deleted:        res.Deleted,
				PriOrder:       res.PriOrder,
				SeverityOrder:  res.SeverityOrder,
				Needconfirm:    res.Needconfirm,
				StatusName:     res.StatusName,
				ProductStatus:  res.ProductStatus,
				Url:            row.Url,
			}

			bug.StdType = stdTypeMappings[bug.Type]
			if bug.StdType == "" {
				bug.StdType = ticket.BUG
			}

			if len(statusMappings) != 0 {
				bug.StdStatus = statusMappings[bug.Status]
			} else {
				bug.StdStatus = ticket.GetStatus(&ticket.StatusRule{
					Done:    []string{"resolved"},
					Default: ticket.IN_PROGRESS,
				}, bug.Status)
			}

			results := make([]interface{}, 0)
			results = append(results, bug)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ExtractExecutions

var ExtractExecutionMeta = plugin.SubTaskMeta{
	Name:             "extractExecutions",
	EntryPoint:       ExtractExecutions,
	EnabledByDefault: true,
	Description:      "extract Zentao executions",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractExecutions(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)

	// this Extract only work for project
	if data.Options.ProjectId == 0 {
		return nil
	}

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ScopeParams(
				data.Options.ConnectionId,
				data.Options.ProjectId,
				data.Options.ProductId,
			),
			Table: RAW_EXECUTION_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			res := &models.ZentaoExecutionRes{}
			err := json.Unmarshal(row.Data, res)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}

			// append product to taskdata
			for _, product := range res.Products {
				data.ProductList[product.ID] = product.Name
			}

			execution := &models.ZentaoExecution{
				ConnectionId:   data.Options.ConnectionId,
				Id:             res.ID,
				Project:        res.Project,
				ProjectId:      res.Project,
				Model:          res.Model,
				Type:           res.Type,
				Lifetime:       res.Lifetime,
				Budget:         res.Budget,
				BudgetUnit:     res.BudgetUnit,
				Attribute:      res.Attribute,
				Percent:        res.Percent,
				Milestone:      res.Milestone,
				Output:         res.Output,
				Auth:           res.Auth,
				Parent:         res.Parent,
				Path:           res.Path,
				Grade:          res.Grade,
				Name:           res.Name,
				Code:           res.Code,
				PlanBegin:      res.PlanBegin,
				PlanEnd:        res.PlanEnd,
				RealBegan:      res.RealBegan,
				RealEnd:        res.RealEnd,
				Status:         res.Status,
				SubStatus:      res.SubStatus,
				Pri:            res.Pri,
				Description:    res.Description,
				Version:        res.Version,
				ParentVersion:  res.ParentVersion,
				PlanDuration:   res.PlanDuration,
				RealDuration:   res.RealDuration,
				OpenedById:     getAccountId(res.OpenedBy),
				OpenedDate:     res.OpenedDate,
				OpenedVersion:  res.OpenedVersion,
				LastEditedById: getAccountId(res.LastEditedBy),
				LastEditedDate: res.LastEditedDate,
				ClosedById:     getAccountId(res.ClosedBy),
				ClosedDate:     res.ClosedDate,
				CanceledById:   getAccountId(res.CanceledBy),
				CanceledDate:   res.CanceledDate,
				SuspendedDate:  res.SuspendedDate,
				POId:           getAccountId(res.PO),
				PMId:           getAccountId(res.PM),
				QDId:           getAccountId(res.QD),
				RDId:           getAccountId(res.RD),
				Team:           res.Team,
				Acl:            res.Acl,
				OrderIn:        res.OrderIn,
				Vision:         res.Vision,
				DisplayCards:   res.DisplayCards,
				FluidBoard:     res.FluidBoard,
				Deleted:        res.Deleted,
				TotalHours:     res.TotalHours,
				TotalEstimate:  res.TotalEstimate,
				TotalConsumed:  res.TotalConsumed,
				TotalLeft:      res.TotalLeft,
				Progress:       res.Progress,
				CaseReview:     res.CaseReview,
			}
			results := make([]interface{}, 0)
			results = append(results, execution)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

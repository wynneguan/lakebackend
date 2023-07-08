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

package archived

import (
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type ZentaoExecution struct {
	ConnectionId   uint64              `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id             int64               `json:"id" gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false"`
	Project        int64               `json:"project"`
	Model          string              `json:"model"`
	Type           string              `json:"type"`
	Lifetime       string              `json:"lifetime"`
	Budget         string              `json:"budget"`
	BudgetUnit     string              `json:"budgetUnit"`
	Attribute      string              `json:"attribute"`
	Percent        int                 `json:"percent"`
	Milestone      string              `json:"milestone"`
	Output         string              `json:"output"`
	Auth           string              `json:"auth"`
	Parent         int64               `json:"parent"`
	Path           string              `json:"path"`
	Grade          int                 `json:"grade"`
	Name           string              `json:"name"`
	Code           string              `json:"code"`
	PlanBegin      *helper.Iso8601Time `json:"begin"`
	PlanEnd        *helper.Iso8601Time `json:"end"`
	RealBegan      *helper.Iso8601Time `json:"realBegan"`
	RealEnd        *helper.Iso8601Time `json:"realEnd"`
	Status         string              `json:"status"`
	SubStatus      string              `json:"subStatus"`
	Pri            string              `json:"pri"`
	Description    string              `json:"desc"`
	Version        int                 `json:"version"`
	ParentVersion  int                 `json:"parentVersion"`
	PlanDuration   int                 `json:"planDuration"`
	RealDuration   int                 `json:"realDuration"`
	OpenedById     int64
	OpenedDate     *helper.Iso8601Time `json:"openedDate"`
	OpenedVersion  string              `json:"openedVersion"`
	LastEditedById int64
	LastEditedDate *helper.Iso8601Time `json:"lastEditedDate"`
	ClosedById     int64
	ClosedDate     *helper.Iso8601Time `json:"closedDate"`
	CanceledById   int64
	CanceledDate   *helper.Iso8601Time `json:"canceledDate"`
	SuspendedDate  *helper.Iso8601Time `json:"suspendedDate"`
	POId           int64
	PMId           int64
	QDId           int64
	RDId           int64
	Team           string  `json:"team"`
	Acl            string  `json:"acl"`
	OrderIn        int     `json:"order"`
	Vision         string  `json:"vision"`
	DisplayCards   int     `json:"displayCards"`
	FluidBoard     string  `json:"fluidBoard"`
	Deleted        bool    `json:"deleted"`
	TotalHours     float64 `json:"totalHours"`
	TotalEstimate  float64 `json:"totalEstimate"`
	TotalConsumed  float64 `json:"totalConsumed"`
	TotalLeft      float64 `json:"totalLeft"`
	ProjectId      int64
	Progress       float64 `json:"progress"`
	CaseReview     bool    `json:"caseReview"`
	archived.NoPKModel
}

func (ZentaoExecution) TableName() string {
	return "_tool_zentao_executions"
}

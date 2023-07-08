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

package models

import (
	"fmt"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type ZentaoProject struct {
	common.NoPKModel `json:"-"`
	ConnectionId     uint64              `json:"connectionid" mapstructure:"connectionid" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id               int64               `json:"id" mapstructure:"id" gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false"`
	Project          int64               `json:"project" mapstructure:"project"`
	Model            string              `json:"model" mapstructure:"model"`
	Type             string              `json:"type" mapstructure:"type"`
	ProjectType      string              `json:"projectType" mapstructure:"projectType"`
	Lifetime         string              `json:"lifetime" mapstructure:"lifetime"`
	Budget           string              `json:"budget" mapstructure:"budget"`
	BudgetUnit       string              `json:"budgetUnit" mapstructure:"budgetUnit"`
	Attribute        string              `json:"attribute" mapstructure:"attribute"`
	Percent          int                 `json:"percent" mapstructure:"percent"`
	Milestone        string              `json:"milestone" mapstructure:"milestone"`
	Output           string              `json:"output" mapstructure:"output"`
	Auth             string              `json:"auth" mapstructure:"auth"`
	Parent           int64               `json:"parent" mapstructure:"parent"`
	Path             string              `json:"path" mapstructure:"path"`
	Grade            int                 `json:"grade" mapstructure:"grade"`
	Name             string              `json:"name" mapstructure:"name"`
	Code             string              `json:"code" mapstructure:"code"`
	PlanBegin        *helper.Iso8601Time `json:"begin" mapstructure:"begin"`
	PlanEnd          *helper.Iso8601Time `json:"end" mapstructure:"end"`
	RealBegan        *helper.Iso8601Time `json:"realBegan" mapstructure:"realBegan"`
	RealEnd          *helper.Iso8601Time `json:"realEnd" mapstructure:"realEnd"`
	Days             int                 `json:"days" mapstructure:"days"`
	Status           string              `json:"status" mapstructure:"status"`
	SubStatus        string              `json:"subStatus" mapstructure:"subStatus"`
	Pri              string              `json:"pri" mapstructure:"pri"`
	Description      string              `json:"desc" mapstructure:"desc"`
	Version          int                 `json:"version" mapstructure:"version"`
	ParentVersion    int                 `json:"parentVersion" mapstructure:"parentVersion"`
	PlanDuration     int                 `json:"planDuration" mapstructure:"planDuration"`
	RealDuration     int                 `json:"realDuration" mapstructure:"realDuration"`
	//OpenedBy       string    `json:"openedBy" mapstructure:"openedBy"`
	OpenedDate    *helper.Iso8601Time `json:"openedDate" mapstructure:"openedDate"`
	OpenedVersion string              `json:"openedVersion" mapstructure:"openedVersion"`
	//LastEditedBy   string              `json:"lastEditedBy" mapstructure:"lastEditedBy"`
	LastEditedDate *helper.Iso8601Time `json:"lastEditedDate" mapstructure:"lastEditedDate"`
	ClosedBy       string
	ClosedByRes    interface{}         `json:"closedBy" mapstructure:"closedBy" gorm:"-"`
	ClosedDate     *helper.Iso8601Time `json:"closedDate" mapstructure:"closedDate"`
	CanceledBy     string
	CanceledByRes  interface{}         `json:"canceledBy" mapstructure:"canceledBy" gorm:"-"`
	CanceledDate   *helper.Iso8601Time `json:"canceledDate" mapstructure:"canceledDate"`
	SuspendedDate  *helper.Iso8601Time `json:"suspendedDate" mapstructure:"suspendedDate"`
	PO             string              `json:"po" mapstructure:"po"`
	PM             `json:"pm" mapstructure:"pm"`
	QD             string `json:"qd" mapstructure:"qd"`
	RD             string `json:"rd" mapstructure:"rd"`
	Team           string `json:"team" mapstructure:"team"`
	Acl            string `json:"acl" mapstructure:"acl"`
	Whitelist      `json:"whitelist" mapstructure:"" gorm:"-"`
	OrderIn        int    `json:"order" mapstructure:"order"`
	Vision         string `json:"vision" mapstructure:"vision"`
	DisplayCards   int    `json:"displayCards" mapstructure:"displayCards"`
	FluidBoard     string `json:"fluidBoard" mapstructure:"fluidBoard"`
	Deleted        bool   `json:"deleted" mapstructure:"deleted"`
	Delay          int    `json:"delay" mapstructure:"delay"`
	Hours          `json:"hours" mapstructure:"hours"`
	TeamCount      int    `json:"teamCount" mapstructure:"teamCount"`
	LeftTasks      string `json:"leftTasks" mapstructure:"leftTasks"`
	//TeamMembers   []interface{} `json:"teamMembers" gorm:"-"`
	TotalEstimate float64 `json:"totalEstimate" mapstructure:"totalEstimate"`
	TotalConsumed float64 `json:"totalConsumed" mapstructure:"totalConsumed"`
	TotalLeft     float64 `json:"totalLeft" mapstructure:"totalLeft"`
	Progress      float64 `json:"progress" mapstructure:"progress"`
	TotalReal     int     `json:"totalReal" mapstructure:"totalReal"`
	ScopeConfigId uint64  `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId"`
}
type PM struct {
	PmId       int64  `json:"id" mapstructure:"id"`
	PmAccount  string `json:"account" mapstructure:"account"`
	PmAvatar   string `json:"avatar" mapstructure:"avatar"`
	PmRealname string `json:"realname" mapstructure:"realname"`
}
type Whitelist []struct {
	WhitelistID       int64  `json:"id" mapstructure:"id"`
	WhitelistAccount  string `json:"account" mapstructure:"account"`
	WhitelistAvatar   string `json:"avatar" mapstructure:"avatar"`
	WhitelistRealname string `json:"realname" mapstructure:"realname"`
}
type Hours struct {
	HoursTotalEstimate float64 `json:"totalEstimate" mapstructure:"totalEstimate"`
	HoursTotalConsumed float64 `json:"totalConsumed" mapstructure:"totalConsumed"`
	HoursTotalLeft     float64 `json:"totalLeft" mapstructure:"totalLeft"`
	HoursProgress      float64 `json:"progress" mapstructure:"progress"`
	HoursTotalReal     float64 `json:"totalReal" mapstructure:"totalReal"`
}

func (p *ZentaoProject) ConvertFix() {
	switch cb := p.ClosedByRes.(type) {
	case string:
		p.ClosedBy = cb
	default:
		if cb == nil {
			p.ClosedBy = ""
		} else {
			p.ClosedBy = fmt.Sprintf("%v", cb)
		}
	}
	p.ClosedByRes = p.ClosedBy

	switch cb := p.CanceledByRes.(type) {
	case string:
		p.CanceledBy = cb
	default:
		if cb == nil {
			p.CanceledBy = ""
		} else {
			p.CanceledBy = fmt.Sprintf("%v", cb)
		}
	}
	p.CanceledByRes = p.CanceledBy
}

func (ZentaoProject) TableName() string {
	return "_tool_zentao_projects"
}

func (p ZentaoProject) ScopeId() string {
	return fmt.Sprintf(`projects/%d`, p.Id)
}

func (p ZentaoProject) ScopeName() string {
	return p.Name
}

func (p ZentaoProject) ScopeParams() interface{} {
	return &ZentaoApiParams{
		ConnectionId: p.ConnectionId,
		ZentaoId:     fmt.Sprintf("projects/%d", p.Id),
	}
}

func (p ZentaoProject) ConvertApiScope() plugin.ToolLayerScope {
	if p.ProjectType == "" {
		p.ProjectType = p.Type
		p.Type = "project"
	}
	return p
}

type ZentaoApiParams struct {
	ConnectionId uint64
	ZentaoId     string
}

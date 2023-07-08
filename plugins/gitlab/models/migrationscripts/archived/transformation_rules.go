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
	"gorm.io/datatypes"
)

type GitlabTransformationRule struct {
	archived.Model
	Name                 string `gorm:"type:varchar(255);index:idx_name_gitlab,unique" validate:"required"`
	PrType               string `mapstructure:"prType" json:"prType" gorm:"type:varchar(255)"`
	PrComponent          string `mapstructure:"prComponent" json:"prComponent" gorm:"type:varchar(255)"`
	PrBodyClosePattern   string `mapstructure:"prBodyClosePattern" json:"prBodyClosePattern" gorm:"type:varchar(255)"`
	IssueSeverity        string `mapstructure:"issueSeverity" json:"issueSeverity" gorm:"type:varchar(255)"`
	IssuePriority        string `mapstructure:"issuePriority" json:"issuePriority" gorm:"type:varchar(255)"`
	IssueComponent       string `mapstructure:"issueComponent" json:"issueComponent" gorm:"type:varchar(255)"`
	IssueTypeBug         string `mapstructure:"issueTypeBug" json:"issueTypeBug" gorm:"type:varchar(255)"`
	IssueTypeIncident    string `mapstructure:"issueTypeIncident" json:"issueTypeIncident" gorm:"type:varchar(255)"`
	IssueTypeRequirement string `mapstructure:"issueTypeRequirement" json:"issueTypeRequirement" gorm:"type:varchar(255)"`
	DeploymentPattern    string `mapstructure:"deploymentPattern" json:"deploymentPattern" gorm:"type:varchar(255)"`
	ProductionPattern    string `mapstructure:"productionPattern,omitempty" json:"productionPattern" gorm:"type:varchar(255)"`
	Refdiff              datatypes.JSONMap
}

func (t GitlabTransformationRule) TableName() string {
	return "_tool_gitlab_transformation_rules"
}

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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/github/models/migrationscripts/archived"
)

type githubRepo20221124 struct {
	TransformationRuleId uint64
	CloneUrl             string `json:"cloneUrl" gorm:"type:varchar(255)" mapstructure:"cloneUrl,omitempty"`
}

func (githubRepo20221124) TableName() string {
	return "_tool_github_repos"
}

type addTransformationRule20221124 struct{}

func (script *addTransformationRule20221124) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &githubRepo20221124{}, &archived.GithubTransformationRule{})
}

func (*addTransformationRule20221124) Version() uint64 {
	return 20221214095902
}

func (*addTransformationRule20221124) Name() string {
	return "add table _tool_github_transformation_rules, add transformation_rule_id&clone_url to _tool_github_repos"
}

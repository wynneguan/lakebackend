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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addCicdScopeDropBuildsJobs struct{}

type cicdPipeline20221107 struct {
	CicdScopeId string
}

func (cicdPipeline20221107) TableName() string {
	return "cicd_pipelines"
}

type cicdTask20221107 struct {
	CicdScopeId string
}

func (cicdTask20221107) TableName() string {
	return "cicd_tasks"
}

func (*addCicdScopeDropBuildsJobs) Up(basicRes context.BasicRes) errors.Error {
	err := basicRes.GetDal().DropTables("builds", "jobs")
	if err != nil {
		return err
	}
	return migrationhelper.AutoMigrateTables(basicRes, &archived.CicdScope{}, &cicdTask20221107{}, &cicdPipeline20221107{})
}

func (*addCicdScopeDropBuildsJobs) Version() uint64 {
	return 20221207000001
}

func (*addCicdScopeDropBuildsJobs) Name() string {
	return "add cicd scope and add cicd_scope_id and drop builds&jobs"
}

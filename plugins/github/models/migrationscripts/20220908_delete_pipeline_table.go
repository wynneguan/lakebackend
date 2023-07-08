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
)

type githubPipeline20220908 struct{}

func (githubPipeline20220908) TableName() string {
	return "_tool_github_pipelines"
}

type deleteGithubPipelineTable struct{}

func (u *deleteGithubPipelineTable) Up(basicRes context.BasicRes) errors.Error {
	// create table
	err := basicRes.GetDal().DropTables(&githubPipeline20220908{})
	if err != nil {
		return errors.Default.Wrap(err, "delete table _tool_github_pipelines error")
	}
	return nil

}

func (*deleteGithubPipelineTable) Version() uint64 {
	return 20220908000001
}

func (*deleteGithubPipelineTable) Name() string {
	return "Github delete github_pipelines table"
}

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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"reflect"
)

const RAW_JOB_TABLE = "jenkins_api_jobs"

var ConvertJobsMeta = plugin.SubTaskMeta{
	Name:             "convertJobs",
	EntryPoint:       ConvertJobs,
	EnabledByDefault: true,
	Description:      "Convert tool layer table jenkins_jobs into  domain layer table jobs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ConvertJobs(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JenkinsTaskData)

	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From("_tool_jenkins_jobs"),
		dal.Where("connection_id = ?", data.Options.ConnectionId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	jobIdGen := didgen.NewDomainIdGenerator(&models.JenkinsJob{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.JenkinsJob{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Params: JenkinsApiParams{
				ConnectionId: data.Options.ConnectionId,
				FullName:     data.Options.JobFullName,
			},
			Ctx:   taskCtx,
			Table: RAW_JOB_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			jenkinsJob := inputRow.(*models.JenkinsJob)
			job := &devops.CicdScope{
				DomainEntity: domainlayer.DomainEntity{
					Id: jobIdGen.Generate(jenkinsJob.ConnectionId, jenkinsJob.FullName),
				},
				Name:        jenkinsJob.FullName,
				Description: jenkinsJob.Description,
				Url:         jenkinsJob.Url,
			}
			return []interface{}{
				job,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

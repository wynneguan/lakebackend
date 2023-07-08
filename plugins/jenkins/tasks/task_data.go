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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

type JenkinsApiParams models.JenkinsApiParams
type JenkinsOptions struct {
	ConnectionId  uint64 `json:"connectionId"`
	ScopeId       string
	ScopeConfigId uint64 `json:"scopeConfigId"`
	JobFullName   string `json:"jobFullName"` // "path1/path2/job name"
	JobName       string `json:"jobName"`     // "job name"
	JobPath       string `json:"jobPath"`     // "job/path1/job/path2"
	TimeAfter     string
	Tasks         []string                   `json:"tasks,omitempty"`
	ScopeConfig   *models.JenkinsScopeConfig `mapstructure:"scopeConfig" json:"scopeConfig"`
}

type JenkinsTaskData struct {
	Options       *JenkinsOptions
	ApiClient     *api.ApiAsyncClient
	Connection    *models.JenkinsConnection
	TimeAfter     *time.Time
	RegexEnricher *api.RegexEnricher
}

func DecodeTaskOptions(options map[string]interface{}) (*JenkinsOptions, errors.Error) {
	var op JenkinsOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "could not decode request parameters")
	}
	return &op, nil
}

func ValidateTaskOptions(op *JenkinsOptions) (*JenkinsOptions, errors.Error) {
	if op.ConnectionId == 0 {
		return nil, errors.BadInput.New("connectionId is invalid")
	}
	if op.JobFullName == "" {
		return nil, errors.BadInput.New("JobFullName is required for Jenkins execution")
	}
	if i := strings.LastIndex(op.JobFullName, `/`); i >= 0 {
		op.JobName = op.JobFullName[i+1:]
		op.JobPath = `job/` + strings.Join(strings.Split(op.JobFullName[:i], `/`), `/job/`)
	} else {
		op.JobName = op.JobFullName
		op.JobPath = `view/all`
	}
	if op.ScopeConfig == nil && op.ScopeConfigId == 0 {
		op.ScopeConfig = new(models.JenkinsScopeConfig)
	}
	return op, nil
}

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
	"github.com/apache/incubator-devlake/core/models/common"
	"time"
)

type GitlabJob struct {
	ConnectionId uint64 `gorm:"primaryKey"`

	GitlabId     int     `gorm:"primaryKey"`
	ProjectId    int     `gorm:"index"`
	PipelineId   int     `gorm:"index"`
	Status       string  `gorm:"type:varchar(255)"`
	Stage        string  `gorm:"type:varchar(255)"`
	Name         string  `gorm:"type:varchar(255)"`
	Ref          string  `gorm:"type:varchar(255)"`
	Tag          bool    `gorm:"type:boolean"`
	AllowFailure bool    `json:"allow_failure"`
	Duration     float64 `gorm:"type:float8"`
	WebUrl       string  `gorm:"type:varchar(255)"`

	GitlabCreatedAt *time.Time
	StartedAt       *time.Time
	FinishedAt      *time.Time

	common.NoPKModel
}

func (GitlabJob) TableName() string {
	return "_tool_gitlab_jobs"
}

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
	"time"
)

// Add gitlab job file for GitlabCI
// Add gitlab_updated_at on gitlab pipeline for GitlabCI
type addGitlabCI struct{}

type gitlabPipeline20220729 struct {
	ConnectionId uint64 `gorm:"primaryKey"`

	GitlabId  int    `gorm:"primaryKey"`
	ProjectId int    `gorm:"index"`
	Status    string `gorm:"type:varchar(100)"`
	Ref       string `gorm:"type:varchar(255)"`
	Sha       string `gorm:"type:varchar(255)"`
	WebUrl    string `gorm:"type:varchar(255)"`
	Duration  int

	GitlabCreatedAt *time.Time
	GitlabUpdatedAt *time.Time
	StartedAt       *time.Time
	FinishedAt      *time.Time
	Coverage        string

	archived.NoPKModel
}

func (gitlabPipeline20220729) TableName() string {
	return "_tool_gitlab_pipelines"
}

type gitlabJob20220729 struct {
	ConnectionId uint64 `gorm:"primaryKey"`

	GitlabId     int     `gorm:"primaryKey"`
	ProjectId    int     `gorm:"index"`
	Status       string  `gorm:"type:varchar(255)"`
	Stage        string  `gorm:"type:varchar(255)"`
	Name         string  `gorm:"type:varchar(255)"`
	Ref          string  `gorm:"type:varchar(255)"`
	Tag          bool    `gorm:"type:boolean"`
	AllowFailure bool    `json:"allow_failure"`
	Duration     float64 `gorm:"type:text"`
	WebUrl       string  `gorm:"type:varchar(255)"`

	GitlabCreatedAt *time.Time
	StartedAt       *time.Time
	FinishedAt      *time.Time

	archived.NoPKModel
}

func (gitlabJob20220729) TableName() string {
	return "_tool_gitlab_jobs"
}

func (*addGitlabCI) Up(baseRes context.BasicRes) errors.Error {
	err := migrationhelper.AutoMigrateTables(
		baseRes,
		&gitlabJob20220729{},
		&gitlabPipeline20220729{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (*addGitlabCI) Version() uint64 {
	return 20220729231236
}

func (*addGitlabCI) Name() string {
	return "pipeline and job"
}

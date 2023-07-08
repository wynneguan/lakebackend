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
	"time"

	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type ZentaoChangelog struct {
	archived.NoPKModel `json:"-"`
	ConnectionId       uint64    `json:"connectionId" mapstructure:"connectionId" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id                 int64     `json:"id" mapstructure:"id" gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false"`
	ObjectId           int       `json:"objectId" mapstructure:"objectId" gorm:"index; NOT NULL"`
	Execution          int       `json:"execution" mapstructure:"execution" `
	Actor              string    `json:"actor" mapstructure:"actor" `
	Action             string    `json:"action" mapstructure:"action"`
	Extra              string    `json:"extra" mapstructure:"extra"`
	ObjectType         string    `json:"objectType" mapstructure:"objectType"`
	Project            int       `json:"project" mapstructure:"project"`
	Product            int       `json:"product" mapstructure:"product"`
	Vision             string    `json:"vision" mapstructure:"vision"`
	Comment            string    `json:"comment" mapstructure:"comment"`
	Efforted           string    `json:"efforted" mapstructure:"efforted"`
	Date               time.Time `json:"date" mapstructure:"date"`
	Read               string    `json:"read" mapstructure:"read"`
}

func (ZentaoChangelog) TableName() string {
	return "_tool_zentao_changelog"
}

type ZentaoChangelogDetail struct {
	archived.NoPKModel `json:"-"`
	ConnectionId       uint64 `json:"connectionId" mapstructure:"connectionId" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id                 int64  `json:"id" mapstructure:"id" gorm:"primaryKey;type:BIGINT  NOT NULL;autoIncrement:false"`
	ChangelogId        int64  `json:"changelogId" mapstructure:"changelogId" gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Field              string `json:"field" mapstructure:"field"`
	Old                string `json:"old" mapstructure:"old"`
	New                string `json:"new" mapstructure:"new"`
	Diff               string `json:"diff" mapstructure:"diff"`
}

func (ZentaoChangelogDetail) TableName() string {
	return "_tool_zentao_changelog_detail"
}

type ZentaoChangelogCom struct {
	Changelog       *ZentaoChangelog
	ChangelogDetail *ZentaoChangelogDetail
}

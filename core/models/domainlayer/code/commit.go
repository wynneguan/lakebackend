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

package code

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
)

type Commit struct {
	common.NoPKModel
	Sha            string `json:"sha" gorm:"primaryKey;type:varchar(40);comment:commit hash"`
	Additions      int    `json:"additions" gorm:"comment:Added lines of code"`
	Deletions      int    `json:"deletions" gorm:"comment:Deleted lines of code"`
	DevEq          int    `json:"deveq" gorm:"comment:Merico developer equivalent from analysis engine"`
	Message        string
	AuthorName     string `gorm:"type:varchar(160)"`
	AuthorEmail    string `gorm:"type:varchar(160)"`
	AuthoredDate   time.Time
	AuthorId       string `gorm:"type:varchar(160)"`
	CommitterName  string `gorm:"type:varchar(160)"`
	CommitterEmail string `gorm:"type:varchar(160)"`
	CommittedDate  time.Time
	CommitterId    string `gorm:"index;type:varchar(160)"`
}

func (Commit) TableName() string {
	return "commits"
}

type CommitFile struct {
	domainlayer.DomainEntity
	CommitSha string `gorm:"index;type:varchar(40)"`
	FilePath  string `gorm:"type:text"`
	Additions int
	Deletions int
}

func (CommitFile) TableName() string {
	return "commit_files"
}

type CommitFileComponent struct {
	common.NoPKModel
	CommitFileId  string `gorm:"primaryKey;type:varchar(160)"`
	ComponentName string `gorm:"type:varchar(160)"`
}

func (CommitFileComponent) TableName() string {
	return "commit_file_components"
}

type CommitLineChange struct {
	domainlayer.DomainEntity
	Id          string `gorm:"type:varchar(160);primaryKey"`
	CommitSha   string `gorm:"type:varchar(40);"`
	NewFilePath string `gorm:"type:varchar(160);"`
	LineNoNew   int    `gorm:"type:int"`
	LineNoOld   int    `gorm:"type:int"`
	OldFilePath string `gorm:"type:varchar(160)"`
	HunkNum     int    `gorm:"type:int"`
	ChangedType string `gorm:"type:varchar(160)"`
	PrevCommit  string `gorm:"type:varchar(160)"`
}

func (CommitLineChange) TableName() string {
	return "commit_line_change"
}

type RepoSnapshot struct {
	common.NoPKModel
	RepoId    string `gorm:"primaryKey;type:varchar(160)"`
	CommitSha string `gorm:"primaryKey;type:varchar(40);"`
	FilePath  string `gorm:"primaryKey;type:varchar(160);"`
	LineNo    int    `gorm:"primaryKey;type:int;"`
}

func (RepoSnapshot) TableName() string {
	return "repo_snapshot"
}

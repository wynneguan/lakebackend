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
)

type GithubAccount struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	Id           int    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Login        string `json:"login" gorm:"type:varchar(255)"`
	Name         string `json:"name" gorm:"type:varchar(255)"`
	Company      string `json:"company" gorm:"type:varchar(255)"`
	Email        string `json:"Email" gorm:"type:varchar(255)"`
	AvatarUrl    string `json:"avatar_url" gorm:"type:varchar(255)"`
	Url          string `json:"url" gorm:"type:varchar(255)"`
	HtmlUrl      string `json:"html_url" gorm:"type:varchar(255)"`
	Type         string `json:"type" gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (GithubAccount) TableName() string {
	return "_tool_github_accounts"
}

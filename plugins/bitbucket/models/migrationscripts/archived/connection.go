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
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type BasicAuth struct {
	Username string `mapstructure:"username" validate:"required" json:"username"`
	Password string `mapstructure:"password" validate:"required" json:"password" encrypt:"yes"`
}

func (ba BasicAuth) GetEncodedToken() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", ba.Username, ba.Password)))
}

type BaseConnection struct {
	Name string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	archived.Model
}

type RestConnection struct {
	BaseConnection   `mapstructure:",squash"`
	Endpoint         string `mapstructure:"endpoint" validate:"required" json:"endpoint"`
	Proxy            string `mapstructure:"proxy" json:"proxy"`
	RateLimitPerHour int    `comment:"api request rate limit per hour" json:"rateLimit"`
}

type TestConnectionRequest struct {
	Endpoint  string `json:"endpoint"`
	Proxy     string `json:"proxy"`
	BasicAuth `mapstructure:",squash"`
}

type ApiUserResponse struct {
	Username      string `json:"username"`
	DisplayName   string `json:"display_name"`
	AccountId     int    `json:"account_id"`
	Uuid          string `json:"uuid"`
	AccountStatus string `json:"account_status"`
}

type BitbucketConnection struct {
	RestConnection `mapstructure:",squash"`
	BasicAuth      `mapstructure:",squash"`
}

func (BitbucketConnection) TableName() string {
	return "_tool_bitbucket_connections"
}

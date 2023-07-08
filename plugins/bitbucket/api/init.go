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

package api

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/go-playground/validator/v10"
)

var vld *validator.Validate
var connectionHelper *api.ConnectionApiHelper
var scopeHelper *api.ScopeApiHelper[models.BitbucketConnection, models.BitbucketRepo, models.BitbucketScopeConfig]
var remoteHelper *api.RemoteApiHelper[models.BitbucketConnection, models.BitbucketRepo, models.BitbucketApiRepo, models.GroupResponse]
var scHelper *api.ScopeConfigHelper[models.BitbucketScopeConfig]
var basicRes context.BasicRes

func Init(br context.BasicRes, p plugin.PluginMeta) {

	basicRes = br
	vld = validator.New()
	connectionHelper = api.NewConnectionHelper(
		basicRes,
		vld,
		p.Name(),
	)
	params := &api.ReflectionParameters{
		ScopeIdFieldName:  "BitbucketId",
		ScopeIdColumnName: "bitbucket_id",
		RawScopeParamName: "FullName",
	}
	scopeHelper = api.NewScopeHelper[models.BitbucketConnection, models.BitbucketRepo, models.BitbucketScopeConfig](
		basicRes,
		vld,
		connectionHelper,
		api.NewScopeDatabaseHelperImpl[models.BitbucketConnection, models.BitbucketRepo, models.BitbucketScopeConfig](
			basicRes, connectionHelper, params),
		params,
		nil,
	)
	remoteHelper = api.NewRemoteHelper[models.BitbucketConnection, models.BitbucketRepo, models.BitbucketApiRepo, models.GroupResponse](
		basicRes,
		vld,
		connectionHelper,
	)
	scHelper = api.NewScopeConfigHelper[models.BitbucketScopeConfig](
		basicRes,
		vld,
		p.Name(),
	)
}

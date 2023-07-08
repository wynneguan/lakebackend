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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

// CreateScopeConfig create scope config for Bamboo
// @Summary create scope config for Bamboo
// @Description create scope config for Bamboo
// @Tags plugins/bamboo
// @Accept application/json
// @Param connectionId path int true "connectionId"
// @Param scopeConfig body models.BambooScopeConfig true "scope config"
// @Success 200  {object} models.BambooScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scope-configs [POST]
func CreateScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeConfigHelper.Create(input)
}

// UpdateScopeConfig update scope config for Bamboo
// @Summary update scope config for Bamboo
// @Description update scope config for Bamboo
// @Tags plugins/bamboo
// @Accept application/json
// @Param id path int true "id"
// @Param scopeConfig body models.BambooScopeConfig true "scope config"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} models.BambooScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scope-configs/{id} [PATCH]
func UpdateScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeConfigHelper.Update(input)
}

// GetScopeConfig return onescope config
// @Summary return onescope config
// @Description return onescope config
// @Tags plugins/bamboo
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} models.BambooScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scope-configs/{id} [GET]
func GetScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeConfigHelper.Get(input)
}

// GetScopeConfigList return all scope configs
// @Summary return all scope configs
// @Description return all scope configs
// @Tags plugins/bamboo
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Param connectionId path int true "connectionId"
// @Success 200  {object} []models.BambooScopeConfig
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scope-configs [GET]
func GetScopeConfigList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeConfigHelper.List(input)
}

// DeleteScopeConfig delete a scope config
// @Summary delete a scope config
// @Description delete a scope config
// @Tags plugins/bamboo
// @Param id path int true "id"
// @Param connectionId path int true "connectionId"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scope-configs/{id} [DELETE]
func DeleteScopeConfig(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeConfigHelper.Delete(input)
}

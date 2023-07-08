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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

type ProductScopeRes struct {
	models.ZentaoProduct
	api.ScopeResDoc[models.ZentaoScopeConfig]
}

type ProductScopeReq api.ScopeReq[models.ZentaoProduct]

type ProjectScopeRes struct {
	models.ZentaoProject
	api.ScopeResDoc[models.ZentaoScopeConfig]
}

type ProjectScopeReq api.ScopeReq[models.ZentaoProject]

// PutProductScope create or update zentao products
// @Summary create or update zentao products
// @Description Create or update zentao products
// @Tags plugins/zentao
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body ProductScopeReq true "json"
// @Success 200  {object} []models.ZentaoProduct
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/product/scopes [PUT]
func PutProductScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return productScopeHelper.Put(input)
}

// PutProjectScope create or update zentao projects
// @Summary create or update zentao projects
// @Description Create or update zentao projects
// @Tags plugins/zentao
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body ProjectScopeReq true "json"
// @Success 200  {object} []models.ZentaoProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/project/scopes [PUT]
func PutProjectScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return projectScopeHelper.Put(input)
}

// PutScope create or update zentao projects
// @Summary create or update zentao projects
// @Description Create or update zentao projects
// @Tags plugins/zentao
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scope body ProjectScopeReq true "json"
// @Success 200  {object} []models.ZentaoProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return PutProjectScope(input)
}

// UpdateProductScope patch to zentao product
// @Summary patch to zentao product
// @Description patch to zentao product
// @Tags plugins/zentao
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param scope body models.ZentaoProduct true "json"
// @Success 200  {object} models.ZentaoProduct
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/product/{scopeId} [PATCH]
func UpdateProductScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return productScopeHelper.Update(input)
}

// UpdateProjectScope patch to zentao project
// @Summary patch to zentao project
// @Description patch to zentao project
// @Tags plugins/zentao
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param scope body models.ZentaoProject true "json"
// @Success 200  {object} models.ZentaoProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/project/{scopeId} [PATCH]
func UpdateProjectScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return projectScopeHelper.Update(input)
}

// UpdateScope patch to zentao project
// @Summary patch to zentao project
// @Description patch to zentao project
// @Tags plugins/zentao
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param scope body models.ZentaoProject true "json"
// @Success 200  {object} models.ZentaoProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return UpdateProjectScope(input)
}

// GetProductScopeList get one product
// @Summary get one product
// @Description get one product
// @Tags plugins/zentao
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Success 200  {object} []ProductScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/scopes/product/ [GET]
func GetProductScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return productScopeHelper.GetScopeList(input)
}

// GetProjectScopeList get Gitlab projects
// @Summary get Gitlab projects
// @Description get Gitlab projects
// @Tags plugins/gitlab
// @Param connectionId path int false "connection ID"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []ProjectScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/scopes/project/ [GET]
func GetProjectScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return projectScopeHelper.GetScopeList(input)
}

// GetScopeList get Gitlab projects
// @Summary get Gitlab projects
// @Description get Gitlab projects
// @Tags plugins/gitlab
// @Param connectionId path int false "connection ID"
// @Param blueprints query bool false "also return blueprints using these scopes as part of the payload"
// @Success 200  {object} []ProjectScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return GetProjectScopeList(input)
}

// GetProductScope get one product
// @Summary get one product
// @Description get one product
// @Tags plugins/zentao
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Success 200  {object} ProductScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/product/{scopeId} [GET]
func GetProductScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return productScopeHelper.GetScope(input)
}

// GetProjectScope get one project
// @Summary get one project
// @Description get one project
// @Tags plugins/zentao
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Success 200  {object} ProjectScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/project/{scopeId} [GET]
func GetProjectScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return projectScopeHelper.GetScope(input)
}

// GetScope get one project
// @Summary get one project
// @Description get one project
// @Tags plugins/zentao
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Success 200  {object} ProjectScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return GetProjectScope(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description delete data associated with plugin scope
// @Tags plugins/zentao
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} api.ScopeRefDoc "References exist to this scope"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/product/{scopeId} [DELETE]
func DeleteProductScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return productScopeHelper.Delete(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description delete data associated with plugin scope
// @Tags plugins/zentao
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 409  {object} api.ScopeRefDoc "References exist to this scope"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/project/{scopeId} [DELETE]
func DeleteProjectScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return projectScopeHelper.Delete(input)
}

// DeleteScope delete plugin data associated with the scope and optionally the scope itself
// @Summary delete plugin data associated with the scope and optionally the scope itself
// @Description delete data associated with plugin scope
// @Tags plugins/zentao
// @Param connectionId path int true "connection ID"
// @Param scopeId path int true "scope ID"
// @Param delete_data_only query bool false "Only delete the scope data, not the scope itself"
// @Success 200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/scopes/{scopeId} [DELETE]
func DeleteScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return DeleteProjectScope(input)
}

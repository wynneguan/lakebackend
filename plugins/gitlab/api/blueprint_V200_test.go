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
	"strconv"
	"testing"
	"time"

	mockdal "github.com/apache/incubator-devlake/mocks/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/unithelper"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/gitlab/tasks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMakeDataSourcePipelinePlanV200(t *testing.T) {
	const testConnectionID uint64 = 1
	const testScopeConfigId uint64 = 2
	const testID int = 37
	const testGitlabEndPoint string = "https://gitlab.com/api/v4/"
	const testHttpUrlToRepo string = "https://this_is_cloneUrl"
	const testToken string = "nddtf"
	const testName string = "gitlab-test"
	const pathWithNamespace string = "nddtf/gitlab-test"
	const testScopeConfigName string = "gitlab scope config"
	const testProxy string = ""

	syncPolicy := &plugin.BlueprintSyncPolicy{}
	bpScopes := []*plugin.BlueprintScopeV200{
		{
			Id:   strconv.Itoa(testID),
			Name: testName,
		},
	}

	var testGitlabProject = &models.GitlabProject{
		ConnectionId:      testConnectionID,
		GitlabId:          testID,
		Name:              testName,
		PathWithNamespace: pathWithNamespace,
		ScopeConfigId:     testScopeConfigId,
		CreatedDate:       &time.Time{},
		HttpUrlToRepo:     testHttpUrlToRepo,
	}

	var testScopeConfig = &models.GitlabScopeConfig{
		ScopeConfig: common.ScopeConfig{
			Model: common.Model{
				ID: testScopeConfigId,
			},
			Entities: []string{plugin.DOMAIN_TYPE_CODE, plugin.DOMAIN_TYPE_TICKET, plugin.DOMAIN_TYPE_CICD},
		},
		Name:   testScopeConfigName,
		PrType: "hey,man,wasup",
		Refdiff: map[string]interface{}{
			"tagsPattern": "pattern",
			"tagsLimit":   10,
			"tagsOrder":   "reverse semver",
		},
	}

	var testGitlabConnection = &models.GitlabConnection{
		BaseConnection: helper.BaseConnection{
			Name: testName,
			Model: common.Model{
				ID: testConnectionID,
			},
		},
		GitlabConn: models.GitlabConn{
			RestConnection: helper.RestConnection{
				Endpoint:         testGitlabEndPoint,
				Proxy:            testProxy,
				RateLimitPerHour: 0,
			},
			AccessToken: helper.AccessToken{
				Token: testToken,
			},
		},
	}

	var expectRepoId = "gitlab:GitlabProject:1:37"

	var testSubTaskMeta = []plugin.SubTaskMeta{
		tasks.ConvertProjectMeta,
		tasks.CollectApiIssuesMeta,
		tasks.ExtractApiIssuesMeta,
		tasks.ConvertIssuesMeta,
		tasks.ConvertIssueLabelsMeta,
		tasks.CollectApiJobsMeta,
		tasks.ExtractApiJobsMeta,
		tasks.CollectApiPipelinesMeta,
		tasks.ExtractApiPipelinesMeta,
	}

	var expectPlans = plugin.PipelinePlan{
		{
			{
				Plugin: "gitlab",
				Subtasks: []string{
					tasks.ConvertProjectMeta.Name,
					tasks.CollectApiIssuesMeta.Name,
					tasks.ExtractApiIssuesMeta.Name,
					tasks.ConvertIssuesMeta.Name,
					tasks.ConvertIssueLabelsMeta.Name,
					tasks.CollectApiJobsMeta.Name,
					tasks.ExtractApiJobsMeta.Name,
					tasks.CollectApiPipelinesMeta.Name,
					tasks.ExtractApiPipelinesMeta.Name,
				},
				Options: map[string]interface{}{
					"connectionId":  uint64(1),
					"projectId":     testID,
					"scopeConfigId": testScopeConfigId,
				},
			},
			{
				Plugin: "gitextractor",
				Options: map[string]interface{}{
					"proxy":  "",
					"repoId": expectRepoId,
					"name":   testName,
					"url":    "https://git:nddtf@this_is_cloneUrl",
				},
			},
		},
		{
			{
				Plugin: "refdiff",
				Options: map[string]interface{}{
					"tagsLimit":   10,
					"tagsOrder":   "reverse semver",
					"tagsPattern": "pattern",
				},
			},
		},
	}

	expectRepo := code.NewRepo(expectRepoId, pathWithNamespace)
	expectRepo.ForkedFrom = testGitlabProject.ForkedFromProjectWebUrl

	expectCicdScope := devops.NewCicdScope(expectRepoId, pathWithNamespace)
	expectCicdScope.Description = ""
	expectCicdScope.Url = ""

	expectBoard := ticket.NewBoard(expectRepoId, pathWithNamespace)
	expectBoard.Description = ""
	expectBoard.Url = ""
	expectBoard.Type = ""

	var err errors.Error

	// register gitlab plugin for NewDomainIdGenerator
	mockMeta := mockplugin.NewPluginMeta(t)
	mockMeta.On("RootPkgPath").Return("github.com/apache/incubator-devlake/plugins/gitlab")
	mockMeta.On("Name").Return("dummy").Maybe()
	err = plugin.RegisterPlugin("gitlab", mockMeta)
	assert.Equal(t, err, nil)

	// Refresh Global Variables and set the sql mock
	mockRes := unithelper.DummyBasicRes(func(mockDal *mockdal.Dal) {
		mockDal.On("First", mock.AnythingOfType("*models.GitlabConnection"), mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.GitlabConnection)
			*dst = *testGitlabConnection
		}).Return(nil)

		mockDal.On("First", mock.AnythingOfType("*models.GitlabProject"), mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.GitlabProject)
			*dst = *testGitlabProject
		}).Return(nil)

		mockDal.On("First", mock.AnythingOfType("*models.GitlabScopeConfig"), mock.Anything).Run(func(args mock.Arguments) {
			dst := args.Get(0).(*models.GitlabScopeConfig)
			*dst = *testScopeConfig
		}).Return(nil)
	})
	Init(mockRes, mockMeta)

	plans, scopes, err := MakePipelinePlanV200(testSubTaskMeta, testConnectionID, bpScopes, syncPolicy)
	assert.Equal(t, err, nil)

	assert.Equal(t, expectPlans, plans)

	// ignore CreatedDate UpdatedDate  CreatedAt UpdatedAt checking
	expectRepo.CreatedDate = scopes[0].(*code.Repo).CreatedDate
	expectRepo.UpdatedDate = scopes[0].(*code.Repo).UpdatedDate
	expectRepo.CreatedAt = scopes[0].(*code.Repo).CreatedAt
	expectRepo.UpdatedAt = scopes[0].(*code.Repo).UpdatedAt

	expectCicdScope.CreatedDate = scopes[1].(*devops.CicdScope).CreatedDate
	expectCicdScope.UpdatedDate = scopes[1].(*devops.CicdScope).UpdatedDate
	expectCicdScope.CreatedAt = scopes[1].(*devops.CicdScope).CreatedAt
	expectCicdScope.UpdatedAt = scopes[1].(*devops.CicdScope).UpdatedAt

	expectBoard.CreatedDate = scopes[2].(*ticket.Board).CreatedDate
	expectBoard.CreatedAt = scopes[2].(*ticket.Board).CreatedAt
	expectBoard.UpdatedAt = scopes[2].(*ticket.Board).UpdatedAt

	var expectScopes = []plugin.Scope{
		expectRepo,
		expectCicdScope,
		expectBoard,
	}

	assert.Equal(t, expectScopes, scopes)
}

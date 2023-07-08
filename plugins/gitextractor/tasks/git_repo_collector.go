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

package tasks

import (
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/gitextractor/parser"
)

type GitExtractorOptions struct {
	RepoId     string `json:"repoId"`
	Name       string `jsno:"name"`
	Url        string `json:"url"`
	User       string `json:"user"`
	Password   string `json:"password"`
	PrivateKey string `json:"privateKey"`
	Passphrase string `json:"passphrase"`
	Proxy      string `json:"proxy"`
}

func (o GitExtractorOptions) Valid() errors.Error {
	if o.RepoId == "" {
		return errors.BadInput.New("empty repoId")
	}
	if o.Url == "" {
		return errors.BadInput.New("empty url")
	}
	url := strings.TrimPrefix(o.Url, "ssh://")
	if !(strings.HasPrefix(o.Url, "http") || strings.HasPrefix(url, "git@") || strings.HasPrefix(o.Url, "/")) {
		return errors.BadInput.New("wrong url")
	}
	return nil
}

func CollectGitCommits(subTaskCtx plugin.SubTaskContext) errors.Error {
	repo := getGitRepo(subTaskCtx)
	if count, err := repo.CountCommits(subTaskCtx.GetContext()); err != nil {
		subTaskCtx.GetLogger().Error(err, "unable to get commit count")
		subTaskCtx.SetProgress(0, -1)
		return err
	} else {
		subTaskCtx.SetProgress(0, count)
	}
	return repo.CollectCommits(subTaskCtx)
}

func CollectGitBranches(subTaskCtx plugin.SubTaskContext) errors.Error {
	repo := getGitRepo(subTaskCtx)
	if count, err := repo.CountBranches(subTaskCtx.GetContext()); err != nil {
		subTaskCtx.GetLogger().Error(err, "unable to get branch count")
		subTaskCtx.SetProgress(0, -1)
		return err
	} else {
		subTaskCtx.SetProgress(0, count)
	}
	return repo.CollectBranches(subTaskCtx)
}

func CollectGitTags(subTaskCtx plugin.SubTaskContext) errors.Error {
	repo := getGitRepo(subTaskCtx)
	if count, err := repo.CountTags(); err != nil {
		subTaskCtx.GetLogger().Error(err, "unable to get tag count")
		subTaskCtx.SetProgress(0, -1)
		return err
	} else {
		subTaskCtx.SetProgress(0, count)
	}
	return repo.CollectTags(subTaskCtx)
}

func CollectGitDiffLines(subTaskCtx plugin.SubTaskContext) errors.Error {
	repo := getGitRepo(subTaskCtx)
	if count, err := repo.CountTags(); err != nil {
		subTaskCtx.GetLogger().Error(err, "unable to get line content")
		subTaskCtx.SetProgress(0, -1)
		return err
	} else {
		subTaskCtx.SetProgress(0, count)
	}
	return repo.CollectDiffLine(subTaskCtx)
}

func getGitRepo(subTaskCtx plugin.SubTaskContext) *parser.GitRepo {
	repo, ok := subTaskCtx.GetData().(*parser.GitRepo)
	if !ok {
		panic("git repo reference not found on context")
	}
	return repo
}

var CollectGitCommitMeta = plugin.SubTaskMeta{
	Name:             "collectGitCommits",
	EntryPoint:       CollectGitCommits,
	EnabledByDefault: true,
	Description:      "collect git commits into Domain Layer Tables",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE, plugin.DOMAIN_TYPE_CROSS},
}

var CollectGitBranchMeta = plugin.SubTaskMeta{
	Name:             "collectGitBranches",
	EntryPoint:       CollectGitBranches,
	EnabledByDefault: true,
	Description:      "collect git branch into Domain Layer Tables",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

var CollectGitTagMeta = plugin.SubTaskMeta{
	Name:             "collectGitTags",
	EntryPoint:       CollectGitTags,
	EnabledByDefault: true,
	Description:      "collect git tag into Domain Layer Tables",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

var CollectGitDiffLineMeta = plugin.SubTaskMeta{
	Name:             "collectDiffLine",
	EntryPoint:       CollectGitDiffLines,
	EnabledByDefault: false,
	Description:      "collect git commit diff line into Domain Layer Tables",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

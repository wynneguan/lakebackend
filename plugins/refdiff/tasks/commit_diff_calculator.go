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
	"fmt"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/refdiff/models"
	"github.com/apache/incubator-devlake/plugins/refdiff/utils"
)

func CalculateCommitsDiff(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*RefdiffTaskData)
	repoId := data.Options.RepoId
	db := taskCtx.GetDal()
	ctx := taskCtx.GetContext()
	logger := taskCtx.GetLogger()

	if data.Options.ProjectName != "" {
		return nil
	}

	// get all data from finish_commits_diffs
	commitPairsSrc := data.Options.AllPairs
	var commitPairs RefCommitPairs
	refCommit := &code.RefCommit{}
	for _, pair := range commitPairsSrc {
		newRefId := fmt.Sprintf("%s:%s", repoId, pair[2])
		oldRefId := fmt.Sprintf("%s:%s", repoId, pair[3])

		count, err := db.Count(
			dal.Select("*"),
			dal.From("_tool_refdiff_finished_commits_diffs"),
			dal.Where("new_commit_sha = ? and old_commit_sha = ?", pair[0], pair[1]))
		if err != nil {
			return err
		}
		if count == 0 {
			commitPairs = append(commitPairs, pair)
		}
		if pair[2] != newRefId || pair[3] != oldRefId {
			refCommit.NewCommitSha = pair[0]
			refCommit.OldCommitSha = pair[1]
			refCommit.NewRefId = newRefId
			refCommit.OldRefId = oldRefId
		}
	}

	if len(commitPairs) == 0 {
		logger.Info("commit pair has been produced.")
		return nil
	}

	commitNodeGraph := utils.NewCommitNodeGraph()
	// mysql limit
	insertCountLimitOfCommitsDiff := int(65535 / reflect.ValueOf(code.CommitsDiff{}).NumField())

	// load commits from db
	commitParent := &code.CommitParent{}
	cursor, err := db.Cursor(
		dal.Select("cp.*"),
		dal.Join("LEFT JOIN repo_commits rc ON (rc.commit_sha = cp.commit_sha)"),
		dal.From("commit_parents cp"),
		dal.Where("rc.repo_id = ?", repoId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}
		err = db.Fetch(cursor, commitParent)
		if err != nil {
			return errors.Default.Wrap(err, "failed to read commit from database")
		}
		commitNodeGraph.AddParent(commitParent.CommitSha, commitParent.ParentCommitSha)
	}

	logger.Info("Create a commit node graph with node count[%d]", commitNodeGraph.Size())

	// calculate diffs for commits pairs and store them into database
	commitsDiff := &code.CommitsDiff{}
	finishedCommitDiff := &models.FinishedCommitsDiff{}
	lenCommitPairs := len(commitPairs)
	taskCtx.SetProgress(0, lenCommitPairs)

	for _, pair := range commitPairs {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}
		// ref might advance, keep commit sha for debugging
		commitsDiff.NewCommitSha = pair[0]
		commitsDiff.OldCommitSha = pair[1]

		finishedCommitDiff.NewCommitSha = pair[0]
		finishedCommitDiff.OldCommitSha = pair[1]

		if commitsDiff.NewCommitSha == commitsDiff.OldCommitSha {
			// different refs might point to a same commit, it is ok
			logger.Info(
				"skipping ref pair due to they are the same %s",
				commitsDiff.NewCommitSha,
			)
			continue
		}

		lostSha, oldCount, newCount := commitNodeGraph.CalculateLostSha(pair[1], pair[0])

		commitsDiffs := []code.CommitsDiff{}
		refCommits := []code.RefCommit{}
		finishedCommitDiffs := []models.FinishedCommitsDiff{}

		commitsDiff.SortingIndex = 1
		for _, sha := range lostSha {
			commitsDiff.CommitSha = sha
			commitsDiffs = append(commitsDiffs, *commitsDiff)

			// sql limit placeholders count only 65535
			if commitsDiff.SortingIndex%insertCountLimitOfCommitsDiff == 0 {
				logger.Info("commitsDiffs count in limited[%d] index[%d]--exec and clean", len(commitsDiffs), commitsDiff.SortingIndex)
				err = db.CreateIfNotExist(commitsDiffs)
				if err != nil {
					return err
				}
				commitsDiffs = []code.CommitsDiff{}
			}

			commitsDiff.SortingIndex++
		}

		if len(commitsDiffs) > 0 {
			logger.Info("insert data count [%d]", len(commitsDiffs))
			err = db.CreateIfNotExist(commitsDiffs)
			if err != nil {
				return err
			}
		}

		refCommits = append(refCommits, *refCommit)
		if len(refCommits) > 0 {
			err = db.CreateIfNotExist(refCommits)
			if err != nil {
				return err
			}
		}

		finishedCommitDiffs = append(finishedCommitDiffs, *finishedCommitDiff)
		if len(finishedCommitDiffs) > 0 {
			err = db.CreateIfNotExist(finishedCommitDiffs)
			if err != nil {
				return err
			}
		}

		logger.Info(
			"total %d commits of difference found between [new][%s] and [old][%s(total:%d)]",
			newCount,
			commitsDiff.NewCommitSha,
			commitsDiff.OldCommitSha,
			oldCount,
		)
		taskCtx.IncProgress(1)
	}
	return nil
}

var CalculateCommitsDiffMeta = plugin.SubTaskMeta{
	Name:             "calculateCommitsDiff",
	EntryPoint:       CalculateCommitsDiff,
	EnabledByDefault: true,
	Description:      "Calculate diff commits between refs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

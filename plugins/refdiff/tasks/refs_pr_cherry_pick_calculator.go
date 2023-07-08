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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/plugin"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type cherryPick struct {
	RepoName             string `gorm:"type:varchar(255)"`
	ParentPrKey          int
	CherrypickBaseBranch string `gorm:"type:varchar(255)"`
	CherrypickPrKey      int
	ParentPrUrl          string `gorm:"type:varchar(255)"`
	ParentPrId           string `gorm:"type:varchar(255)"`
	CreatedDate          time.Time
}

func CalculatePrCherryPick(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*RefdiffTaskData)
	repoId := data.Options.RepoId
	ctx := taskCtx.GetContext()
	db := taskCtx.GetDal()

	if data.Options.ProjectName != "" {
		return nil
	}
	var prTitleRegex *regexp.Regexp
	var err error
	prTitlePattern := taskCtx.GetConfig("GITHUB_PR_TITLE_PATTERN")
	if len(prTitlePattern) > 0 {
		prTitleRegex, err = regexp.Compile(prTitlePattern)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile prTitlePattern failed")
		}
	}

	cursor, err := db.Cursor(
		dal.From(&code.PullRequest{}),
		dal.Join("left join repos on pull_requests.base_repo_id = repos.id"),
		dal.Where("repos.id = ?", repoId),
	)
	if err != nil {
		return errors.Convert(err)
	}

	defer cursor.Close()

	pr := &code.PullRequest{}
	var parentPrKeyInt int
	taskCtx.SetProgress(0, -1)

	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}

		err = db.Fetch(cursor, pr)
		if err != nil {
			return errors.Convert(err)
		}

		parentPrKey := ""
		if prTitleRegex != nil {
			groups := prTitleRegex.FindStringSubmatch(pr.Title)
			if len(groups) > 1 {
				parentPrKey = groups[1]
			}
		}

		if parentPrKeyInt, err = strconv.Atoi(parentPrKey); err != nil {
			continue
		}

		var parentPrId string
		err = db.Pluck("id", &parentPrId,
			dal.Where("pull_request_key = ? and base_repo_id = ?", parentPrKeyInt, repoId),
			dal.From("pull_requests"),
		)
		if err != nil {
			return errors.Convert(err)
		}
		if len(parentPrId) == 0 {
			continue
		}
		pr.ParentPrId = parentPrId

		err = db.Update(pr)
		if err != nil {
			return errors.Convert(err)
		}
		taskCtx.IncProgress(1)
	}

	cursor2, err := db.Cursor(
		dal.Select(`
			pr2.pull_request_key                 AS parent_pr_key,
			pr1.parent_pr_id                     AS parent_pr_id,
			pr1.base_ref                         AS cherrypick_base_branch,
			pr1.pull_request_key                 AS cherrypick_pr_key,
			repos.NAME                           AS repo_name,
			Concat(repos.url, '/pull/', pr2.pull_request_key) AS parent_pr_url,
			pr2.created_date
		`),
		dal.From(`pull_requests pr1`),
		dal.Join(`LEFT JOIN pull_requests pr2 ON pr1.parent_pr_id = pr2.id`),
		dal.Join(`LEFT JOIN repos ON pr2.base_repo_id = repos.id`),

		dal.Where("pr1.parent_pr_id != ''"),
		dal.Orderby(`
			pr1.parent_pr_id,
			pr2.created_date,
			pr1.base_ref ASC
		`),
	)
	if err != nil {
		return errors.Convert(err)
	}
	defer cursor2.Close()

	var refsPrCherryPick *code.RefsPrCherrypick
	var lastParentPrId string
	var lastCreatedDate time.Time
	var cherrypickBaseBranches []string
	var cherrypickPrKeys []string
	for cursor2.Next() {
		var item cherryPick
		err = db.Fetch(cursor2, &item)
		if err != nil {
			return errors.Convert(err)
		}
		if item.ParentPrId == lastParentPrId && item.CreatedDate == lastCreatedDate {
			cherrypickBaseBranches = append(cherrypickBaseBranches, item.CherrypickBaseBranch)
			cherrypickPrKeys = append(cherrypickPrKeys, strconv.Itoa(item.CherrypickPrKey))
		} else {
			if refsPrCherryPick != nil {
				refsPrCherryPick.CherrypickBaseBranches = strings.Join(cherrypickBaseBranches, ",")
				refsPrCherryPick.CherrypickPrKeys = strings.Join(cherrypickPrKeys, ",")
				err = db.CreateOrUpdate(refsPrCherryPick)
				if err != nil {
					return errors.Convert(err)
				}
			}
			lastParentPrId = item.ParentPrId
			lastCreatedDate = item.CreatedDate
			cherrypickBaseBranches = []string{item.CherrypickBaseBranch}
			cherrypickPrKeys = []string{strconv.Itoa(item.CherrypickPrKey)}
			refsPrCherryPick = &code.RefsPrCherrypick{
				RepoName:    item.RepoName,
				ParentPrKey: item.ParentPrKey,
				ParentPrUrl: item.ParentPrUrl,
				ParentPrId:  item.ParentPrId,
			}
		}
	}

	if refsPrCherryPick != nil {
		err = db.CreateOrUpdate(refsPrCherryPick)
		if err != nil {
			return errors.Convert(err)
		}
	}

	return nil
}

var CalculatePrCherryPickMeta = plugin.SubTaskMeta{
	Name:             "calculatePrCherryPick",
	EntryPoint:       CalculatePrCherryPick,
	EnabledByDefault: true,
	Description:      "Calculate pr cherry pick",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

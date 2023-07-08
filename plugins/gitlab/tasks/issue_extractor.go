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
	"encoding/json"
	"regexp"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiIssuesMeta)
}

var ExtractApiIssuesMeta = plugin.SubTaskMeta{
	Name:             "extractApiIssues",
	EntryPoint:       ExtractApiIssues,
	EnabledByDefault: true,
	Description:      "Extract raw Issues data into tool layer table gitlab_issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	Dependencies:     []*plugin.SubTaskMeta{&CollectApiIssuesMeta},
}

type IssuesResponse struct {
	ProjectId int `json:"project_id"`
	Milestone struct {
		Due_date    string
		Project_id  int
		State       string
		Description string
		Iid         int
		Id          int
		Title       string
		CreatedAt   api.Iso8601Time
		UpdatedAt   api.Iso8601Time
	}
	Author *struct {
		State     string
		WebUrl    string `json:"web_url"`
		AvatarUrl string `json:"avatar_url"`
		Username  string
		Id        int    `json:"id"`
		Name      string `json:"name"`
	}
	Description string
	State       string
	Iid         int
	Assignees   []struct {
		AvatarUrl string `json:"avatar_url"`
		WebUrl    string `json:"web_url"`
		State     string
		Username  string
		Id        int
		Name      string
	}
	Assignee *struct {
		AvatarUrl string
		WebUrl    string
		State     string
		Username  string
		Id        int
		Name      string
	}
	Type               string
	Labels             []string `json:"labels"`
	UpVotes            int
	DownVotes          int
	MergeRequestsCount int
	Id                 int `json:"id"`
	Title              string
	GitlabUpdatedAt    api.Iso8601Time  `json:"updated_at"`
	GitlabCreatedAt    api.Iso8601Time  `json:"created_at"`
	GitlabClosedAt     *api.Iso8601Time `json:"closed_at"`
	ClosedBy           struct {
		State     string
		WebUrl    string
		AvatarUrl string
		Username  string
		Id        int
		Name      string
	}
	UserNotesCount int
	DueDate        api.Iso8601Time
	WebUrl         string `json:"web_url"`
	References     struct {
		Short    string
		Relative string
		Full     string
	}
	TimeStats struct {
		TimeEstimate        int64
		TotalTimeSpent      int64
		HumanTimeEstimate   string
		HumanTotalTimeSpent string
	}
	HasTasks         bool
	TaskStatus       string
	Confidential     bool
	DiscussionLocked bool
	IssueType        string
	Serverity        string
	Links            struct {
		Self       string `json:"url"`
		Notes      string
		AwardEmoji string
		Project    string
	}
	TaskCompletionStatus struct {
		Count          int
		CompletedCount int
	}
}

func ExtractApiIssues(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_TABLE)
	config := data.Options.ScopeConfig
	var issueSeverityRegex *regexp.Regexp
	var issueComponentRegex *regexp.Regexp
	var issuePriorityRegex *regexp.Regexp
	var issueTypeBugRegex *regexp.Regexp
	var issueTypeRequirementRegex *regexp.Regexp
	var issueTypeIncidentRegex *regexp.Regexp
	var issueSeverity = config.IssueSeverity
	var err error
	if len(issueSeverity) > 0 {
		issueSeverityRegex, err = regexp.Compile(issueSeverity)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile issueSeverity failed")
		}
	}
	var issueComponent = config.IssueComponent
	if len(issueComponent) > 0 {
		issueComponentRegex, err = regexp.Compile(issueComponent)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile issueComponent failed")
		}
	}
	var issuePriority = config.IssuePriority
	if len(issuePriority) > 0 {
		issuePriorityRegex, err = regexp.Compile(issuePriority)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile issuePriority failed")
		}
	}
	var issueTypeBug = config.IssueTypeBug
	if len(issueTypeBug) > 0 {
		issueTypeBugRegex, err = regexp.Compile(issueTypeBug)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile issueTypeBug failed")
		}
	}
	var issueTypeRequirement = config.IssueTypeRequirement
	if len(issueTypeRequirement) > 0 {
		issueTypeRequirementRegex, err = regexp.Compile(issueTypeRequirement)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile issueTypeRequirement failed")
		}
	}
	var issueTypeIncident = config.IssueTypeIncident
	if len(issueTypeIncident) > 0 {
		issueTypeIncidentRegex, err = regexp.Compile(issueTypeIncident)
		if err != nil {
			return errors.Default.Wrap(err, "regexp Compile issueTypeIncident failed")
		}
	}
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			body := &IssuesResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}

			if body.ProjectId == 0 {
				return nil, nil
			}

			results := make([]interface{}, 0, 2)
			gitlabIssue, err := convertGitlabIssue(body, data.Options.ProjectId)
			if err != nil {
				return nil, err
			}

			for _, label := range body.Labels {
				results = append(results, &models.GitlabIssueLabel{
					IssueId:      gitlabIssue.GitlabId,
					LabelName:    label,
					ConnectionId: data.Options.ConnectionId,
				})
				if issueSeverityRegex != nil {
					groups := issueSeverityRegex.FindStringSubmatch(label)
					if len(groups) > 1 {
						gitlabIssue.Severity = groups[1]
					}
				}

				if issueComponentRegex != nil {
					groups := issueComponentRegex.FindStringSubmatch(label)
					if len(groups) > 1 {
						gitlabIssue.Component = groups[1]
					}
				}

				if issuePriorityRegex != nil {
					groups := issuePriorityRegex.FindStringSubmatch(label)
					if len(groups) > 1 {
						gitlabIssue.Priority = groups[1]
					}
				}

				if issueTypeBugRegex != nil {
					if ok := issueTypeBugRegex.MatchString(label); ok {
						gitlabIssue.StdType = ticket.BUG
						gitlabIssue.Type = label
					}
				}

				if issueTypeRequirementRegex != nil {
					if ok := issueTypeRequirementRegex.MatchString(label); ok {
						gitlabIssue.StdType = ticket.REQUIREMENT
						gitlabIssue.Type = label
					}
				}

				if issueTypeIncidentRegex != nil {
					if ok := issueTypeIncidentRegex.MatchString(label); ok {
						gitlabIssue.StdType = ticket.INCIDENT
						gitlabIssue.Type = label
					}
				}
			}
			gitlabIssue.ConnectionId = data.Options.ConnectionId
			if body.Author != nil {
				gitlabAuthor, err := convertGitlabAuthor(body, data.Options.ConnectionId)
				if err != nil {
					return nil, err
				}
				results = append(results, gitlabAuthor)
			}
			results = append(results, gitlabIssue)

			for _, v := range body.Assignees {
				assignee := &models.GitlabAccount{
					ConnectionId: data.Options.ConnectionId,
					Username:     v.Username,
					Name:         v.Name,
					State:        v.State,
					AvatarUrl:    v.AvatarUrl,
					WebUrl:       v.WebUrl,
				}
				issueAssignee := &models.GitlabIssueAssignee{
					ConnectionId: data.Options.ConnectionId,
					GitlabId:     gitlabIssue.GitlabId,
					ProjectId:    gitlabIssue.ProjectId,
					AssigneeId:   v.Id,
					AssigneeName: v.Username,
				}
				results = append(results, assignee, issueAssignee)
			}

			return results, nil
		},
	})

	if err != nil {
		return errors.Convert(err)
	}

	return extractor.Execute()
}

func convertGitlabIssue(issue *IssuesResponse, projectId int) (*models.GitlabIssue, errors.Error) {
	gitlabIssue := &models.GitlabIssue{
		GitlabId:        issue.Id,
		ProjectId:       projectId,
		Number:          issue.Iid,
		State:           issue.State,
		Title:           issue.Title,
		Body:            issue.Description,
		Url:             issue.Links.Self,
		ClosedAt:        api.Iso8601TimeToTime(issue.GitlabClosedAt),
		GitlabCreatedAt: issue.GitlabCreatedAt.ToTime(),
		GitlabUpdatedAt: issue.GitlabUpdatedAt.ToTime(),
		TimeEstimate:    issue.TimeStats.TimeEstimate,
		TotalTimeSpent:  issue.TimeStats.TotalTimeSpent,
		CreatorId:       issue.Author.Id,
		CreatorName:     issue.Author.Username,
	}

	if issue.Assignee != nil {
		gitlabIssue.AssigneeId = issue.Assignee.Id
		gitlabIssue.AssigneeName = issue.Assignee.Username
	}
	if issue.Author != nil {
		gitlabIssue.CreatorId = issue.Author.Id
		gitlabIssue.CreatorName = issue.Author.Username
	}
	if issue.GitlabClosedAt != nil {
		gitlabIssue.LeadTimeMinutes = uint(issue.GitlabClosedAt.ToTime().Sub(issue.GitlabCreatedAt.ToTime()).Minutes())
	}

	return gitlabIssue, nil
}

func convertGitlabAuthor(issue *IssuesResponse, connectionId uint64) (*models.GitlabAccount, errors.Error) {
	gitlabAuthor := &models.GitlabAccount{
		ConnectionId: connectionId,
		GitlabId:     issue.Author.Id,
		Username:     issue.Author.Username,
		Name:         issue.Author.Name,
		State:        issue.Author.State,
		AvatarUrl:    issue.Author.AvatarUrl,
		WebUrl:       issue.Author.WebUrl,
	}

	return gitlabAuthor, nil
}

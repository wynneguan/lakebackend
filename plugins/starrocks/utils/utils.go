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

package utils

import (
	"fmt"
	"strings"
)

// GetTablesByDomainLayer return the tables of the DomainLayer
func GetTablesByDomainLayer(domainLayer string) []string {
	switch domainLayer {
	case "code":
		return []string{
			"commits_diffs",
			"pull_requests",
			"commits",
			"refs_pr_cherrypicks",
			"repos",
			"refs",
			"pull_request_commits",
			"repo_commits",
			"pull_request_labels",
			"commit_parents",
			"notes",
			"pull_request_comments",
			"commit_files",
		}
	case "crossdomain":
		return []string{
			"pull_request_issues",
			"users",
			"issue_commits",
			"issue_repo_commits",
			"refs_issues_diffs",
			"board_repos",
		}
	case "devops":
		return []string{
			"builds",
			"jobs",
		}
	case "ticket":
		return []string{
			"board_issues",
			"boards",
			"changelogs",
			"issue_comments",
			"issue_labels",
			"issues",
			"sprints",
			"issue_worklogs",
			"board_sprints",
			"sprint_issues",
		}

	}
	return nil
}

func hasPrefixes(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func stringIn(s string, l ...string) bool {
	for _, item := range l {
		if s == item {
			return true
		}
	}
	return false
}

// GetStarRocksDataType analysis and return the data type of StarRocks
func GetStarRocksDataType(dataType string) string {
	dataType = strings.ToLower(dataType)
	starrocksDatatype := "string"
	if hasPrefixes(dataType, "datetime", "timestamp") {
		starrocksDatatype = "datetime"
	} else if stringIn(dataType, "date") {
		starrocksDatatype = "date"
	} else if strings.HasPrefix(dataType, "bigint") || stringIn(dataType, "bigserial") {
		starrocksDatatype = "bigint"
	} else if stringIn(dataType, "char") {
		starrocksDatatype = "char"
	} else if stringIn(dataType, "int", "integer", "serial") {
		starrocksDatatype = "int"
	} else if stringIn(dataType, "tinyint(1)", "boolean") {
		starrocksDatatype = "boolean"
	} else if stringIn(dataType, "smallint", "smallserial") {
		starrocksDatatype = "smallint"
	} else if stringIn(dataType, "real") {
		starrocksDatatype = "float"
	} else if stringIn(dataType, "numeric", "double precision") {
		starrocksDatatype = "double"
	} else if stringIn(dataType, "decimal") {
		starrocksDatatype = "decimal"
	} else if stringIn(dataType, "json", "jsonb") {
		starrocksDatatype = "json"
	} else if dataType == "uuid" {
		starrocksDatatype = "char(36)"
	} else if strings.HasSuffix(dataType, "[]") {
		starrocksDatatype = fmt.Sprintf("array<%s>", GetStarRocksDataType(strings.Split(dataType, "[]")[0]))
	}
	return starrocksDatatype
}

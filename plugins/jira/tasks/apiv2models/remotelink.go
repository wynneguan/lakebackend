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

package apiv2models

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"gorm.io/datatypes"
)

type RemoteLink struct {
	ID          uint64 `json:"id"`
	Self        string `json:"self"`
	GlobalID    string `json:"globalId"`
	Application struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"application"`
	Relationship string `json:"relationship"`
	Object       struct {
		URL     string `json:"url"`
		Title   string `json:"title"`
		Summary string `json:"summary"`
		Icon    struct {
			URL16X16 string `json:"url16x16"`
			Title    string `json:"title"`
		} `json:"icon"`
		Status struct {
			Resolved bool `json:"resolved"`
			Icon     struct {
				URL16X16 string `json:"url16x16"`
				Title    string `json:"title"`
				Link     string `json:"link"`
			} `json:"icon"`
		} `json:"status"`
	} `json:"object"`
}

func (r RemoteLink) ToToolLayer(connectionId, issueId uint64, raw json.RawMessage) *models.JiraRemotelink {
	return &models.JiraRemotelink{
		ConnectionId: connectionId,
		RemotelinkId: r.ID,
		IssueId:      issueId,
		Self:         r.Self,
		Title:        r.Object.Title,
		Url:          r.Object.URL,
		RawJson:      datatypes.JSON(raw),
	}
}

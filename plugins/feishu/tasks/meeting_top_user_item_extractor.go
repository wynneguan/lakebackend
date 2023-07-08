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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/feishu/models"
)

var _ plugin.SubTaskEntryPoint = ExtractMeetingTopUserItem

func ExtractMeetingTopUserItem(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*FeishuTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: FeishuApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Table: RAW_MEETING_TOP_USER_ITEM_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			body := &models.FeishuMeetingTopUserItem{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}
			rawInput := &api.DatePair{}
			rawErr := errors.Convert(json.Unmarshal(row.Input, rawInput))
			if rawErr != nil {
				return nil, rawErr
			}
			results := make([]interface{}, 0)
			results = append(results, &models.FeishuMeetingTopUserItem{
				ConnectionId:    data.Options.ConnectionId,
				StartTime:       rawInput.PairStartTime,
				MeetingCount:    body.MeetingCount,
				MeetingDuration: body.MeetingDuration,
				Name:            body.Name,
				UserType:        body.UserType,
			})
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractMeetingTopUserItemMeta = plugin.SubTaskMeta{
	Name:             "extractMeetingTopUserItem",
	EntryPoint:       ExtractMeetingTopUserItem,
	EnabledByDefault: true,
	Description:      "Extract raw top user meeting data into tool layer table feishu_meeting_top_user_item",
}

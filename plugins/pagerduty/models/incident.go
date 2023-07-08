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

package models

import (
	"github.com/apache/incubator-devlake/core/models/common"
	"time"
)

const (
	IncidentStatusAcknowledged IncidentStatus = "acknowledged"
	IncidentStatusTriggered    IncidentStatus = "triggered"
	IncidentStatusResolved     IncidentStatus = "resolved"
)

type (
	IncidentUrgency string
	IncidentStatus  string

	Incident struct {
		common.NoPKModel
		ConnectionId uint64 `gorm:"primaryKey"`
		Number       int    `gorm:"primaryKey"`
		Url          string
		ServiceId    string
		Summary      string
		Status       IncidentStatus  //acknowledged, triggered, resolved
		Urgency      IncidentUrgency //high or low
		CreatedDate  time.Time
		UpdatedDate  time.Time
	}
)

func (Incident) TableName() string {
	return "_tool_pagerduty_incidents"
}

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

package raw

import "time"

type Notifications struct {
	// Address corresponds to the JSON schema field "address".
	Address *string `json:"address,omitempty"`

	// Id corresponds to the JSON schema field "id".
	Id *string `json:"id,omitempty"`

	// StartedAt corresponds to the JSON schema field "started_at".
	StartedAt *time.Time `json:"started_at,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type *string `json:"type,omitempty"`

	// User corresponds to the JSON schema field "user".
	User *NotificationsUser `json:"user,omitempty"`
}

type NotificationsUser struct {
	// HtmlUrl corresponds to the JSON schema field "html_url".
	HtmlUrl *string `json:"html_url,omitempty"`

	// Id corresponds to the JSON schema field "id".
	Id *string `json:"id,omitempty"`

	// Self corresponds to the JSON schema field "self".
	Self *string `json:"self,omitempty"`

	// Summary corresponds to the JSON schema field "summary".
	Summary *string `json:"summary,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type *string `json:"type,omitempty"`
}

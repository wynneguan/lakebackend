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

type Incident struct {
	// Acknowledgements corresponds to the JSON schema field "acknowledgements".
	Acknowledgements []IncidentsAcknowledgementsElem `json:"acknowledgements,omitempty"`

	// Alerts corresponds to the JSON schema field "alerts".
	Alerts []IncidentsAlertsElem `json:"alerts,omitempty"`

	// Assignments corresponds to the JSON schema field "assignments".
	Assignments []IncidentsAssignmentsElem `json:"assignments,omitempty"`

	// ConferenceBridge corresponds to the JSON schema field "conference_bridge".
	ConferenceBridge *IncidentsConferenceBridge `json:"conference_bridge,omitempty"`

	// CreatedAt corresponds to the JSON schema field "created_at".
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// EscalationPolicy corresponds to the JSON schema field "escalation_policy".
	EscalationPolicy *IncidentsEscalationPolicy `json:"escalation_policy,omitempty"`

	// FirstTriggerLogEntry corresponds to the JSON schema field
	// "first_trigger_log_entry".
	FirstTriggerLogEntry *IncidentsFirstTriggerLogEntry `json:"first_trigger_log_entry,omitempty"`

	// HtmlUrl corresponds to the JSON schema field "html_url".
	HtmlUrl *string `json:"html_url,omitempty"`

	// Id corresponds to the JSON schema field "id".
	Id *string `json:"id,omitempty"`

	// IncidentKey corresponds to the JSON schema field "incident_key".
	IncidentKey *string `json:"incident_key,omitempty"`

	// IncidentNumber corresponds to the JSON schema field "incident_number".
	IncidentNumber *int `json:"incident_number,omitempty"`

	// LastStatusChangeAt corresponds to the JSON schema field
	// "last_status_change_at".
	LastStatusChangeAt *time.Time `json:"last_status_change_at,omitempty"`

	// LastStatusChangeBy corresponds to the JSON schema field
	// "last_status_change_by".
	LastStatusChangeBy *IncidentsLastStatusChangeBy `json:"last_status_change_by,omitempty"`

	// LogEntries corresponds to the JSON schema field "log_entries".
	LogEntries []IncidentsLogEntriesElem `json:"log_entries,omitempty"`

	// Priority corresponds to the JSON schema field "priority".
	Priority *IncidentsPriority `json:"priority,omitempty"`

	// Self corresponds to the JSON schema field "self".
	Self *string `json:"self,omitempty"`

	// Service corresponds to the JSON schema field "service".
	Service *IncidentsService `json:"service,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status *string `json:"status,omitempty"`

	// Summary corresponds to the JSON schema field "summary".
	Summary *string `json:"summary,omitempty"`

	// Teams corresponds to the JSON schema field "teams".
	Teams []IncidentsTeamsElem `json:"teams,omitempty"`

	// Title corresponds to the JSON schema field "title".
	Title *string `json:"title,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type *string `json:"type,omitempty"`

	// Urgency corresponds to the JSON schema field "urgency".
	Urgency *string `json:"urgency,omitempty"`
}

type IncidentsAcknowledgementsElem struct {
	// Acknowledger corresponds to the JSON schema field "acknowledger".
	Acknowledger *IncidentsAcknowledgementsElemAcknowledger `json:"acknowledger,omitempty"`

	// At corresponds to the JSON schema field "at".
	At *time.Time `json:"at,omitempty"`
}

type IncidentsAcknowledgementsElemAcknowledger struct {
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

type IncidentsAlertsElem struct {
	// AlertKey corresponds to the JSON schema field "alert_key".
	AlertKey *string `json:"alert_key,omitempty"`

	// Body corresponds to the JSON schema field "body".
	Body *IncidentsAlertsElemBody `json:"body,omitempty"`

	// CreatedAt corresponds to the JSON schema field "created_at".
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// HtmlUrl corresponds to the JSON schema field "html_url".
	HtmlUrl *string `json:"html_url,omitempty"`

	// Id corresponds to the JSON schema field "id".
	Id *string `json:"id,omitempty"`

	// Incident corresponds to the JSON schema field "incident".
	Incident *IncidentsAlertsElemIncident `json:"incident,omitempty"`

	// Integration corresponds to the JSON schema field "integration".
	Integration *IncidentsAlertsElemIntegration `json:"integration,omitempty"`

	// Self corresponds to the JSON schema field "self".
	Self *string `json:"self,omitempty"`

	// Service corresponds to the JSON schema field "service".
	Service *IncidentsAlertsElemService `json:"service,omitempty"`

	// Severity corresponds to the JSON schema field "severity".
	Severity *string `json:"severity,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status *string `json:"status,omitempty"`

	// Summary corresponds to the JSON schema field "summary".
	Summary *string `json:"summary,omitempty"`

	// Suppressed corresponds to the JSON schema field "suppressed".
	Suppressed *bool `json:"suppressed,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type *string `json:"type,omitempty"`
}

type IncidentsAlertsElemBody struct {
	// Contexts corresponds to the JSON schema field "contexts".
	Contexts []IncidentsAlertsElemBodyContextsElem `json:"contexts,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type *string `json:"type,omitempty"`
}

type IncidentsAlertsElemBodyContextsElem struct {
	// Href corresponds to the JSON schema field "href".
	Href *string `json:"href,omitempty"`

	// Src corresponds to the JSON schema field "src".
	Src *string `json:"src,omitempty"`

	// Text corresponds to the JSON schema field "text".
	Text *string `json:"text,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type *string `json:"type,omitempty"`
}

type IncidentsAlertsElemIncident struct {
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

type IncidentsAlertsElemIntegration struct {
	// HtmlUrl corresponds to the JSON schema field "html_url".
	HtmlUrl *string `json:"html_url,omitempty"`

	// Id corresponds to the JSON schema field "id".
	Id *string `json:"id,omitempty"`

	// Name corresponds to the JSON schema field "name".
	Name *string `json:"name,omitempty"`

	// Self corresponds to the JSON schema field "self".
	Self *string `json:"self,omitempty"`

	// Service corresponds to the JSON schema field "service".
	Service *IncidentsAlertsElemIntegrationService `json:"service,omitempty"`

	// Summary corresponds to the JSON schema field "summary".
	Summary *string `json:"summary,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type *string `json:"type,omitempty"`
}

type IncidentsAlertsElemIntegrationService struct {
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

type IncidentsAlertsElemService struct {
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

type IncidentsAssignmentsElem struct {
	// Assignee corresponds to the JSON schema field "assignee".
	Assignee *IncidentsAssignmentsElemAssignee `json:"assignee,omitempty"`

	// At corresponds to the JSON schema field "at".
	At *time.Time `json:"at,omitempty"`
}

type IncidentsAssignmentsElemAssignee struct {
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

type IncidentsConferenceBridge struct {
	// ConferenceNumber corresponds to the JSON schema field "conference_number".
	ConferenceNumber *string `json:"conference_number,omitempty"`

	// ConferenceUrl corresponds to the JSON schema field "conference_url".
	ConferenceUrl *string `json:"conference_url,omitempty"`
}

type IncidentsEscalationPolicy struct {
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

type IncidentsFirstTriggerLogEntry struct {
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

type IncidentsLastStatusChangeBy struct {
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

type IncidentsLogEntriesElem struct {
	// Agent corresponds to the JSON schema field "agent".
	Agent *IncidentsLogEntriesElemAgent `json:"agent,omitempty"`

	// Channel corresponds to the JSON schema field "channel".
	Channel *IncidentsLogEntriesElemChannel `json:"channel,omitempty"`

	// CreatedAt corresponds to the JSON schema field "created_at".
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// EventDetails corresponds to the JSON schema field "event_details".
	EventDetails *IncidentsLogEntriesElemEventDetails `json:"event_details,omitempty"`

	// HtmlUrl corresponds to the JSON schema field "html_url".
	HtmlUrl *string `json:"html_url,omitempty"`

	// Id corresponds to the JSON schema field "id".
	Id *string `json:"id,omitempty"`

	// Incident corresponds to the JSON schema field "incident".
	Incident *IncidentsLogEntriesElemIncident `json:"incident,omitempty"`

	// Self corresponds to the JSON schema field "self".
	Self *string `json:"self,omitempty"`

	// Summary corresponds to the JSON schema field "summary".
	Summary *string `json:"summary,omitempty"`

	// Teams corresponds to the JSON schema field "teams".
	Teams []IncidentsLogEntriesElemTeamsElem `json:"teams,omitempty"`

	// Type corresponds to the JSON schema field "type".
	Type *string `json:"type,omitempty"`
}

type IncidentsLogEntriesElemAgent struct {
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

type IncidentsLogEntriesElemChannel struct {
	// Type corresponds to the JSON schema field "type".
	Type *string `json:"type,omitempty"`
}

type IncidentsLogEntriesElemEventDetails struct {
	// Description corresponds to the JSON schema field "description".
	Description *string `json:"description,omitempty"`
}

type IncidentsLogEntriesElemIncident struct {
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

type IncidentsLogEntriesElemTeamsElem struct {
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

type IncidentsPriority struct {
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

type IncidentsService struct {
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

type IncidentsTeamsElem struct {
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

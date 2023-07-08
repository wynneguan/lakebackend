#!/bin/sh
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

. "$(dirname $0)/../../vars/active-vars.sh"

jira_endpoint=${1-$JIRA_ENDPOINT}
jira_username=${2-$JIRA_USERNAME}
jira_password=${3-$JIRA_PASSWORD}

curl -sv $LAKE_ENDPOINT/plugins/jira/connections --data @- <<JSON | jq
{
    "name": "testjira",
    "endpoint": "$jira_endpoint",
    "username": "$jira_username",
    "password": "$jira_password"
}
JSON

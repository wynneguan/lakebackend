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

package main

import (
	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/plugins/jenkins/impl"
	"github.com/spf13/cobra"
)

var PluginEntry impl.Jenkins

func main() {
	jenkinsCmd := &cobra.Command{Use: "jenkins"}
	connectionId := jenkinsCmd.Flags().Uint64P("connection", "c", 1, "jenkins connection id")
	jobFullName := jenkinsCmd.Flags().StringP("jobFullName", "j", "", "jenkins job full name")
	deployTagPattern := jenkinsCmd.Flags().String("deployTagPattern", "(?i)deploy", "deploy tag name")
	timeAfter := jenkinsCmd.Flags().StringP("timeAfter", "a", "", "collect data that are created after specified time, ie 2006-05-06T07:08:09Z")

	jenkinsCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"connectionId":     *connectionId,
			"jobFullName":      *jobFullName,
			"deployTagPattern": *deployTagPattern,
			"timeAfter":        *timeAfter,
		})
	}
	runner.RunCmd(jenkinsCmd)
}

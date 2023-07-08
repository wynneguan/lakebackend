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

package helper

import (
	"fmt"
	"os"

	"github.com/apache/incubator-devlake/core/config"
)

// UseMySQL FIXME
func UseMySQL(host string, port int) string {
	conn := fmt.Sprintf("mysql://merico:merico@%s:%d/lake?charset=utf8mb4&parseTime=True", host, port)
	_ = os.Setenv("E2E_DB_URL", conn)
	config.GetConfig().Set("E2E_DB_URL", conn)
	return conn
}

// UsePostgres FIXME
func UsePostgres(host string, port int) string {
	conn := fmt.Sprintf("postgres://merico:merico@%s:%d/lake", host, port)
	_ = os.Setenv("E2E_DB_URL", conn)
	config.GetConfig().Set("E2E_DB_URL", conn)
	return conn
}

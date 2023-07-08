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

package api

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/core/errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type CSTTimeRecord struct {
	Created CSTTime
}

type CSTTimeRecordP struct {
	Created *CSTTime
}

func TestCSTTime(t *testing.T) {
	pairs := map[string]time.Time{
		`{ "Created": "2021-07-30 19:14:33" }`: TimeMustParse("2021-07-30T11:14:33Z"),
		`{ "Created": "2021-07-30" }`:          TimeMustParse("2021-07-29T16:00:00Z"),
	}

	for input, expected := range pairs {
		var record CSTTimeRecord
		err := errors.Convert(json.Unmarshal([]byte(input), &record))
		assert.Nil(t, err)
		assert.Equal(t, expected, time.Time(record.Created).UTC())
	}
}

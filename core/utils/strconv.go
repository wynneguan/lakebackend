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
	"github.com/apache/incubator-devlake/core/errors"
	"strconv"
	"time"
)

// StrToIntOr Return defaultValue if text is empty, or try to convert it to int
func StrToIntOr(text string, defaultValue int) (int, errors.Error) {
	if text == "" {
		return defaultValue, nil
	}
	return errors.Convert01(strconv.Atoi(text))
}

// StrToDurationOr Return defaultValue if text is empty, or try to convert it to time.Duration
func StrToDurationOr(text string, defaultValue time.Duration) (time.Duration, errors.Error) {
	if text == "" {
		return defaultValue, nil
	}
	return errors.Convert01(time.ParseDuration(text))
}

// StrToBoolOr Return defaultValue if text is empty, or try to convert it to bool
func StrToBoolOr(text string, defaultValue bool) (bool, errors.Error) {
	if text == "" {
		return defaultValue, nil
	}
	return errors.Convert01(strconv.ParseBool(text))
}

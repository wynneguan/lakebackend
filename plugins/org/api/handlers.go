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
	"encoding/csv"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/gocarina/gocsv"
	"net/http"
)

const maxMemory = 32 << 20 // 32 MB

type Handlers struct {
	store store
}

func NewHandlers(basicRes context.BasicRes) *Handlers {
	return &Handlers{store: NewDbStore(basicRes.GetDal(), basicRes)}
}

func (h *Handlers) unmarshal(r *http.Request, items interface{}) errors.Error {
	if r == nil {
		return errors.Default.New("request is nil")
	}
	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(maxMemory); err != nil {
			return errors.Convert(err)
		}
	}
	f, fh, err := r.FormFile("file")
	if err != nil {
		return errors.Convert(err)
	}
	f.Close()
	file, err := fh.Open()
	if err != nil {
		return errors.Convert(err)
	}
	defer file.Close()
	return errors.Convert(gocsv.UnmarshalCSV(csv.NewReader(file), items))
}

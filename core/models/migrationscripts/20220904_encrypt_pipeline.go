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

package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*encryptPipeline)(nil)

type encryptPipeline struct{}

type PipelineEncryption0904 struct {
	archived.Model
	Plan string
}

func (script *encryptPipeline) Up(basicRes context.BasicRes) errors.Error {
	encryptionSecret := basicRes.GetConfig(plugin.EncodeKeyEnvStr)
	if encryptionSecret == "" {
		return errors.BadInput.New("invalid encryptionSecret")
	}
	err := migrationhelper.TransformColumns(
		basicRes,
		script,
		"_devlake_pipelines",
		[]string{"plan"},
		func(src *PipelineEncryption0904) (*PipelineEncryption0904, errors.Error) {
			plan, err := plugin.Encrypt(encryptionSecret, src.Plan)
			if err != nil {
				return nil, err
			}
			return &PipelineEncryption0904{
				Model: src.Model,
				Plan:  plan,
			}, nil
		},
	)
	return err
}

func (*encryptPipeline) Version() uint64 {
	return 20220904162121
}

func (*encryptPipeline) Name() string {
	return "encrypt Pipeline"
}

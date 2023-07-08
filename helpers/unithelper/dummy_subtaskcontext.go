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

package unithelper

import (
	"github.com/apache/incubator-devlake/core/dal"
	mockplugin "github.com/apache/incubator-devlake/mocks/core/plugin"
	"github.com/stretchr/testify/mock"
)

// DummySubTaskContext FIXME ...
func DummySubTaskContext(db dal.Dal) *mockplugin.SubTaskContext {
	mockCtx := new(mockplugin.SubTaskContext)
	mockCtx.On("GetDal").Return(db)
	mockCtx.On("GetLogger").Return(DummyLogger())
	mockCtx.On("SetProgress", mock.Anything, mock.Anything)
	mockCtx.On("IncProgress", mock.Anything, mock.Anything)
	mockCtx.On("GetName").Return("test")
	return mockCtx
}

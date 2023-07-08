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

package tasks

import (
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"reflect"
)

var ConnectUserAccountsExactMeta = plugin.SubTaskMeta{
	Name:             "connectUserAccountsExact",
	EntryPoint:       ConnectUserAccountsExact,
	EnabledByDefault: true,
	Description:      "associate users and accounts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func ConnectUserAccountsExact(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*TaskData)
	var users []crossdomain.User
	err := db.All(&users)
	if err != nil {
		return err
	}
	emails := make(map[string]string)
	names := make(map[string]string)
	for _, user := range users {
		if user.Email != "" {
			emails[user.Email] = user.Id
		}
		if user.Name != "" {
			names[user.Name] = user.Id
		}
	}
	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From(&crossdomain.Account{}),
		dal.Where("id NOT IN (SELECT account_id FROM user_accounts)"),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(crossdomain.Account{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: Params{
				ConnectionId: data.Options.ConnectionId,
			},
			Table: "users",
		},

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			account := inputRow.(*crossdomain.Account)
			if userId, ok := emails[account.Email]; account.Email != "" && ok {
				return []interface{}{
					&crossdomain.UserAccount{
						UserId:    userId,
						AccountId: account.Id,
					},
				}, nil
			}
			if userId, ok := names[account.FullName]; account.FullName != "" && ok {
				return []interface{}{
					&crossdomain.UserAccount{
						UserId:    userId,
						AccountId: account.Id,
					},
				}, nil
			}
			if userId, ok := names[account.UserName]; account.UserName != "" && ok {
				return []interface{}{
					&crossdomain.UserAccount{
						UserId:    userId,
						AccountId: account.Id,
					},
				}, nil
			}
			return nil, nil
		},
	})
	if err != nil {
		return err
	}
	return converter.Execute()
}

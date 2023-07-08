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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"reflect"
)

// DataConvertHandler Accept row from source cursor, return list of entities that need to be stored
type DataConvertHandler func(row interface{}) ([]interface{}, errors.Error)

// DataConverterArgs includes the arguments about DataConverter.
// This will be used in Creating a DataConverter.
//
//	DataConverterArgs {
//				InputRowType: 		type of inputRow ,
//				Input:        		dal cursor,
//				RawDataSubTaskArgs: args about raw data task
//				Convert: 			main function including conversion logic
//				BatchSize: 			batch size
type DataConverterArgs struct {
	RawDataSubTaskArgs
	// Domain layer entity Id prefix, i.e. `jira:JiraIssue:1`, `github:GithubIssue`
	InputRowType reflect.Type
	Input        dal.Rows
	Convert      DataConvertHandler
	BatchSize    int
}

// DataConverter helps you convert Data from Tool Layer Tables to Domain Layer Tables
// It reads rows from specified Iterator, and feed it into `Converter` handler
// you can return arbitrary domain layer entities from this handler, ApiConverter would
// first delete old data by their RawDataOrigin information, and then perform a
// batch save operation for you.
type DataConverter struct {
	*RawDataSubTask
	args *DataConverterArgs
}

// NewDataConverter function helps you create a DataConverter using DataConverterArgs.
// You can see the usage in plugins/github/tasks/pr_issue_convertor.go or other convertor file.
func NewDataConverter(args DataConverterArgs) (*DataConverter, errors.Error) {
	rawDataSubTask, err := NewRawDataSubTask(args.RawDataSubTaskArgs)
	if err != nil {
		return nil, err
	}
	// process args
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	return &DataConverter{
		RawDataSubTask: rawDataSubTask,
		args:           &args,
	}, nil
}

// Execute function implements Subtask interface.
// It loads data from Tool Layer Tables using `Ctx.GetDal()`, convert Data using `converter.args.Convert` handler
// Then save data to Domain Layer Tables using BatchSaveDivider
func (converter *DataConverter) Execute() errors.Error {
	// load data from database
	db := converter.args.Ctx.GetDal()

	// batch save divider
	RAW_DATA_ORIGIN := "RawDataOrigin"
	divider := NewBatchSaveDivider(converter.args.Ctx, converter.args.BatchSize, converter.table, converter.params)

	// set progress
	converter.args.Ctx.SetProgress(0, -1)

	cursor := converter.args.Input
	defer cursor.Close()
	ctx := converter.args.Ctx.GetContext()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}
		inputRow := reflect.New(converter.args.InputRowType).Interface()
		err := db.Fetch(cursor, inputRow)
		if err != nil {
			return errors.Default.Wrap(err, "error fetching rows")
		}

		results, err := converter.args.Convert(inputRow)
		if err != nil {
			return errors.Default.Wrap(err, "error calling Converter plugin implementation")
		}

		for _, result := range results {
			// get the batch operator for the specific type
			batch, err := divider.ForType(reflect.TypeOf(result))
			if err != nil {
				return errors.Default.Wrap(err, "error getting batch from result")
			}
			// set raw data origin field
			origin := reflect.ValueOf(result).Elem().FieldByName(RAW_DATA_ORIGIN)
			if origin.IsValid() {
				origin.Set(reflect.ValueOf(inputRow).Elem().FieldByName(RAW_DATA_ORIGIN))
			}
			// records get saved into db when slots were max outed
			err = batch.Add(result)
			if err != nil {
				return errors.Default.Wrap(err, "error adding result to batch")
			}
		}
		converter.args.Ctx.IncProgress(1)
	}

	// save the last batches
	return divider.Close()
}

// Check if DataConverter implements SubTask interface
var _ plugin.SubTask = (*DataConverter)(nil)

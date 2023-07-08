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

package util

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GenerateAllFormatVar fill all format var into values
func GenerateAllFormatVar(values map[string]string, baseVarName, baseValue string) {
	values[strcase.ToLowerCamel(baseVarName)] = strcase.ToLowerCamel(baseValue)
	values[strcase.ToCamel(baseVarName)] = strcase.ToCamel(baseValue)
	values[strcase.ToSnake(baseVarName)] = strcase.ToSnake(baseValue)
	values[strcase.ToScreamingSnake(baseVarName)] = strcase.ToScreamingSnake(baseValue)
	values[strcase.ToKebab(baseVarName)] = strcase.ToKebab(baseValue)
	values[strcase.ToScreamingKebab(baseVarName)] = strcase.ToScreamingKebab(baseValue)
}

// ReadTemplate read a file to string
func ReadTemplate(templateFile string) string {
	f, err := os.ReadFile(templateFile)
	cobra.CheckErr(err)
	return string(f)
}

// WriteTemplates write some strings to files
func WriteTemplates(path string, templates map[string]string) {
	err := os.MkdirAll(path, 0755)
	cobra.CheckErr(err)
	for name, template := range templates {
		err := os.MkdirAll(filepath.Dir(filepath.Join(path, name)), 0755)
		cobra.CheckErr(err)
		err = os.WriteFile(filepath.Join(path, name), []byte(template), 0600)
		cobra.CheckErr(err)
		println(filepath.Join(path, name), ` generated`)
	}
}

// ReplaceVarInFile replacte var into file without reading
func ReplaceVarInFile(filename string, reg *regexp.Regexp, new string) {
	f, err := os.ReadFile(filename)
	cobra.CheckErr(err)
	f = reg.ReplaceAll(f, []byte(new))

	err = os.WriteFile(filename, f, 0777)
	cobra.CheckErr(err)
	println(filename, ` updated`)
}

// DetectExistVars filter the used vars in templates
func DetectExistVars(templates map[string]string, values map[string]string) (newValues map[string]string) {
	newValues = map[string]string{}
	for varName, value := range values {
		for _, template := range templates {
			if strings.Contains(template, varName) {
				newValues[varName] = value
				break
			}
		}
	}
	return newValues
}

// ReplaceVarInTemplates replace var with templates into templates
func ReplaceVarInTemplates(templates map[string]string, valueMap map[string]string) {
	for i, template := range templates {
		templates[i] = ReplaceVars(template, valueMap)
	}
}

// ReplaceVars will replace s with valueMap and return it
func ReplaceVars(s string, valueMap map[string]string) string {
	for varName, value := range valueMap {
		s = ReplaceVar(s, varName, value)
	}
	return s
}

// ReplaceVar will replace s with value and return it
func ReplaceVar(s, varName, value string) string {
	return strings.ReplaceAll(s, fmt.Sprintf(`{{ .%s }}`, varName), value)
}

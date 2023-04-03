// Copyright 2023 Shenry Tech AB
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shdbcli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/shenrytech/shdb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"sigs.k8s.io/yaml"
)

func outputJson(tr *shdb.TypeRegistry, obj shdb.IObject) error {
	o := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	fmt.Println(o.Format(obj))
	return nil
}

func outputYaml(tr *shdb.TypeRegistry, obj shdb.IObject) error {
	o := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	j, err := o.Marshal(obj)
	if err != nil {
		return err
	}
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		return err
	}
	fmt.Println(string(y))
	return nil
}

func defaultTemplate(tr *shdb.TypeRegistry, obj shdb.IObject, format string) error {
	switch format {
	case "brief":
		return outputTmpl(tr, obj, "Type: {{.metadata.type}}, UUID: {{.metadata.uuid}}, Last Updated: {{.metadata.updated_at}}")
	case "detailed":
		fmt.Println(prototext.Format(obj))
		return nil
	case "list":
		return outputTmpl(tr, obj, "{{.metadata.uuid}}\t{{.metadata.type}}\t{{.metadata.updated_at}}")
	}

	return fmt.Errorf("default format %s not found", format)
}

func outputDefault(tr *shdb.TypeRegistry, obj shdb.IObject, format string) error {
	mi, err := tr.GetMessageInfo(obj.GetMetadata().TypeId().TypeKey())
	if err != nil {
		return fmt.Errorf("type with fullname %s not found in type registry", obj.GetMetadata().ProtoReflect().Descriptor().FullName())
	}
	t, ok := mi.PrintTemplates[format]
	if !ok {
		return defaultTemplate(tr, obj, format)
	}
	return outputTmpl(tr, obj, t)
}

func outputTmpl(tr *shdb.TypeRegistry, obj shdb.IObject, tmpl string) error {
	o := protojson.MarshalOptions{
		UseProtoNames: true,
	}
	data, err := o.Marshal(obj)
	if err != nil {
		return err
	}
	a := map[string]interface{}{}
	err = json.Unmarshal(data, &a)
	if err != nil {
		return err
	}
	t, err := template.New("test").Parse(tmpl)
	if err != nil {
		return err
	}
	err = t.Execute(os.Stdout, a)
	fmt.Println()
	return err

}

func output(tr *shdb.TypeRegistry, obj shdb.IObject, format string) error {
	switch format {
	case "json":
		return outputJson(tr, obj)
	case "yaml":
		return outputYaml(tr, obj)
	case "brief", "detailed", "list":
		return outputDefault(tr, obj, format)
	default:
		return outputTmpl(tr, obj, format)
	}
}

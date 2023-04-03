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
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/shenrytech/shdb"
	"github.com/spf13/cobra"
)

func completeType(cli *shdb.Client, toComplete string) ([]string, cobra.ShellCompDirective) {
	rsp, err := cli.GetTypeNames()
	if err != nil {
		log.Printf("failed to retrieve types and aliases %v", err)
		return nil, cobra.ShellCompDirectiveError
	}
	res := []string{}
	for k, v := range rsp {
		res = append(res, v...)
		res = append(res, k)
	}
	return res, cobra.ShellCompDirectiveNoFileComp
}

func completeId(cli *shdb.Client, toa string, toComplete string) ([]string, cobra.ShellCompDirective) {
	tk, err := cli.TypeRegistry().GetTypeKeyFromToA(toa)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	res := []string{}
	ch, err := cli.SearchRef(tk, func(obj *shdb.ObjRef) bool {
		id, err := uuid.FromBytes(obj.Uuid)
		if err != nil {
			log.Printf("invalid uuid in database %x", obj.Uuid)
			return false
		}
		return strings.HasPrefix(id.String(), toComplete)
	})
	if err != nil {
		panic(err)
	}
	for ref := range ch {
		id, err := uuid.FromBytes(ref.Uuid)
		if err != nil {
			panic(err)
		}
		res = append(res, id.String())
	}

	return res, cobra.ShellCompDirectiveNoFileComp
}

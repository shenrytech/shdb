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
	"fmt"

	"github.com/google/uuid"
	"github.com/shenrytech/shdb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Cli struct {
}

var ccAccessor func() *grpc.ClientConn

func ValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cc := ccAccessor()
	cli := shdb.NewClient(cmd.Context(), cc)
	switch len(args) {
	case 0:
		return completeType(cli, toComplete)
	case 1:
		return completeId(cli, args[0], toComplete)
	}
	return nil, cobra.ShellCompDirectiveError
}

func get(cmd *cobra.Command, args []string) error {
	cc := ccAccessor()
	cli := shdb.NewClient(cmd.Context(), cc)

	id, err := uuid.Parse(args[1])
	if err != nil {
		return err
	}
	bid, err := id.MarshalBinary()
	if err != nil {
		return err
	}
	tk, err := cli.TypeRegistry().GetTypeKeyFromToA(args[0])
	if err != nil {
		return err
	}
	ref := &shdb.ObjRef{
		Type: tk[:],
		Uuid: bid,
	}
	obj, err := cli.Get(cmd.Context(), *ref.TypeId())
	if err != nil {
		return err
	}
	fmt.Println(obj)
	return nil
}

// versionCmd represents the version command
var getCmd = &cobra.Command{
	Use:               "get <fullname|alias> id",
	Short:             "retrieve an IObject",
	RunE:              get,
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: ValidArgsFunction,
}

func AddGet(ctx context.Context, parent *cobra.Command, ccAccess func() *grpc.ClientConn) {
	ccAccessor = ccAccess
	parent.AddCommand(getCmd)
}

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
	"github.com/google/uuid"
	"github.com/shenrytech/shdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Cli struct {
}

var ccAccessor func() *grpc.ClientConn

func ValidTypeIdArgFn(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

func ValidTypeArgFn(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cc := ccAccessor()
	cli := shdb.NewClient(cmd.Context(), cc)
	return completeType(cli, toComplete)
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
	return output(cli.TypeRegistry(), obj, viper.GetString("output"))
}

func list(cmd *cobra.Command, args []string) error {
	cc := ccAccessor()
	cli := shdb.NewClient(cmd.Context(), cc)

	tk, err := cli.TypeRegistry().GetTypeKeyFromToA(args[0])
	if err != nil {
		return err
	}

	obj, err := cli.List(cmd.Context(), tk)
	if err != nil {
		return err
	}
	tr := cli.TypeRegistry()
	for _, v := range obj {
		if err := output(tr, v, "list"); err != nil {
			return err
		}
	}
	return nil
}

// versionCmd represents the version command
var getCmd = &cobra.Command{
	Use:               "get <fullname|alias> id",
	Short:             "retrieve an IObject",
	RunE:              get,
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: ValidTypeIdArgFn,
}

var listCmd = &cobra.Command{
	Use:               "list <fullname|alias>",
	Short:             "list objects of a specific type",
	RunE:              list,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: ValidTypeArgFn,
}

func AddGet(ctx context.Context, parent *cobra.Command, ccAccess func() *grpc.ClientConn) {
	ccAccessor = ccAccess
	getCmd.PersistentFlags().StringP("output", "o", "yaml", "output format [json|yaml|brief|list|detail|\"<go template>\"]")
	viper.BindPFlag("output", getCmd.PersistentFlags().Lookup("output"))
	parent.AddCommand(getCmd)
	parent.AddCommand(listCmd)
}

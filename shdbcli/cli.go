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
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/shenrytech/shdb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Cli struct {
}

type ShdbContextKey struct{}
type ShdbContextVal struct {
	ClientConnAccessor func(state interface{}) *grpc.ClientConn
}

func ValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ctxVal, ok := cmd.Context().Value(ShdbContextKey{}).(ShdbContextVal)
	if !ok {
		fmt.Printf("failed to connect to server")
		return nil, cobra.ShellCompDirectiveError
	}
	cc := ctxVal.ClientConnAccessor(nil)
	cli := shdb.NewClient(cmd.Context(), cc)
	s := strings.Split(toComplete, " ")
	switch len(s) {
	case 1:
		return completeType(cli, s[0])
	case 2:
		return completeId(cli, s[0], s[1])
	}
	return nil, cobra.ShellCompDirectiveError
}

func get(cmd *cobra.Command, args []string) error {
	ctxVal, ok := cmd.Context().Value(ShdbContextKey{}).(ShdbContextVal)
	if !ok {
		return errors.New("failed to connect to server")
	}
	cc := ctxVal.ClientConnAccessor(nil)
	cli := shdb.NewObjectServiceClient(cc)
	schema, err := cli.GetSchema(cmd.Context(), &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to fetch db schema: [%v]", err)
	}
	tr := shdb.NewTypeRegistry()
	if err := tr.UseFileDescriptorSet(schema); err != nil {
		fmt.Printf("failed to use db schema [%v]", err)
	}
	id, err := uuid.Parse(args[1])
	if err != nil {
		return err
	}
	bid, err := id.MarshalBinary()
	if err != nil {
		return err
	}
	tk, err := tr.GetTypeKeyFromToA(args[0])
	if err != nil {
		return err
	}
	ref := &shdb.ObjRef{
		Type: tk[:],
		Uuid: bid,
	}
	obj, err := cli.Get(cmd.Context(), &shdb.GetReq{Ref: ref})
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

func AddGet(ctx context.Context, parent *cobra.Command, ccAccessor func(state interface{}) *grpc.ClientConn) {
	ct := context.WithValue(ctx, ShdbContextKey{}, ShdbContextVal{ClientConnAccessor: ccAccessor})
	getCmd.SetContext(ct)
	parent.AddCommand(getCmd)
}

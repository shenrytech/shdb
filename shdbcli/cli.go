package shdbcli

import (
	"fmt"
	"log"

	"github.com/shenrytech/shdb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Cli struct {
}

type CtxClientConn struct{}

func GetTypesAndAliases(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	cc := cmd.Context().Value(CtxClientConn{}).(*grpc.ClientConn)
	cli := shdb.NewObjectServiceClient(cc)
	rsp, err := cli.GetTypeNames(cmd.Context(), &emptypb.Empty{})
	if err != nil {
		log.Printf("failed to retrieve types and aliases %v", err)
		return nil, cobra.ShellCompDirectiveError
	}
	res := []string{}
	for _, v := range rsp.TypeAliases {
		res = append(res, v.Aliased...)
		res = append(res, v.Fullname)
	}
	return res, cobra.ShellCompDirectiveNoFileComp
}

func get(cmd *cobra.Command, args []string) {
	fullNameOrAlias := args[0]
	cc := cmd.Context().Value(CtxClientConn{}).(*grpc.ClientConn)
	cli := shdb.NewObjectServiceClient(cc)
	schema, err := cli.GetSchema(cmd.Context(), &emptypb.Empty{})
	if err != nil {
		fmt.Printf("failed to fetch db schema: [%v]", err)
		return
	}
	tr := shdb.NewTypeRegistry()
	if err := tr.UseFileDescriptorSet(schema); err != nil {
		fmt.Printf("failed to use db schema [%v]", err)
	}
	// TODO: We create and list object.

}

// versionCmd represents the version command
var getCmd = &cobra.Command{
	Use:               "get <fullname|alias>",
	Short:             "retrieve an IObject",
	Run:               get,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: GetTypesAndAliases,
}

func AddGet(ctx context.Context, parent *cobra.Command, cc *grpc.ClientConn) {
	ct := context.WithValue(ctx, CtxClientConn{}, cc)
	getCmd.SetContext(ct)
	parent.AddCommand(getCmd)
}

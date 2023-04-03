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

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shenrytech/shdb/shdbcli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cc *grpc.ClientConn

var rootCmd = &cobra.Command{
	Use:           "shdbcli",
	Short:         "cli for shdb database servers",
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		cc, err = grpc.Dial(viper.GetString("address"),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		return
	},
}

func main() {

	defer func() {
		if cc != nil {
			cc.Close()
		}
	}()

	rootCmd.PersistentFlags().String("address", "localhost:3335", "address to database server")
	viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))
	shdbcli.AddGet(context.Background(), rootCmd, func() *grpc.ClientConn { return cc })

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

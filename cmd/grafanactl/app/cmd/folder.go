// Copyright © 2019 Juan Mesa <linuxven@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gosuri/uitable"

	gapi "github.com/overdrive3000/go-grafana-api"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	id  int64
	uid string
)

// SetUpClient set up a new grafana client
func SetUpClient() (*gapi.Client, error) {
	log.Debugf("Setting up grafana client with url %s and key %s", viper.GetString("url"), viper.GetString("apiKey"))
	return gapi.New(
		viper.GetString("apiKey"),
		viper.GetString("url"),
	)
}

// folderCmd represents the folder command
func folderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "folder <get, list, create, delete>",
		Short: "Grafana Folders",
		Long: `Perform operations against Grafana Folders:

* Create Folders
* List Folders
* Search Folders`,
	}
	cmd.AddCommand(getFolderCmd())
	cmd.AddCommand(listFoldersCmd())

	return cmd
}

func listFoldersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all Folders",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug("Listing Grafana Folders")
			client, _ := SetUpClient()
			var folders []gapi.Folder
			var out []byte
			folders, err := client.Folders()
			if err != nil {
				log.Error(err)
			}

			switch cmd.Flag("output").Value.String() {
			case "json":
				out, err = json.Marshal(folders)
				if err != nil {
					log.Error(err)
				}
			case "table":
				out = formatAsTable(folders, 60)
			default:
				log.Error(errors.New(fmt.Sprintf("unknown output format %q", cmd.Flag("output").Value.String())))
			}

			fmt.Fprintln(os.Stdout, string(out))
		},
	}

	return cmd
}

func getFolderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Search a folder",
		Long: `Search a folder either by ID or UID
grafanactl folder get 
[--id <value>]
[--uid <value>]`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flag("id").Changed && !cmd.Flag("uid").Changed {
				return errors.New("Either --id or --uid must be specified")
			}
			if cmd.Flag("id").Changed && cmd.Flag("uid").Changed {
				return errors.New("Just one flag --id or --uid is allowed at a time")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug("Getting Grafana Folder")
			client, _ := SetUpClient()
			var err error
			var folder *gapi.Folder
			var out []byte
			if cmd.Flag("id").Changed {
				folder, err = client.Folder(id)
			}
			if cmd.Flag("uid").Changed {
				folder, err = client.FolderByUID(uid)
			}
			if err != nil {
				log.Error(err)
			}
			switch cmd.Flag("output").Value.String() {
			case "json":
				out, err = json.Marshal(folder)
				if err != nil {
					log.Error(err)
				}
			case "table":
				folders := []gapi.Folder{*folder}
				out = formatAsTable(folders, 60)
			default:
				log.Error(errors.New(fmt.Sprintf("unknown output format %q", cmd.Flag("output").Value.String())))
			}

			fmt.Fprintln(os.Stdout, string(out))
		},
	}
	cmd.Flags().Int64Var(&id, "id", 0, "Folder ID to search")
	cmd.Flags().StringVar(&uid, "uid", "", "Folder UID to search")

	return cmd
}

func formatAsTable(folders []gapi.Folder, colWidth uint) []byte {
	tbl := uitable.New()
	log.Debugf("%#v", folders)
	tbl.MaxColWidth = colWidth
	tbl.AddRow("ID", "UID", "TITLE")
	for i := 0; i <= len(folders)-1; i++ {
		f := folders[i]
		tbl.AddRow(f.Id, f.Uid, f.Title)
	}
	return tbl.Bytes()
}

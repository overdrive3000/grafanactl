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
	"io/ioutil"
	"os"

	"github.com/gosuri/uitable"

	gapi "github.com/overdrive3000/go-grafana-api"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	cmd.AddCommand(createFolderCmd())
	cmd.AddCommand(deleteFolderCmd())

	return cmd
}

func deleteFolderCmd() *cobra.Command {
	var uid string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a Folder",
		Long: `Delete a grafana folder
grafanactl delete --uid <value>`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug("Deleting Folder")
			client, _ := SetUpClient()
			err := client.DeleteFolder(uid)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			fmt.Printf("Folder %s deleted\n", uid)
		},
	}

	cmd.Flags().StringVar(&uid, "uid", "", "Folder UID")
	cmd.MarkFlagRequired("uid")

	return cmd
}

func createFolderCmd() *cobra.Command {
	var file, uid, title string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Folder",
		Long: `Creates a new Grafana Folder
grafanactl create [--title <value> [--uid <value]]
grafanactl create [--file|-f <value>]`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flag("file").Changed && (cmd.Flag("uid").Changed || cmd.Flag("title").Changed) {
				return errors.New("Flag --file cannot be using with --uid or --title")
			}
			if !cmd.Flag("file").Changed && !cmd.Flag("title").Changed {
				return errors.New("Either --file or --title must be specified")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Debug("Creating a new Folder")
			client, _ := SetUpClient()
			var folder gapi.Folder
			var err error
			if cmd.Flag("file").Changed {
				file, err := ioutil.ReadFile(cmd.Flag("file").Value.String())
				if err != nil {
					log.Debug(err)
					os.Exit(1)
				}
				_ = json.Unmarshal(file, &folder)
				uid = folder.Uid
				title = folder.Title
			}
			switch uid {
			case "":
				folder, err = client.NewFolder(title)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			default:
				folder, err = client.NewFolderWithUID(title, uid)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}
			fmt.Printf("folder id %d created\n", folder.Id)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Create folder with title <value>")
	cmd.Flags().StringVar(&uid, "uid", "", "Create folder with UID <value>")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Create folder from json file <value>")

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
				os.Exit(1)
			}

			switch cmd.Flag("output").Value.String() {
			case "json":
				out, err = json.Marshal(folders)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			case "table":
				out = formatAsTable(folders, 60)
			default:
				log.Error(errors.New(fmt.Sprintf("unknown output format %q", cmd.Flag("output").Value.String())))
				os.Exit(1)
			}

			fmt.Fprintln(os.Stdout, string(out))
		},
	}

	return cmd
}

func getFolderCmd() *cobra.Command {
	var id int64
	var uid string
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
				os.Exit(1)
			}
			switch cmd.Flag("output").Value.String() {
			case "json":
				out, err = json.Marshal(folder)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			case "table":
				folders := []gapi.Folder{*folder}
				out = formatAsTable(folders, 60)
			default:
				log.Error(errors.New(fmt.Sprintf("unknown output format %q", cmd.Flag("output").Value.String())))
				os.Exit(1)
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

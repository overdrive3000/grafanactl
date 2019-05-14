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
	"fmt"

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

	return cmd
}

// folderByUIDCmd search a folder by UID in Grafana
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
			if cmd.Flag("id").Changed {
				folder, err = client.Folder(id)
			}
			if cmd.Flag("uid").Changed {
				folder, err = client.FolderByUID(uid)
			}
			if err != nil {
				log.Error(err)
			}
			fmt.Println(folder)
		},
	}
	cmd.Flags().Int64Var(&id, "id", 0, "Folder ID to search")
	cmd.Flags().StringVar(&uid, "uid", "", "Folder UID to search")

	return cmd
}

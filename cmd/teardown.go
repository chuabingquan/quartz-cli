/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"quartz/api"

	"github.com/spf13/cobra"
)

// teardownCmd represents the teardown command
var teardownCmd = &cobra.Command{
	Use:   "teardown",
	Short: "Teardown current Quartz project",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		js := api.NewJobService()

		projectName := args[0]

		jobs, err := js.Jobs()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		projectExist := false
		foundProject := api.Job{}
		for _, job := range jobs {
			if job.Name == projectName {
				projectExist = true
				foundProject = job
				break
			}
		}

		if !projectExist {
			fmt.Println("Project to teardown does not exist.")
			return
		}

		res, err := js.DeleteJob(foundProject.ID)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(teardownCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// teardownCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// teardownCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

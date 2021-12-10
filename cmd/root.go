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
	"os/exec"

	"github.com/spf13/cobra"
)

var destPath string
var excludeRepoString string
var includeRepoString string
var parallelism int
var requiredEnvs = []string{
	"GITHUB_TOKEN",
}
var skipArchived bool

func checkRequirements() error {
	for _, env := range requiredEnvs {
		_, exist := os.LookupEnv(env)
		if !exist {
			errMsg := fmt.Sprintf("Environment variable %s is required", env)
			return errors.New(errMsg)
		}
	}

	_, err := exec.LookPath("git")
	if err != nil {
		errMsg := fmt.Sprintf("Could not find command git. Please install it: %s\n", err)
		return errors.New(errMsg)
	}

	if includeRepoString != "" && excludeRepoString != "" {
		errMsg := fmt.Sprint("You can't use both --exclude-repos and --include-repos")
		return errors.New(errMsg)
	}

	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "github-org-repos-sync <github-org-name>",
	Short: "Sync github org repos",
	Long:  "Sync github org repos.",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkRequirements()
	},
	Run: func(cmd *cobra.Command, args []string) {
		main(args)
	},
	Args:    cobra.ExactArgs(1),
	Example: "github-org-repos-sync floorpunch",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().BoolVar(&skipArchived, "skip-archived", false, "Skip archived repos?")
	rootCmd.Flags().StringVar(&excludeRepoString, "exclude-repos", "", "Exclude repos that contain string")
	rootCmd.Flags().StringVarP(&destPath, "destination-path", "d", ".", "Destionation path for repos")
	rootCmd.Flags().StringVar(&includeRepoString, "include-repos", "", "Include only repos that contain string")
	rootCmd.Flags().IntVarP(&parallelism, "parallelism", "p", 1, "Number of parallel git operations")
}

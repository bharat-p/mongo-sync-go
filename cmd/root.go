// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

	"github.com/bharat-p/goutils/cli"

	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	Local struct {
		Username     string
		Password     string
		Host         string
		Port         int
		AuthDatabase string // authentication database
		Database     string // database to dump/restore
	}
	Remote struct {
		Username string
		Password string

		Host         string
		Port         int
		AuthDatabase string // authentication database
		Database     string // database to dump/restore
	}
	General struct {
		MongodumpOptions    string `mapstructure:"mongodump_options"`
		MongorestoreOptions string `mapstructure:"mongorestore_options"`
	}
}

var Config config
var cfgFile string

var DryRun = true
var Verbose = false

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "go-mongo-sync",
	Short: "Sync local and remote Mongo Database instances",
	Long:  `Sync data b/w local and remote Mongo Databases using mongodump and mongorestore.`,
}

func testBinary(executableName string) bool {

	exists := true
	cmd := exec.Command("type", executableName)
	_, err := cmd.CombinedOutput()

	if err != nil {
		exists = false
		fmt.Errorf("%s", err)
	}

	return exists
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mongo-sync.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	if testBinary("mongodump") == false {
		log.Fatalf("mongodump not found. Please make sure you have mongodump command installed.")
	}

	if testBinary("mongorestore") == false {
		log.Fatalf("mongorestore not found. Please make sure you have mongodump command installed.")
	}

	RootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "D", false, "Dry Run")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "V", false, "Debug mode")

	RootCmd.PersistentFlags().StringP("local.host", "", "", "Local mongo host(overrides value from config file)")
	RootCmd.PersistentFlags().IntP("local.port", "", 0, "Local mongo port (overrides value from config file)")
	RootCmd.PersistentFlags().StringP("local.database", "", "", "Local mongo database name to push/pull(overrides value from config file) ")

	RootCmd.PersistentFlags().StringP("remote.host", "", "", "Remote mongo host (overrides value from config file)")
	RootCmd.PersistentFlags().IntP("remote.port", "", 0, "Local mongo port(overrides value from config file)")
	RootCmd.PersistentFlags().StringP("remote.database", "", "", "Remote mongo database name to push/pull (overrides value from config file)")

	RootCmd.PersistentFlags().StringP("mongodump-args", "", "", "Additional arguments to pass to mongodump command(overrides value from config file) ")
	RootCmd.PersistentFlags().StringP("mongorestore-args", "", "", "Additional arguments to pass to mongorestore command(overrides value from config file) ")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if Verbose {
		cli.Debug = true
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}
	if DryRun {
		cli.DryRun = true
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(".mongo-sync") // name of config file (without extension)
	viper.AddConfigPath("./")          // search in current working directory first
	viper.AddConfigPath("$HOME")       // adding home directory as second search path

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file:%s", viper.ConfigFileUsed())
	}
	viper.Unmarshal(&Config)

	if val, _ := RootCmd.Flags().GetString("local.database"); val != "" {
		Config.Local.Database, _ = RootCmd.Flags().GetString("local.database")
	}
	if val, _ := RootCmd.Flags().GetString("local.host"); val != "" {
		Config.Local.Host, _ = RootCmd.Flags().GetString("local.host")
	}

	if val, _ := RootCmd.Flags().GetInt("local.port"); val > 0 {
		Config.Local.Port, _ = RootCmd.Flags().GetInt("local.port")
	}

	if val, _ := RootCmd.Flags().GetString("remote.database"); val != "" {
		Config.Remote.Database, _ = RootCmd.Flags().GetString("remote.database")
	}

	if val, _ := RootCmd.Flags().GetString("remote.host"); val != "" {
		Config.Remote.Host, _ = RootCmd.Flags().GetString("remote.host")
	}

	if val, _ := RootCmd.Flags().GetInt("remote.port"); val > 0 {
		Config.Remote.Port, _ = RootCmd.Flags().GetInt("remote.port")
	}

	if Config.Local.Port == 0 {
		Config.Local.Port = 27017
	}

	if Config.Remote.Port == 0 {
		Config.Remote.Port = 27017
	}

}

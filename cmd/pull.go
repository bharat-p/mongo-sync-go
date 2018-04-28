package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bharat-p/goutils/cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rmiCmd represents the rmi command
var rmiCmd = &cobra.Command{
	Use:   "pull",
	Short: "Sync data from remote mongo database to local mongo database",
	Long: `Sync data from remote mongo database to local mongo database
:
go-mongo-sync pull --database=mydb
`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Running go-mongo-sync pull")
		localMongoHost := Config.Local.Host
		localMongoPort := Config.Local.Port
		localDatabaseName := Config.Local.Database

		remoteMongoHost := Config.Remote.Host
		remoteMongoPort := Config.Remote.Port
		remoteDatabaseName := Config.Remote.Database

		if remoteDatabaseName == "" {
			log.Error("Remote Database name can not be empty for pull command, please specify --remote.database or --local.database.")
			cmd.Usage()
			os.Exit(1)
		}

		if remoteMongoHost == "" {
			log.Error("No remote server specified. Make sure you have created config file .mongo-sync.yaml with appropriate configuration")
			cmd.Usage()
			os.Exit(1)
		}

		if localMongoHost == "" {
			log.Error("No local server specified. Make sure you have created config file .mongo-sync.yaml with appropriate configuration")
			cmd.Usage()
			os.Exit(1)
		}

		if localDatabaseName == "" {
			localDatabaseName = remoteDatabaseName
		}

		log.Infof("Syncing mongo database: %s from [%s:%d] to [%s:%d] as database: %s",
			remoteDatabaseName, remoteMongoHost, remoteMongoPort, localMongoHost, localMongoPort, localDatabaseName)

		dir, err := ioutil.TempDir("", fmt.Sprintf("%s-", Config.Remote.Database))
		if err != nil {
			log.Fatal(err)
		}

		defer os.RemoveAll(dir) // clean up

		dumpCmdArgs := fmt.Sprintf(`--host=%s --port=%d --username="%s" --password="%s" 
			--authenticationDatabase=%s --out=%s --db=%s %s`,
			remoteMongoHost, remoteMongoPort, Config.Remote.Username, Config.Remote.Password,
			Config.Remote.AuthDatabase, dir, remoteDatabaseName, Config.General.MongodumpOptions)

		exitCode, err := cli.RunCommand("mongodump", strings.Fields(dumpCmdArgs)...)
		if exitCode != 0 {
			log.Errorf("Failed to create database dump.[%s]", err)
			os.Exit(1)
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s", dir, remoteDatabaseName)); err == nil {
			restoreCmdArgs := fmt.Sprintf(`--host=%s --port=%d --username=%s --password=%s 
			--authenticationDatabase=%s --dir=%s/%s --db=%s %s`,
				localMongoHost, localMongoPort, Config.Local.Username, Config.Local.Password,
				Config.Local.AuthDatabase, dir, remoteDatabaseName, localDatabaseName, Config.General.MongorestoreOptions)

			exitCode, err = cli.RunCommand("mongorestore", strings.Fields(restoreCmdArgs)...)
			if exitCode != 0 {
				log.Errorf("Failed to restore database.[%s]", err)
				os.Exit(1)
			} else {
				log.Infof("Done syncing database=%s from host=%s to host=%s as database=%s",
					remoteDatabaseName, remoteMongoHost, localMongoHost, localDatabaseName)
			}
		} else {
			log.Errorf("Database dump directory: %s/%s not found. Probably trying to sync empty database?", dir, remoteDatabaseName)
			os.Exit(1)
		}

	},
}

func init() {
	RootCmd.AddCommand(rmiCmd)
}

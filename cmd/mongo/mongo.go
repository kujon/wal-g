package mongo

import (
	"fmt"
	"os"
	"strings"

	"github.com/wal-g/wal-g/internal"
	"github.com/wal-g/wal-g/internal/webserver"

	"github.com/spf13/cobra"
	"github.com/wal-g/tracelog"
)

var DBShortDescription = "MongoDB backup tool"

// These variables are here only to show current version. They are set in makefile during build process
var WalgVersion = "devel"
var GitRevision = "devel"
var BuildDate = "devel"
var webServer webserver.WebServer

var Cmd = &cobra.Command{
	Use:     "wal-g",
	Short:   DBShortDescription, // TODO : improve description
	Version: strings.Join([]string{WalgVersion, GitRevision, BuildDate, "MongoDB"}, "\t"),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := internal.AssertRequiredSettingsSet()
		tracelog.ErrorLogger.FatalOnError(err)

		httpListenAddr, httpListen := internal.GetSetting(internal.HttpListen)
		if httpListen {
			webServer = webserver.NewSimpleWebServer(httpListenAddr)
			tracelog.ErrorLogger.FatalOnError(webServer.Serve())
		}

		exposePprof, err := internal.GetBoolSetting(internal.HttpExposePprof, false)
		tracelog.ErrorLogger.FatalOnError(err)
		if exposePprof {
			internal.RequiredSettings[internal.HttpListen] = true
			err := internal.AssertRequiredSettingsSet()
			tracelog.ErrorLogger.FatalfOnError(internal.HttpExposePprof + " failed: %v", err)
			webserver.EnablePprofEndpoints(webServer)
		}
	},
}

func Execute() {
	if err := Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(internal.InitConfig, internal.Configure)

	internal.RequiredSettings[internal.MongoDBUriSetting] = true
	Cmd.PersistentFlags().StringVar(&internal.CfgFile, "config", "", "config file (default is $HOME/.wal-g.yaml)")
	Cmd.InitDefaultVersionFlag()
	internal.AddConfigFlags(Cmd)
}

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/shawnsilva/piccolo/log"
	"github.com/shawnsilva/piccolo/piccolo"
	"github.com/shawnsilva/piccolo/utils"
	"github.com/shawnsilva/piccolo/version"
)

var (
	flagConfigFile       = flag.String("config", "conf/config.json", "Path to config file")
	flagDumpConfigFormat = flag.Bool("dumpconf", false, "If enabled, piccolo will dump a sample config file and exit. Uses config as path.")
	flagVersion          = flag.Bool("version", false, "Print the version and exit.")

	appVersion = &version.Info{}
	conf       *utils.Config
	bot        *piccolo.Bot
)

func usage() {
	fmt.Fprintf(os.Stderr, "\nUsage of %s:\n\n", (os.Args[0]))
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Parse()

	var err error
	appVersion.ParseVersion()
	if *flagVersion {
		fmt.Println(os.Args[0], "version:", appVersion.GetVersionString())
		os.Exit(0)
	}

	*flagConfigFile = filepath.ToSlash(*flagConfigFile)

	if *flagDumpConfigFormat {
		dumpFilePath := utils.StrConcat([]string{*flagConfigFile, ".sample"})
		log.WithFields(log.Fields{
			"dumpFilePath": filepath.FromSlash(dumpFilePath),
		}).Info("Dumping sample config file.")
		err = utils.DumpConfigFormat(dumpFilePath)
		if err != nil {
			log.WithFields(log.Fields{
				"dumpFilePath": dumpFilePath,
				"error":        err,
			}).Fatal("Error Writing sample config.")
		}
		os.Exit(0)
	}

	conf, err = utils.LoadConfig(*flagConfigFile)
	if err != nil {
		log.WithFields(log.Fields{
			"configFile": *flagConfigFile,
			"error":      err,
		}).Fatal("Error Loading Config.")
	}

	bot = piccolo.NewBot(conf, appVersion)
}

func main() {

	bot.Start()

	sigIntChannel := make(chan os.Signal, 1)
	cleanupDoneChannel := make(chan bool)
	signal.Notify(sigIntChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		for _ = range sigIntChannel {
			fmt.Println("\nReceived Shutdown Request, shutting down...")
			// do stuff
			bot.Stop()
			cleanupDoneChannel <- true
		}
	}()
	<-cleanupDoneChannel

}

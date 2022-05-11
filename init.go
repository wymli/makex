package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/wymli/makex/internal/config"
)

func init() {
	initLog()
	initConfig()
}

func initConfig() {
	if err := config.InitMakexConfig(); err != nil {
		log.Fatal(err)
	}
}

func initLog() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	formatter := log.TextFormatter{
		EnvironmentOverrideColors: true,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		DisableLevelTruncation:    true,
		QuoteEmptyFields:          true,
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			seps := strings.Split(frame.File, string(filepath.Separator))
			for i := len(seps); i < 3; i++ {
				seps = append(seps, "")
			}

			seps = seps[len(seps)-3:]

			fileName := filepath.Join(seps...) + ":" + strconv.Itoa(frame.Line)
			return "", fileName
		},
	}
	log.SetFormatter(&formatter)
}

package logger

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Read environment variables
	Verbose, _ := strconv.ParseBool(os.Getenv("VERBOSE"))

	// Initialize logger
	Log = logrus.New()
	Log.SetOutput(os.Stdout)

	if Verbose {
		Log.SetLevel(logrus.DebugLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			FullTimestamp:             true,
		})
		Log.Debug("Verbose logging enabled")
	} else {
		Log.SetLevel(logrus.InfoLevel)
		Log.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			DisableTimestamp:          true,
		})
	}
}

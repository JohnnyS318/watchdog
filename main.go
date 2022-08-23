package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type WatchdogConfig struct {
	Receiver    []string
	Watchtarget string
	Interval    time.Duration
	Cache       []byte
}

var Logger *zap.SugaredLogger
var Config *WatchdogConfig

func main() {

	//logging
	registerLogger()

	Logger.Infof("Hello, world! Starting up...")

	// ENV
	godotenv.Load()

	watch := flag.String("w", "https://google.com", "watches this url. example: https:google.com")
	receivers := flag.String("t", "", "sends a mail to this address. example: \"abc@gmail.com, def@gmail.com\"")
	interval := flag.Int("i", 10, "checks the website every x seconds. example: 10")

	flag.Parse()

	Logger.Infof("Hello Watching.... %s", *watch)

	parsedList := strings.Split(*receivers, ",")
	recList := make([]string, 0)
	for _, receiver := range parsedList {
		if receiver == "" || receiver == " " {
			continue
		}

		if !CheckEmailValidity(receiver) {
			Logger.Fatalf("Invalid email address: %s", receiver)
			panic("Invalid receiver email address")
		}

		Logger.Infof("Registered receiver: %s", receiver)
		recList = append(recList, receiver)
	}

	Config = &WatchdogConfig{
		Receiver:    recList,
		Watchtarget: *watch,
		Interval:    time.Duration(*interval) * time.Second,
		Cache:       nil,
	}

	// EMAIL Registration
	client, err := RegisterMailClient()

	if err != nil {
		Logger.Errorf("error while registering mail client: %v. Exiting now....", err)
		return
	}

	// err = SendMail(client)

	// if err != nil {
	// 	Logger.Errorf("Error while sending mail: %v", err)
	// }
	// Logger.Debugf("Email sent successfully")

	// Watchdog timer

	ticker := time.NewTicker(Config.Interval)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	defer func() {
		Logger.Infof("Defer")
	}()

	Logger.Infof("Starting Watchdog timer...")

	errCount := 0

	go func() {
		for {
			select {
			case <-ticker.C:
				Logger.Debugf("Watching now...")
				data, err := CopyWeb(Config.Watchtarget)

				if err != nil {
					errCount += 1
					Logger.Errorf("Error while copying website: %v", err)
					Logger.Debugf("Skipping this iteration with Error Count: %d", errCount)
				}else {
					if Config.Cache == nil {
						Config.Cache = data
					} else if !CompareBuf(Config.Cache, data) {
						Logger.Infof("Website has changed!")
						err = SendMail(client)
						if err != nil {
							Logger.Errorf("Error while sending mail: %v", err)
						}
						Logger.Infof("Email sent successfully")
						Config.Cache = data
					}
				}


			case <-quit:
				Logger.Infof("Bye Bye")
				ticker.Stop()
				return
			}
		}

	}()

	<-quit
	Logger.Infof("Bye Bye Main")
	ticker.Stop()
	Logger.Infof("Exit main")
}

func registerLogger() {
	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	Logger = logger.Sugar()
}

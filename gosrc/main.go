package main

import (
	"flag"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	// versionflag : Flag for display version
	versionflag bool
	// debugLevel : Flag for enable debug level
	debugLevel bool
	// version : Default version
	version = "N/A"
	// poolingInterval : pooling interval in minutes
	poolingInterval = 30
	// command args
	commandArgs = ""
	// telegramToken : To create a bot, please contact @BotFather on telegram
	telegramToken = "None"
	// telegramID : To find an id, please contact @myidbot on telegram
	telegramID = 0

	// Icons
	icon = map[int]string{
		0: "\u2705",       // ":white_check_mark:",
		1: "\u26a0\ufe0f", // ":warning:",
		2: "\u274c",       // ":x:",
		3: "\u2754",       // ":grey_question:",
	}
)

func init() {
	// Global
	flag.BoolVar(&versionflag, "v", false, "Print build id")
	flag.BoolVar(&debugLevel, "d", false, "debug mode")

	flag.IntVar(&poolingInterval, "poolinginterval", getIntEnv("POOLING_INTERVAL", poolingInterval),
		"Pooling Interval (or use env variable : POOLING_INTERVAL)")

	// command Args
	flag.StringVar(&commandArgs, "commandargs", getStringEnv("COMMAND_ARGS", commandArgs),
		"Arguments to use for check_adaptec_raid (or use env variable : COMMAND_ARGS)")

	// Telegram
	flag.StringVar(&telegramToken, "telegramtoken", getStringEnv("TELEGRAM_TOKEN", telegramToken),
		"To create a bot, please contact @BotFather on telegram (or use env variable : TELEGRAM_TOKEN)")
	flag.IntVar(&telegramID, "telegramid", getIntEnv("TELEGRAM_ID", telegramID),
		"To find an id, please contact @myidbot on telegram (or use env variable : TELEGRAM_ID)")

	flag.Parse()

	log.SetOutput(os.Stdout)

}

func main() {
	log.Info("Starting...")
	if debugLevel {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Fatal("Cannot instantiate Telegram BOT")
	}

	laststate := -1
	args := strings.Split(commandArgs, " ")
	for {
		log.WithFields(log.Fields{
			"args": args,
		}).Debug("Running command")
		cmd := exec.Command("check_adaptec_raid", args...)
		var waitStatus syscall.WaitStatus

		out, err := cmd.CombinedOutput()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				waitStatus = exitError.Sys().(syscall.WaitStatus)

				log.WithFields(log.Fields{
					"exitcode": waitStatus.ExitStatus(),
					"out":      out,
				}).Debug("Statuscheck error")
			} else {
				log.WithFields(log.Fields{
					"err": err.Error(),
				}).Fatal("Exec error")
			}
		} else {
			// Success
			waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
			log.WithFields(log.Fields{
				"exitcode": waitStatus.ExitStatus(),
				"out":      out,
			}).Debug("Statuscheck error")
		}

		if waitStatus.ExitStatus() != laststate {
			log.Debug("State Transition")

			msg := tgbotapi.NewMessage(int64(telegramID), getIcon(waitStatus.ExitStatus())+string(out[:]))
			_, err = bot.Send(msg)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err.Error(),
				}).Debug("Notification error rescheduling in 5 minutes")
				time.Sleep(time.Duration(int(time.Minute) * 5))
				continue
			}
			log.Debug("Notification sent")
			laststate = waitStatus.ExitStatus()
		}

		log.Debugf("Waiting %d minutes", poolingInterval)
		time.Sleep(time.Duration(int(time.Minute) * poolingInterval))
	}
}

func getIcon(exitcode int) string {
	if i, ok := icon[exitcode]; ok {
		return i + " "
	}
	return ""
}

func getStringEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		log.WithFields(log.Fields{"key": key}).Info("[main] : Use custom value")
		return value
	}
	log.WithFields(log.Fields{"key": key}).Info("[main] : Use default value")
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.WithFields(log.Fields{"key": key, "err": err}).Fatal("[main] : Invalid value")
			return fallback
		}
		log.WithFields(log.Fields{"key": key}).Info("[main] : Use custom value")
		return int(v)
	}
	log.WithFields(log.Fields{"key": key}).Info("[main] : Use default value")
	return fallback
}

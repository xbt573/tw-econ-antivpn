package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/xbt573/tw-econ-antivpn/antivpn"
	"github.com/xbt573/tw-econ-antivpn/econ"
	"github.com/xbt573/tw-econ-antivpn/env"
)

var (
	playerJoinedRegex = regexp.MustCompile(`ClientID=(\d+).*?(\d+\.\d+\.\d+\.\d+)`)

	host        = env.GetDefault("TW_HOST", "localhost")
	port        = intMustParse(env.GetDefault("TW_PORT", "8303"))
	password    = env.Get("TW_PASSWORD")
	token       = env.Get("API_TOKEN")
	kickMessage = env.GetDefault("KICK_MESSAGE", "Kicked for VPN")
	banMessage  = env.GetDefault("BAN_MESSAGE", "Banned for VPN")
	banTime     = intMustParse(env.GetDefault("BAN_TIME", "60"))

	console = econ.NewECON(host, password, port)
	vpn     = antivpn.NewAntiVPN(token)

	signalChannel = make(chan os.Signal, 1)
)

func mainLoop() {
	for console.Connected {
		message, err := console.Read()
		if err != nil {
			log.Fatalln(err)
		}

		if strings.Contains(message, "player has entered the game") {
			match := playerJoinedRegex.FindStringSubmatch(message)
			checkResult, err := vpn.CheckVPN(match[2])
			if err != nil {
				log.Fatalln(err)
			}

			id := intMustParse(match[1])

			if checkResult.Ban {
				err := console.Ban(id, banTime, banMessage)
				if err != nil {
					log.Fatalln(err)
				}
			} else if checkResult.IsVPN {
				err := console.Kick(id, kickMessage)
				if err != nil {
					log.Fatalln(err)
				}
			}

			switch {
			case checkResult.Ban:
				log.Printf("Banned %v\n", match[2])

			case checkResult.IsVPN && checkResult.Cached:
				log.Printf("Kicked %v (cached)\n", match[2])
			case checkResult.IsVPN:
				log.Printf("Kicked %v\n", match[2])
			}
		}
	}
}

func init() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	signal.Notify(signalChannel, os.Interrupt)
}

func main() {
	log.Println("Starting tw-econ-antivpn...")

	err := console.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	go mainLoop()

	log.Println("Started! Waiting for server shutdown or interrupt...")

	select {
	case <-signalChannel:
		break

	case <-console.Completed:
		break
	}

	log.Println("Shutting down...")
}

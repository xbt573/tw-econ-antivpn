package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

var (
	playerJoinedRegex = regexp.MustCompile(`ClientID=(\d+) addr=(.*:\d+)`)
)

type VPNApiResponse struct {
	Security struct {
		VPN   bool `json:"vpn"`
		Proxy bool `json:"proxy"`
		Tor   bool `json:"tor"`
		Relay bool `json:"relay"`
	} `json:"security"`
}

func checkVPN(token string, ip string) (bool, error) {
	url := "https://vpnapi.io/api/" + strings.Split(ip, ":")[0] + "?key=" + token
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var response VPNApiResponse
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return false, err
	}

	if response.Security.VPN || response.Security.Proxy ||
		response.Security.Tor || response.Security.Relay {
		return true, nil
	}

	return false, nil
}

func main() {
	host, exists := os.LookupEnv("TW_HOST")
	if !exists {
		log.Fatalln("TW_HOST not set")
	}

	port, exists := os.LookupEnv("TW_PORT")
	if !exists {
		log.Fatalln("TW_PORT not set")
	}

	password, exists := os.LookupEnv("TW_PASSWORD")
	if !exists {
		log.Fatalln("TW_PASSWORD not set")
	}

	token, exists := os.LookupEnv("API_TOKEN")
	if !exists {
		log.Fatalln("API_TOKEN not set")
	}

	kickMessage, exists := os.LookupEnv("KICK_MESSAGE")
	if !exists {
		kickMessage = "Kicked for VPN"
	}

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatalln(err)
	}

	message := string(buffer[:n])
	if strings.Contains(message, "Enter password") {
		_, err = conn.Write([]byte(password + "\n"))
		if err != nil {
			log.Fatalln(err)
		}

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Fatalln(err)
		}

		message := string(buffer[:n])
		if !strings.Contains(message, "Authentication successful") {
			log.Fatalln("Wrong password or timeout")
		}
	}

	go func() {
		for {
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				log.Fatalln(err)
			}

			message := string(buffer[:n])
			if strings.Contains(message, "player has entered the game") {
				match := playerJoinedRegex.FindStringSubmatch(message)
				isVPN, err := checkVPN(token, match[2])
				if err != nil {
					log.Fatalln(err)
				}

				if isVPN {
					_, err := conn.Write(
						[]byte("kick " + match[1] + " " + kickMessage + "\n"),
					)
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		break
	}
	log.Println("Shutting down...")
}

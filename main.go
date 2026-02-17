package main

import (
	"log"
	"os/exec"
	"time"
)

// Default installation paths for Discord and TeamSpeak
var apps = map[string]string{
	"Discord": "C:\\Users\\%USERNAME%\\AppData\\Local\\Discord\\Discord.exe",
	"TS3":     "C:\\Program Files\\TeamSpeak 3 Client\\ts3client_win64.exe",
}

// Block or unblock the app via Windows Firewall
func setFirewall(appName, appPath string, block bool) {
	// action := "block"
	if !block {
		// Delete the rule
		cmd := exec.Command("cmd", "/C", `netsh advfirewall firewall delete rule name="Block `+appName+`"`)
		if err := cmd.Run(); err != nil {
			log.Printf("Error removing firewall rule for %s: %v", appName, err)
		} else {
			log.Printf("Firewall rule removed for %s", appName)
		}
		return
	}

	// Add the block rule
	cmd := exec.Command("cmd", "/C",
		`netsh advfirewall firewall add rule name="Block `+appName+`" dir=out program="`+appPath+`" action=block`,
	)
	if err := cmd.Run(); err != nil {
		log.Printf("Error adding firewall rule for %s: %v", appName, err)
	} else {
		log.Printf("Firewall rule added for %s", appName)
	}
}

func main() {
	log.Println("Firewall blocker started. Active hours: 22:00–02:00")

	for {
		now := time.Now()
		hour := now.Hour()

		if hour >= 22 || hour < 2 {
			// Active hours — block apps
			for name, path := range apps {
				setFirewall(name, path, true)
			}
		} else {
			// Outside active hours — unblock apps
			for name, path := range apps {
				setFirewall(name, path, false)
			}
		}

		time.Sleep(60 * time.Second) // check every minute
	}
}

package main

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

var targetProcesses = []string{"Discord.exe", "ts3client_win64.exe"}

func main() {
	// Optional: Hide console window by detaching from it
	// hideConsole()

	log.Println("Process killer started. Running between 22:00 and 02:00")
	runLoop()
}

func runLoop() {
	for {
		now := time.Now()
		hour := now.Hour()

		// Check if current time is between 22:00 (10 PM) and 02:00 (2 AM)
		if hour >= 22 || hour < 2 {
			log.Printf("Active hours: %02d:%02d - checking for target processes", now.Hour(), now.Minute())
			checkAndKillProcesses()
		} else {
			log.Printf("Outside active hours (%02d:%02d) - sleeping", now.Hour(), now.Minute())
		}

		time.Sleep(5 * time.Second)
	}
}

func checkAndKillProcesses() {
	procs, err := process.Processes()
	if err != nil {
		log.Printf("Error getting processes: %v", err)
		return
	}

	for _, proc := range procs {
		name, err := proc.Name()
		if err != nil {
			continue
		}

		for _, target := range targetProcesses {
			if name == target {
				log.Printf("Killing process: %s (PID %d)", name, proc.Pid)
				err := proc.Kill()
				if err != nil {
					log.Printf("Failed to kill %s: %v", name, err)
				}
			}
		}
	}
}

// Optional function to hide console window
func hideConsole() {
	// This is a Windows-specific approach to hide the console
	// Note: This may not work in all cases and requires syscall
	// Alternative: Compile with -ldflags="-H windowsgui"

	// For a simpler approach, just use the compiler flag mentioned above
}

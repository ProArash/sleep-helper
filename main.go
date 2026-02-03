package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/process"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

var targetProcesses = []string{"Discord.exe", "ts3client_win64.exe"}

func main() {
	isSvc, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if running in service: %v", err)
	}

	if isSvc {
		runService("Win64")
	} else {
		runLoop()
	}
}
func runLoop() {
	for {
		now := time.Now()
		hour := now.Hour()
		if hour >= 23 || hour < 2 {
			checkAndKillProcesses()
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
				fmt.Printf("Killing process: %s (PID %d)\n", name, proc.Pid)
				err := proc.Kill()
				if err != nil {
					log.Printf("Failed to kill %s: %v", name, err)
				}
			}
		}
	}
}

func runService(name string) {
	elog, err := eventlog.Open(name)
	if err != nil {
		return
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("%s service started.", name))

	runLoop() // Keep the loop running
}

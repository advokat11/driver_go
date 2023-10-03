package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

func main() {
	args := os.Args[1:]
	loggingEnabled := len(args) > 0 && args[0] == "log"

	var driverFiles []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".inf" {
			driverFiles = append(driverFiles, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	numDrivers := int64(len(driverFiles))
	bar := progressbar.Default(numDrivers)

	var successfulInstalls uint64
	var failedInstalls uint64
	var wg sync.WaitGroup

	var logFile *os.File
	if loggingEnabled {
		logFile, _ = os.Create("log.txt")
		defer logFile.Close()
	}

	for _, path := range driverFiles {
		driverName := filepath.Base(path)
		msg := fmt.Sprintf("Installing driver: %s", driverName)
		bar.Describe(msg)

		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			cmd := exec.Command("pnputil.exe", "/add-driver", path, "/install")
			if logFile != nil {
				cmd.Stdout = logFile
				cmd.Stderr = logFile
			}

			if err := cmd.Run(); err != nil {
				atomic.AddUint64(&failedInstalls, 1)
			} else {
				atomic.AddUint64(&successfulInstalls, 1)
			}
			bar.Add(1)
		}(path)
	}

	wg.Wait()

	bar.Finish()

	fmt.Printf("Successful installs: %d, Failed installs: %d\n", successfulInstalls, failedInstalls)
	fmt.Printf("Driver installation is complete, press any button to exit.")
	fmt.Scanln()
}

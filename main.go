package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/schollz/progressbar/v3"
)

func main() {
	logFile, err := openLogFile()
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	driverFiles, err := findDriverFiles()
	if err != nil {
		log.Fatal(err)
	}

	bar := progressbar.Default(int64(len(driverFiles)))
	var successfulInstalls uint64
	var failedInstalls uint64

	for _, path := range driverFiles {
		err := installDriver(path, logFile)
		if err != nil {
			failedInstalls++
			writeToLogFile(logFile, fmt.Sprintf("Установка драйвера %s не удалась: %v\n", path, err))
		} else {
			successfulInstalls++
			writeToLogFile(logFile, fmt.Sprintf("Драйвер %s успешно установлен\n", path))
		}
		bar.Add(1)
	}

	bar.Finish()
	printStats(successfulInstalls, failedInstalls)
}

func openLogFile() (*os.File, error) {
	logFile, err := os.Create("driver_install.log")
	if err != nil {
		return nil, err
	}

	return logFile, nil
}

func findDriverFiles() ([]string, error) {
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
	return driverFiles, err
}

func installDriver(path string, logFile *os.File) error {
	cmd := exec.Command("pnputil", "/add-driver", path, "/install")
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func writeToLogFile(logFile *os.File, text string) {
	writer := bufio.NewWriter(transform.NewWriter(logFile, charmap.Windows1251.NewEncoder()))
	writer.WriteString(text)
	writer.Flush()
}

func printStats(success uint64, failed uint64) {
	fmt.Printf("Успешно установлено: %d, Установить не удалось: %d\n", success, failed)
	fmt.Print("Установка драйверов завершена.")
}

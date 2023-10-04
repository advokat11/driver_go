package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func main() {

	driverFiles, err := findDriverFiles()
	if err != nil {
		log.Fatal(err)
	}

	bar := progressbar.Default(int64(len(driverFiles)))

	var successfulInstalls uint64
	var failedInstalls uint64

	for _, path := range driverFiles {

		driverName := filepath.Base(path)
		msg := fmt.Sprintf("Устанавливаемый драйвер: %s", driverName)

		bar.Describe(msg)

		if installDriver(path) != nil {
			failedInstalls++
		} else {
			successfulInstalls++
		}

		bar.Add(1)
	}

	bar.Finish()

	printStats(successfulInstalls, failedInstalls)

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

func installDriver(path string) error {

	cmd := exec.Command("pnputil", "/add-driver", path, "/install")
	err := cmd.Run()

	return err

}

func printStats(success uint64, failed uint64) {

	fmt.Printf("Успешно установлено: %d, Установить не удалось: %d\n", success, failed)

	fmt.Print("Установка драйверов завершена. Для выхода нажмите любую клавишу.")
	fmt.Scanln()

}

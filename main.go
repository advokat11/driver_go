package main

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	printLogo()

	logFile, err := openLogFile()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			log.Printf("Ошибка при закрытии файла: %v", err)
		}
	}()

	driverFiles, err := findDriverFiles()
	if err != nil {
		log.Fatal(err)
	}

	// Указываем вывод в os.Stdout и обновляем шаблон прогресс-бара
	bar := pb.New(len(driverFiles)).
		SetWriter(os.Stdout).
		SetRefreshRate(100 * time.Millisecond).
		SetTemplateString(`{{ bar . "[" "■" ">" " " "]" }} {{percent .}} | {{counters .}}`)
	bar.Start()

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
		bar.Increment()
	}

	bar.Finish()
	printStats(successfulInstalls, failedInstalls)
}

func printLogo() {
	// Create a beautiful ASCII-art for AQ
	logo := figure.NewFigure("AQ", "slant", true)
	color.Set(color.FgHiCyan)
	logo.Print()
	color.Unset()

	// Display "drivers" centered under the logo
	signature := "drivers"
	logoWidth := len(strings.Split(logo.String(), "\n")[0])
	padding := (logoWidth - len(signature)) / 2
	color.Set(color.FgHiCyan)
	fmt.Printf("%s%s\n", strings.Repeat(" ", padding), signature)

	// Display the version centered below "drivers"
	version := "v3.6"
	paddingVersion := (logoWidth - len(version)) / 2
	fmt.Printf("%s%s\n", strings.Repeat(" ", paddingVersion), version)
	color.Unset()
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
	if _, err := writer.WriteString(text); err != nil {
		log.Printf("Ошибка при записи в лог-файл: %v", err)
	}
	if err := writer.Flush(); err != nil {
		log.Printf("Ошибка при сохранении данных в лог-файле: %v", err)
	}
}

func printStats(success uint64, failed uint64) {
	color.Set(color.FgHiGreen)
	fmt.Printf("Успешно установлено: %d\n", success)
	color.Set(color.FgHiRed)
	fmt.Printf("Установить не удалось: %d\n", failed)
	color.Unset()
	fmt.Println("Установка драйверов завершена.")
}

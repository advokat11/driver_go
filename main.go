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
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var logMutex sync.Mutex

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

	// Создаем прогресс-бар
	bar := pb.New(len(driverFiles)).
		SetWriter(os.Stdout).
		SetRefreshRate(100 * time.Millisecond).
		SetTemplateString(`{{ bar . "[" "■" ">" " " "]" }} {{percent .}} | {{counters .}}`)
	bar.Start()

	var successfulInstalls uint64
	var failedInstalls uint64

	// Создаем канал для путей к драйверам
	driverChan := make(chan string)

	// Определяем количество рабочих горутин (например, количество CPU)
	numWorkers := runtime.NumCPU()

	// Создаем WaitGroup
	var wg sync.WaitGroup

	// Запускаем рабочие горутины
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range driverChan {
				err := installDriver(path, logFile)
				if err != nil {
					atomic.AddUint64(&failedInstalls, 1)
					writeToLogFile(logFile, fmt.Sprintf("Установка драйвера %s не удалась: %v\n", path, err))
				} else {
					atomic.AddUint64(&successfulInstalls, 1)
					writeToLogFile(logFile, fmt.Sprintf("Драйвер %s успешно установлен\n", path))
				}
				bar.Increment()
			}
		}()
	}

	// Отправляем пути к драйверам в канал
	for _, path := range driverFiles {
		driverChan <- path
	}
	close(driverChan)

	// Ожидаем завершения всех горутин
	wg.Wait()

	bar.Finish()

	successful := atomic.LoadUint64(&successfulInstalls)
	failed := atomic.LoadUint64(&failedInstalls)
	printStats(successful, failed)
}

func printLogo() {
	// Создаем красивый ASCII-арт для AQ
	logo := figure.NewFigure("AQ", "slant", true)
	color.Set(color.FgHiCyan)
	logo.Print()
	color.Unset()

	// Отображаем "drivers" по центру под логотипом
	signature := "drivers"
	logoWidth := len(strings.Split(logo.String(), "\n")[0])
	padding := (logoWidth - len(signature)) / 2
	color.Set(color.FgHiCyan)
	fmt.Printf("%s%s\n", strings.Repeat(" ", padding), signature)

	// Отображаем версию по центру под "drivers"
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
	logMutex.Lock()
	defer logMutex.Unlock()
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

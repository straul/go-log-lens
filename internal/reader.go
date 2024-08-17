package internal

import (
	"bufio"
	"context"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/semaphore"
	"os"
	"sync"
)

// StreamLogsWithProgress reads the content of a single log file line by line with a progress bar and processes each line with the provided processFunc.
func StreamLogsWithProgress(filePath string, processFunc func(string), bar *progressbar.ProgressBar) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		processFunc(scanner.Text())
		if bar != nil {
			bar.Add(len(scanner.Bytes()))
		}
	}

	return scanner.Err()
}

// StreamLogsConcurrentlyWithProgress reads the content of multiple log files concurrently with a progress bar and processes each line with the provided processFunc.
func StreamLogsConcurrentlyWithProgress(filePaths []string, processFunc func(string), bar *progressbar.ProgressBar, concurrency int) error {
	var wg sync.WaitGroup
	logChannel := make(chan string, 1000)
	sem := semaphore.NewWeighted(int64(concurrency))

	for _, filePath := range filePaths {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			sem.Acquire(context.Background(), 1) // 控制并发数
			defer sem.Release(1)

			file, err := os.Open(filePath)
			if err != nil {
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				logChannel <- scanner.Text()
				if bar != nil {
					bar.Add(len(scanner.Bytes()))
				}
			}
		}(filePath)
	}

	go func() {
		wg.Wait()
		close(logChannel)
	}()

	for log := range logChannel {
		processFunc(log)
	}

	return nil
}

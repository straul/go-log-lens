package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var levels = []string{"INFO", "WARNING", "ERROR", "DEBUG"}
var messages = []string{
	"System started",
	"User logged in",
	"File not found",
	"Connection lost",
	"Transaction completed",
	"Memory usage high",
	"Disk almost full",
	"New user registered",
	"Backup completed",
	"Permission denied",
}

var meaninglessWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "ut", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea",
	"commodo", "consequat", "duis", "aute", "irure", "dolor", "in", "reprehenderit",
	"in", "voluptate", "velit", "esse", "cillum", "dolore", "eu", "fugiat", "nulla",
	"pariatur", "excepteur", "sint", "occaecat", "cupidatat", "non", "proident",
	"sunt", "in", "culpa", "qui", "officia", "deserunt", "mollit", "anim", "id", "est", "laborum",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 定义命令行参数
	numFiles := flag.Int("files", 5, "Number of log files to generate")
	numLines := flag.Int("lines", 100, "Number of lines per log file")
	flag.Parse()

	// 创建带时间戳的目录
	outputDir := filepath.Join("generated_logs", time.Now().Format("20060102_150405"))
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// 生成日志文件
	for i := 0; i < *numFiles; i++ {
		fileName := fmt.Sprintf("logfile_%d.log", i+1)
		filePath := filepath.Join(outputDir, fileName)
		generateLogFile(filePath, *numLines)
	}

	fmt.Printf("Logs generated in directory: %s\n", outputDir)
}

func generateLogFile(filePath string, numLines int) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating log file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	for i := 0; i < numLines; i++ {
		logEntry := generateLogEntry()
		file.WriteString(logEntry)
	}
}

func generateLogEntry() string {
	timestamp := time.Now().Format(time.DateTime)
	level := levels[rand.Intn(len(levels))]
	message := messages[rand.Intn(len(messages))]

	baseLog := fmt.Sprintf("%s [%s] %s", timestamp, level, message)

	// 计算需要填充的额外字符数
	remainingLength := rand.Intn(100) + 250 - len(baseLog)
	additionalWords := generateRandomWords(remainingLength)

	return fmt.Sprintf("%s %s\n", baseLog, additionalWords)
}

func generateRandomWords(targetLength int) string {
	var words []string
	for len(strings.Join(words, " ")) < targetLength {
		words = append(words, meaninglessWords[rand.Intn(len(meaninglessWords))])
	}
	return strings.Join(words, " ")
}

package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	// 读取配置文件
	config := ReadYamlConfig("config.yaml")
	wg := new(sync.WaitGroup)
	startTime := time.Now()
	var count int64 = 0
	// 读取日志文件
	for _, dir := range config.SourceDir {
		var files, _ = GetFiles(dir)
		for i, source := range files {
			count += 1
			if config.Verbose {
				fmt.Printf("Processing log file(%d/%d): %s\n", i+1, len(files), source)
			} else {
				fmt.Printf("Processing log files：%d/%d\n", i+1, len(files))
			}
			var dest string
			if strings.ToLower(config.OutputDir) == "default" {
				dest = strings.Split(source, ".")[0] + ".json"
			} else {
				dest = filepath.Join(config.OutputDir, strings.Split(filepath.Base(source), ".")[0]+".json")
			}

			// 解析
			if config.Format == "string" {
				if config.MultiThreads {
					wg.Add(1)
					go AnalyzerStr(source, dest, config.MarshalIndent, func(filename string) {
						if config.Verbose {
							fmt.Println("Parsing log file done: ", filename)
						}
						wg.Done()
					})
				} else {
					AnalyzerStr(source, dest, config.MarshalIndent, func(filename string) {
						if config.Verbose {
							fmt.Println("Parsing log file done: ", filename)
						}
					})
				}
			} else {
				if config.MultiThreads {
					wg.Add(1)
					go Analyzer(source, dest, config.MarshalIndent, func(filename string) {
						if config.Verbose {
							fmt.Println("Parsing log file done: ", filename)
						}
						wg.Done()
					})
				} else {
					Analyzer(source, dest, config.MarshalIndent, func(filename string) {
						if config.Verbose {
							fmt.Println("Parsing log file done: ", filename)
						}
					})
				}
			}
		}
	}
	wg.Wait()
	// 统计解析用时
	if config.Timing {
		elapsedTime := time.Since(startTime)
		fmt.Printf("Total time: %s\n", elapsedTime)
		fmt.Printf("Tasks: %d\n", count)
		fmt.Printf("Average time: %s\n", time.Duration(elapsedTime.Milliseconds()/count)*time.Millisecond)
	}
	fmt.Println("It’s all done!")
}

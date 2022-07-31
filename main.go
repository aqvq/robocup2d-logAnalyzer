package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	// 如果不存在配置文件则创建
	var configFilename = "config.yaml"
	if !PathExists(configFilename) {
		WriteYamlConfig(configFilename)
	}
	// 读取配置文件
	config := ReadYamlConfig(configFilename)
	ParseCmd(config)
	if config.Verbose {
		PrintConfig(config)
	}
	// 声明变量
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
			if PathExists(dest) && !config.Overwrite {
				// 当已经存在文件时，不进行覆盖重写操作，直接跳过下面的代码
				if config.Verbose {
					fmt.Println("Skipping existing log file: ", dest)
				}
				continue
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
		if count != 0 {
			fmt.Printf("Average time: %s\n", time.Duration(elapsedTime.Nanoseconds()/count)*time.Nanosecond)
		}
	}
	fmt.Println("It’s all done!")
}

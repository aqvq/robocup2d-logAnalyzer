package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Config struct {
	SourceDir     []string `yaml:"sourceDir"`
	OutputDir     string   `yaml:"outputDir"`
	MarshalIndent bool     `yaml:"marshalIndent"`
	Verbose       bool     `yaml:"verbose"`
	Format        string   `yaml:"format"`
	MultiThreads  bool     `yaml:"multiThreads"`
	Timing        bool     `yaml:"timing"`
}

// GetFiles 获取给定路径下的所有文件
func GetFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".rcg" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ReadYamlConfig 读取配置文件信息
func ReadYamlConfig(path string) *Config {
	fmt.Println("Reading config file...")
	config := new(Config)
	if file, err := ioutil.ReadFile(path); err != nil {
		panic("Error opening configuration file")
	} else {
		err := yaml.Unmarshal(file, config)
		if err != nil {
			panic("Error parsing configuration file")
		}
	}
	config.Format = strings.ToLower(config.Format)

	if config.Verbose {
		fmt.Println("Configuration:")
		fmt.Println("- Source directory: ")
		for _, dir := range config.SourceDir {
			fmt.Println("  - " + dir)
		}
		fmt.Println("- Output directory: " + config.OutputDir)
		if config.MarshalIndent {
			fmt.Println("- Marshal indent: true")
		} else {
			fmt.Println("- Marshal indent: false")
		}
		if config.Verbose {
			fmt.Println("- Verbose: true")
		} else {
			fmt.Println("- Verbose: false")
		}
		if config.Format == "string" {
			fmt.Println("- Format: string")
		} else if config.Format == "numeric" {
			fmt.Println("- Format: numeric")
		} else {
			panic("Error parsing configuration file")
		}
		if config.MultiThreads {
			fmt.Println("- MultiThreads: true")
		} else {
			fmt.Println("- MultiThreads: false")
		}
		if config.Timing {
			fmt.Println("- Timing: true")
		} else {
			fmt.Println("- Timing: false")
		}
	}
	return config
}

package main

import (
	"flag"
	"fmt"
	"github.com/go-yaml/yaml"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Source         []string `yaml:"source"`
	Output         string   `yaml:"output"`
	Formatting     bool     `yaml:"formatting"`
	Verbose        bool     `yaml:"verbose"`
	DataType       string   `yaml:"datatype"`
	Multithreading bool     `yaml:"multithreading"`
	Timing         bool     `yaml:"timing"`
	Overwrite      bool     `yaml:"overwrite"`
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

// PathExists 判断一个文件或文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// ReadYamlConfig 读取配置文件信息
func ReadYamlConfig(path string) *Config {
	config := new(Config)
	if file, err := ioutil.ReadFile(path); err != nil {
		panic("Error opening configuration file")
	} else {
		err := yaml.Unmarshal(file, config)
		if err != nil {
			panic("Error parsing configuration file")
		}
	}
	config.DataType = strings.ToLower(config.DataType)
	return config
}

// ParseCmd 解析命令行参数
func ParseCmd(config *Config) {
	flag.StringVar(&config.DataType, "datatype", config.DataType, "json文件记录数据的类型. string表示以字符类型记录数据, numeric表示以数值类型记录数据. 使用string类型会使解析速度加快")
	flag.StringVar(&config.DataType, "d", config.DataType, "--datatype")
	flag.StringVar(&config.Output, "output", config.Output, "指定json文件的输出目录, default表示输出到源文件所在目录")
	flag.StringVar(&config.Output, "o", config.Output, "--output")
	flag.BoolVar(&config.Verbose, "verbose", config.Verbose, "是否输出详细信息")
	flag.BoolVar(&config.Verbose, "v", config.Verbose, "--verbose")
	flag.BoolVar(&config.Timing, "timing", config.Timing, "是否记录解析用时")
	flag.BoolVar(&config.Timing, "t", config.Timing, "--timing")
	flag.BoolVar(&config.Formatting, "formatting", config.Formatting, "是否格式化缩进json文件, 格式化可以让json文件结构更清晰, 但会降低解析速度并增大输出文件")
	flag.BoolVar(&config.Formatting, "f", config.Formatting, "--formatting")
	flag.BoolVar(&config.Overwrite, "overwrite", config.Overwrite, "当检测到已存在同名json文件时, 是否进行覆写操作")
	flag.BoolVar(&config.Overwrite, "w", config.Overwrite, "--overwrite")
	flag.BoolVar(&config.Multithreading, "multithreading", config.Multithreading, "是否启用多线程来提高解析速度")
	flag.BoolVar(&config.Multithreading, "m", config.Multithreading, "--multithreading")
	flag.Parse()
	if flag.NArg() > 0 {
		config.Source = flag.Args()
	}
}

// PrintConfig 打印Config结构体中的数据
func PrintConfig(config *Config) {
	fmt.Println("Configuration:")
	fmt.Println("- Source: ")
	for _, dir := range config.Source {
		fmt.Println("  - " + dir)
	}
	fmt.Println("- Output: " + config.Output)
	if config.Formatting {
		fmt.Println("- Formatting: true")
	} else {
		fmt.Println("- Formatting: false")
	}
	if config.Verbose {
		fmt.Println("- Verbose: true")
	} else {
		fmt.Println("- Verbose: false")
	}
	if config.DataType == "string" {
		fmt.Println("- DataType: string")
	} else if config.DataType == "numeric" {
		fmt.Println("- DataType: numeric")
	} else {
		panic("Error parsing configuration file")
	}
	if config.Multithreading {
		fmt.Println("- Multithreading: true")
	} else {
		fmt.Println("- Multithreading: false")
	}
	if config.Timing {
		fmt.Println("- Timing: true")
	} else {
		fmt.Println("- Timing: false")
	}
	if config.Overwrite {
		fmt.Println("- Overwrite: true")
	} else {
		fmt.Println("- Overwrite: false")
	}
}

// WriteYamlConfig 创建配置文件并写入原始数据
func WriteYamlConfig(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic("Error Creating configuration file")
	}
	_, err = file.WriteString(`# config.yaml
# created at 2022年7月31日 by Shang

# source:
# 指定包含日志文件的源文件夹, 程序会自动遍历该目录及其子目录下的所有日志文件. 注意格式:
#  - 这里输入一个目录
#  - 可以再添加一个目录
source:
  - .

# output: 指定json文件的输出目录, default表示输出到源文件所在目录
output: default

# formatting: true/false 是否格式化缩进json文件, 格式化可以让json文件结构更清晰, 但会降低解析速度并增大输出文件
formatting: false

# verbose: true/false 是否输出详细信息
verbose: true

# datatype: string/numeric json文件记录数据的类型. string表示以字符类型记录数据, numeric表示以数值类型记录数据. 使用string类型会使解析速度加快
datatype: string

# multithreading: true/false 是否启用多线程来提高解析速度
multithreading: true

# timing: true/false 是否记录解析用时
timing: true

# overwrite: true/false 当检测到已存在同名json文件时, 是否进行覆写操作
overwrite: false
`)
	if err != nil {
		panic("Error writing configuration file")
	}
}

func TrimBracket(str string) string {
	if len(str) == 0 {
		return ""
	}
	end := len(str) - 1
	start := 0
	for {
		if start < len(str) && str[start] == '(' {
			start += 1
		} else {
			break
		}
	}
	for {
		if end >= 0 && str[end] == ')' {
			end -= 1
		} else {
			end += 1
			break
		}
	}
	if start >= end {
		return ""
	} else {
		return str[start:end]
	}
}

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
	SourceDir     []string `yaml:"sourceDir"`
	OutputDir     string   `yaml:"outputDir"`
	MarshalIndent bool     `yaml:"marshalIndent"`
	Verbose       bool     `yaml:"verbose"`
	Format        string   `yaml:"format"`
	MultiThreads  bool     `yaml:"multiThreads"`
	Timing        bool     `yaml:"timing"`
	Overwrite     bool     `yaml:"overwrite"`
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
	config.Format = strings.ToLower(config.Format)
	return config
}

// ParseCmd 解析命令行参数
func ParseCmd(config *Config) {
	flag.StringVar(&config.Format, "format", config.Format, "json文件的数据类型: string表示输出的数据全部是字符串格式, numeric表示输出的数据是数字类型. 使用string类型会使解析速度加快")
	flag.StringVar(&config.Format, "f", config.Format, "--format")
	flag.StringVar(&config.OutputDir, "output", config.OutputDir, "指定输出文件目录, default表示输出到源日志文件所在目录")
	flag.StringVar(&config.OutputDir, "o", config.OutputDir, "--output")
	flag.BoolVar(&config.Verbose, "verbose", config.Verbose, "输出详细信息")
	flag.BoolVar(&config.Verbose, "v", config.Verbose, "--verbose")
	flag.BoolVar(&config.Timing, "timing", config.Timing, "计算解析用时")
	flag.BoolVar(&config.Timing, "t", config.Timing, "--timing")
	flag.BoolVar(&config.MarshalIndent, "indent", config.MarshalIndent, "json文件格式化缩进, 格式化可以让json文件结构更清晰, 但解析速度会变慢, 并且文件也会变大")
	flag.BoolVar(&config.MarshalIndent, "i", config.MarshalIndent, "--indent")
	flag.BoolVar(&config.Overwrite, "overwrite", config.Overwrite, "覆盖同名json文件")
	flag.BoolVar(&config.Overwrite, "w", config.Overwrite, "--overwrite")
	flag.BoolVar(&config.MultiThreads, "multithreads", config.MultiThreads, "启用多线程来提高解析速度")
	flag.BoolVar(&config.MultiThreads, "m", config.MultiThreads, "--multithreads")
	flag.Parse()
	if flag.NArg() > 0 {
		config.SourceDir = flag.Args()
	}
}

// PrintConfig 打印Config结构体中的数据
func PrintConfig(config *Config) {
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
	_, err = file.WriteString(`# Config.yaml
# created at 2022年7月31日 by Shang
#
# 注意:
# 冒号后面一定要有空格
# 缩进要用空格
# 选项不区分大小写

# sourceDir:
# 这里指定包含日志文件的源文件夹, 程序会自动遍历该目录下的所有日志文件, 包括子目录, 注意格式:
#  - "这里输入一个目录"
#  - "可以再添加一个目录"
#  - “”
# 举例:
#  - C:\Users\shang\GolandProjects\logAnalyzer\logs
#  - D:\logs
sourceDir:
  - ./logs

# outputDir: "这里指定输出文件目录, default表示输出到源日志文件所在目录"
outputDir: default

# marshalIndent: true/false 表示json文件是否格式化缩进, 格式化可以让json文件结构更清晰, 但解析速度会变慢, 并且文件也会变大
marshalIndent: false

# verbose: true/false 是否输出详细信息
verbose: true

# format: string/numeric 输出json数据类型: string表示输出的数据全部是字符串格式, numeric表示输出的数据是数字类型. 使用string类型会使解析速度加快
format: string

# multiThreads: true/false 表示是否启用多线程来提高解析速度
multiThreads: true

# timing: true/false 是否计算解析用时
timing: true

# overwrite: true/false 当检测到已存在同名输出文件时，是否进行覆盖重新生成
overwrite: false
`)
	if err != nil {
		panic("Error writing configuration file")
	}
}

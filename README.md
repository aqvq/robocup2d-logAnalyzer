# logAnalyzer

解析Robocup2d日志文件(.rcg)生成json文件。

本项目使用go语言编写

后续可使用其他工具解析生成的json文件

# 配置文件

通过修改`config.yaml`配置文件，可修改程序默认行为。

具体内容如下：

```yaml
# Config.yaml
# created at 2022年7月27日 by Shang
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

```

# 命令行参数

输入`--help`或`-h`获取详细信息

支持的命令行参数格式有以下几种形式：

| 标志| 说明         |
| ---- |------------|
| -flag xxx | 使用空格，一个-符号 |
| --flag xxx | 使用空格，两个-符号 |
| -flag=xxx | 使用等号，一个-符号 |
| --flag=xxx | 使用等号，两个-符号 |

具体参数如下：

```bash
  -f string
        --format (default "string")
  -format string
        json文件的数据类型: string表示输出的数据全部是字符串格式, numeric表示输出的数据是数字类型. 使用string类型会使解析速度加快 (default "string")
  -i    --indent
  -indent
        json文件格式化缩进, 格式化可以让json文件结构更清晰, 但解析速度会变慢, 并且文件也会变大
  -m    --multithreads (default true)
  -multithreads
        启用多线程来提高解析速度 (default true)
  -o string
        --output (default "default")
  -output string
        指定输出文件目录, default表示输出到源日志文件所在目录 (default "default")
  -overwrite
        覆盖同名json文件
  -t    --timing (default true)
  -timing
        计算解析用时 (default true)
  -v    --verbose (default true)
  -verbose
        输出详细信息 (default true)
  -w    --overwrite

```
> 注意：
> 
> 布尔类型的参数必须使用等号的方式指定, 为`true`时可以只写参数名如`-w`
> 
> 源日志文件夹目录参数（可指定多个）必须在所有参数给定之后给出
> 

# Goland配置

安装好go之后，主要是要记住设置代理。

在`设置` -> `Go` -> `Go模块`中，添加环境变量`GOPROXY`，下拉选择`https://goproxy.cn,direct`
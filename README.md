# LogAnalyzer

解析Robocup2d日志文件(.rcg)生成json文件

本项目使用go语言编写

生成的json文件中的每一行数据都是一个结构体，后续可以使用其他工具解析并处理其中的数据

# 配置文件

通过修改`config.yaml`配置文件，可修改程序默认行为

配置文件在首次运行时会自动生成

具体内容如下：

```yaml
# config.yaml
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

具体参数内容如下：

```bash
  -d string
        --datatype (default "string")
  -datatype string
        json文件记录数据的类型. string表示以字符类型记录数据, numeric表示以数值类型记录数据. 使用string类型会使解析速度加快 (default "string")
  -f    --formatting
  -formatting
        是否格式化缩进json文件, 格式化可以让json文件结构更清晰, 但会降低解析速度并增大输出文件
  -m    --multithreading (default true)
  -multithreading
        是否启用多线程来提高解析速度 (default true)
  -o string
        --output (default "default")
  -output string
        指定json文件的输出目录, default表示输出到源文件所在目录 (default "default")
  -overwrite
        当检测到已存在同名json文件时, 是否进行覆写操作
  -t    --timing (default true)
  -timing
        是否记录解析用时 (default true)
  -v    --verbose (default true)
  -verbose
        是否输出详细信息 (default true)
  -w    --overwrite

```

> 注意：
> 
> 布尔类型的参数必须使用等号的方式指定, 为`true`时可以只写参数名, 如`-w`
> 
> 源日志文件目录参数（可指定多个）必须在所有参数给定之后给出
> 

# Goland配置

安装好go之后，主要记住设置代理。

在`设置` -> `Go` -> `Go模块`中，添加环境变量`GOPROXY`，下拉选择`https://goproxy.cn,direct`

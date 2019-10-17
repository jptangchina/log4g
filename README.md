## Log4g是什么？

Log4g是一个基于[logrus](https://github.com/sirupsen/logrus)与[lumberjack](https://github.com/natefinch/lumberjack)的日志工具，保留原有logrus的扩展功能，并支持日志输出配置。

至于为什么基于logrus，别问！问就是用star最多的库。

## Log4g特性

1. 得益于lumberjack，日志文件可以自动进行切割输出。理论上支持所有lumberjack配置
2. 采用日志文件输出时，不会在控制台打印日志，如果需要，可以通过Hook方式扩展控制台与文件同时输出。
3. 使用TOML作为配置文件
4. 不同日志级别分文件输出
5. 同一日志级别可以输出到任意多个文件中



## 配置Log4g

```toml
[[appender]]
filename = "/Users/xxxx/Desktop/logs/info.log"
maxSize = 1
maxAge = 1
maxBackups = 5
localTime = false
compress = false
minLevel = "INFO"
maxLevel = "INFO"
```

更多配置可参考`lumberjack`

如果希望在控制台输出，可以删掉或者修改appender配置。对于没有appender配置的日志级别，Log4g会使用logrus的默认输出。

## 使用Log4g

```go
package main

import log "github.com/jptangchina/log4g"

func main() {
  log.Info("Test info output")
}
```



## TODO

1. 测试文件还没写，还没仔细进行全面的测试，哈哈
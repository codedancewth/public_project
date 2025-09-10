## 公开go的脚手架

#### proto生成/使用

```shell
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

使用后执行命令

```
make publicproject 
```

即可生成对应的proto相关的请求配置文件





#### mac使用brew 安装mysql/redis

mysql

```
// 更新brew
brew update

// brew 安装mysql
brew install mysql
mysql --version

// brew启动mysql
brew services start mysql
```



redis

```
brew install redis
redis-cli --version        # 例：redis-cli 7.2.4

redis-server // 临时前台运行（调试用，Ctrl-C 即停）

brew services start redis // 后台服务方式（推荐，开机自启）

brew services list | grep redis // 检查状态
```


# Pubilc_Project

公开的go的脚手架

proto生成

go install github.com/zeromicro/go-zero/tools/goctl@latest


使用后执行命令

make publicproject 即可生成对应的proto相关的请求配置文件


mac使用brew 安装mysql/redis
brew update

brew install mysql
mysql --version

brew services start mysql


brew install redis
redis-cli --version        # 例：redis-cli 7.2.4

临时前台运行（调试用，Ctrl-C 即停）
redis-server
后台服务方式（推荐，开机自启）
brew services start redis
检查状态：
brew services list | grep redis
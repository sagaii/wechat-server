@echo off
setlocal

:: 获取本地日期和时间
for /f "tokens=2 delims==" %%I in ('wmic os get localdatetime /value') do set datetime=%%I

:: 提取年、月、日
set year=%datetime:~0,4%
set month=%datetime:~4,2%
set day=%datetime:~6,2%

:: 格式化日期为 yyyyMMdd
set formattedDate=%year%%month%%day%

:: 显示格式化后的日期
echo %formattedDate%

:: 构建 Docker 镜像
docker build -t wechat-server .

:: 运行 Docker 容器
docker run -d -p 3000:3000 wechat-server

pause

:: 登录到 Docker Hub
docker login

:: 为镜像打标签，使用格式化后的日期作为标签
docker tag wechat-server:latest sagaii/wechat-server:%formattedDate%

endlocal

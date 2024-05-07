@echo off

set PATH=%PATH%;D:\WSL\TDM-GCC-64\bin 
:: 构建前端产物
echo 安装前端依赖
call pnpm install
echo 构建前端产物
call pnpm build
:: 构建后端产物
echo 构建后端产物
go env -w GO111MODULE=on
go install -a github.com/go-bindata/go-bindata/...@latest
go install -a github.com/elazarl/go-bindata-assetfs/...@latest
cd service
go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/...
go build -o ./target/sun-panel.exe --ldflags="-X sun-panel/global.RUNCODE=release" main.go
copy ../dist ./target/web
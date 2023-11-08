## 将静态资源编译成go文件打包到可执行文件内

原版教程来源:https://blog.enianteam.com/u/sun/content/211

> 为了简化部署和减少出错的几率，将前端文件打包到可执行文件中，最终程序发布之后只有一个可执行文件

### 1. 安装
** 注意：`...`必须带上 **
```ssh
go get github.com/go-bindata/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...

# go版本>=1.17 使用intsall方式
go install -a -v github.com/go-bindata/go-bindata/...@latest
go install -a -v github.com/elazarl/go-bindata-assetfs/...@latest
```
### 2. 安装成功后将 `GOPATH/bin` 加入环境变量

参考各自系统环境变量配置即可


### 3. 压缩静态文件 到 asset目录
以下命令在Windows的`powershell`可能会报错，可使用`cmd`执行
```ssh
# 开发环境，并非真实将所有文件编译，修改静态文件可以及时生效
go-bindata-assetfs -debug -o=assets/bindata.go -pkg=assets static/... view/... # 多个
go-bindata-assetfs -debug -o=assets/bindata.go -pkg=assets assets/... 

# 正式环境，修改静态文件后需要重新编译
go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/... 
```
> 正式环境需要 去掉` -debug `

#### 参考文章
Go | Go 语言打包静态文件以及如何与Gin一起使用Go-bindata
https://www.jianshu.com/p/a7f5885679ef

[golang]Go内嵌静态资源go-bindata的安装及使用
https://www.cnblogs.com/landv/p/11577213.html
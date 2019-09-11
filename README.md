# simple_project
this is simple Golang project

Help By =>https://book.eddycjy.com/golang/

go get -u github.com/swaggo/swag/cmd/swag  报错，解决方法如下：

在go.mod文件末尾加一行：

replace github.com/urfave/cli v1.22.0 => github.com/urfave/cli v1.21.0

# goconfig
自用的读取ini文件，模仿一些*GIT项目*，精简了一些不常用的功能。

**安装**：

  ```go get github.com/y451309839/goconfig```


**使用**：
  
  ```go
  import "github.com/y451309839/gopinyin"

  ......

  cfg := goconfig.NewConfig("程序配置", "test.ini")
  
  //user := cfg.String("DEFAULT", "user")
  
  user := cfg.String("", "user")
  ```
  

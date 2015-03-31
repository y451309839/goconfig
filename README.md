goconfig
========

自用的**读取ini文件**，模仿一些*GIT项目*，因为任性，所以精简了一些自认为不常用的功能。

## 安装：

  ```go get github.com/y451309839/goconfig```

##使用：
  
  ```go
  import "github.com/y451309839/gopinyin"

  ......

  cfg := goconfig.NewConfig("程序配置", "test.ini")
  
  //user := cfg.String("DEFAULT", "user")
  
  user := cfg.String("", "user")
  ```

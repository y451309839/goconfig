# goconfig
自用的读取ini文件，精简了配置使用变量和一些不常用的获取，因为我 用 不 到！

安装：

  go get github.com/y451309839/goconfig


使用：

  cfg := goconfig.NewConfig("程序配置", "test.ini")
  
  //user := cfg.String("DEFAULT", "user")
  
  user := cfg.String("", "user")
  
  

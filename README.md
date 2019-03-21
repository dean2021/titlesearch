# TitleSearch

批量抓取域名title工具 (辅助挖洞小工具系列)

## 使用说明

    Usage of ./titlesearch:
      -c int
            Number of concurrent consumers. (default 30)
      -f string
            Input from list of domain list. (default "domains.txt")
      -o string
            Output results to file
      -t int
            http request timeout (default 10)
       
## 编译安装

    go get github.com/dean2021/titlesearch
    cd $GOPATH/src/github.com/dean2021/titlesearch
    go get .
    go build
    ./titlesearch

## License

This project is copyleft of [CSOIO](https://www.csoio.com/) and released under the GPL 3 license.


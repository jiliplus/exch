<!-- markdownlint-disable MD041 -->
<h1 align="center">:men_wrestling: exchanges</h1>
<p align="center">
<!--  -->
<a href="https://github.com/jujili/ex/releases"> <img src="https://img.shields.io/github/v/tag/jujili/ex?include_prereleases&sort=semver" alt="Release" title="Release"></a>
<!--  -->
<a href="https://www.travis-ci.org/jujili/ex"><img src="https://www.travis-ci.org/jujili/ex.svg?branch=master"/></a>
<!--  -->
<a href="https://codecov.io/gh/jujili/ex"><img src="https://codecov.io/gh/jujili/ex/branch/master/graph/badge.svg"/></a>
<!--  -->
<a href="https://goreportcard.com/report/github.com/jujili/ex"><img src="https://goreportcard.com/badge/github.com/jujili/ex" alt="Go Report Card" title="Go Report Card"/></a>
<!--  -->
<a href="http://godoc.org/github.com/jujili/ex"><img src="https://img.shields.io/badge/godoc-ta-blue.svg" alt="Go Doc" title="Go Doc"/></a>
<!--  -->
<br/>
<!--  -->
<a href="https://github.com/jujili/ex/blob/master/CHANGELOG.md"><img src="https://img.shields.io/badge/Change-Log-blueviolet.svg" alt="Change Log" title="Change Log"/></a>
<!--  -->
<a href="https://golang.google.cn"><img src="https://img.shields.io/github/go-mod/go-version/jujili/ex" alt="Go Version" title="Go Version"/></a>
<!--  -->
<a href="https://github.com/jujili/ex/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="MIT License" title="MIT License"/></a>
<!--  -->
<br/>
<!--  -->
<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=7f61280435c41608fb8cb96cf8af7d31ef0007c44b223c9e3596ce84dec329bc"><img border="0" src="https://img.shields.io/badge/QQ%20群-23%2053%2000%2093-blue.svg" alt="jili交流QQ群:23530093" title="jili交流QQ群:23530093"></a>
<!--  -->
<a href="https://mp.weixin.qq.com/s?__biz=MzA4MDU4NDI5Mw==&mid=2455230332&idx=1&sn=8086c43e259b0012596ed63d6ecd7d10&chksm=88017c76bf76f5604f2f3280ffd96029b5ccaf99db48d18066d3e3bc9bc8a2e1a05de1a3225f&mpshare=1&scene=1&srcid=&sharer_sharetime=1578553397373&sharer_shareid=5ce52651949258759d82d1bf31b455b5#rd"><img src="https://img.shields.io/badge/微信公众号-jujili-success.svg" alt="微信公众号：jujili" title="微信公众号：jujili"/></a>
<!--  -->
<a href="https://zhuanlan.zhihu.com/jujili"><img src="https://img.shields.io/badge/知乎专栏-jili-blue.svg" alt="知乎专栏：jili" title="知乎专栏：jili"/></a>
<!--  -->
</p>

ta 用于计算技术分析指标。

- [安装与更新](#%e5%ae%89%e8%a3%85%e4%b8%8e%e6%9b%b4%e6%96%b0)
- [指标简介](#%e6%8c%87%e6%a0%87%e7%ae%80%e4%bb%8b)
	- [EWMA](#ewma)

## 安装与更新

在命令行中输入以下内容，可以获取到最新版

```shell
go get -u github.com/jujili/ex
```

## 指标简介

### EWMA

[EWMA(指数移动平均，Exponentially Weighted Moving Average)](https://zh.wikipedia.org/zh-cn/%E7%A7%BB%E5%8B%95%E5%B9%B3%E5%9D%87#%E6%8C%87%E6%95%B8%E7%A7%BB%E5%8B%95%E5%B9%B3%E5%9D%87) 是一种求取时间序列均值的常用方法。其计算公式为

> S<sub>t</sub> = α * Y<sub>t</sub> + (1 - α) * S<sub>t-1</sub>

其中：

- S<sub>t</sub>   : t 时刻的平均值
- Y<sub>t</sub>   : t 时刻的实际值
- α               : 最新值的加权值
- S<sub>t-1</sub> : t-1 时刻的平均值

与简单移动平均（simple moving average，SMA）和加权移动平均（weighted moving average，WMA）相比，旧数据的权重指数级别衰减。因此，新数据的权重更高，均值更接近于实际值。

其中 α 也可以使用周期数 N 来表示：

> α = 2/(N+1)

通常 N = 5

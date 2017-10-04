[English](#scel) | [中文](#搜狗词库)
--------|-----


# scel
Sougou scel dict - 搜狗 scel 词库工具

* Scel to Protobuf
    * Easy to use in web
    * Easy to use in other language
    * No need to parse
* Useful `scel` command
    * `info` - Info about dict
    * `conv` - Conversion between format
* Provide typescript scel reader

# 搜狗词库

* 将 Scel 转换为 Protobuf 文件格式
    * 简化 web 使用
    * 简化其它语言使用
    * 不需要解析
* 非常有用的 `scel` 命令行工具
    * `info` - 词库信息
    * `conv` - 格式转换
* 提供 typescript 的 scel 解析器

## CLI

```bash
# Install
go get github.com/wenerme/scel/cmd/scel

# Info
scel info 全国省市区县地名.scel 
# file: 全国省市区县地名.scel
# name: 全国省市区县地名
# type: 单位机构名
# desc: 比搜狗自带的还全！！！
# e.g.: 澳门 重庆 福建 河北 黑龙江 江西 

# Conversion
# Convert scel to pb without `ext` data
scel conv -E 全国省市区县地名.scel city.pb
```

## JS

```bash
# Build js
yarn build
```

## Dev

```bash
# Generate pb
protoc --go_out=plugins=grpc,import_path=telattr:$HOME/go/src/ *.proto
```

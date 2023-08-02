# 飞书多维表格插件SDK
旨在让开发者调用该SDK，能快速开发飞书多维表格插件，降低开发者门槛

## 安装
```go
go get -u github.com/larksuite/base-sdk-go/v3
```

## Base Client
开发者在使用SDK前，需要先创建一个Base Client

### 创建Base Client
```go
var client = lark.NewClient("personalBaseToken", "appToken")
```

### 配置Base Client(可选)
创建Base Client时，可对创建Base Client进行一定的配置，比如设置日志级别，HTTP请求超时时间，使用飞书还是Lark品牌
```go
var client = lark.NewClient("appID", "appSecret",
    lark.WithLogLevel(larkcore.LogLevelDebug),
    lark.WithReqTimeout(3*time.Second),
    lark.WithHttpClient(http.DefaultClient))
```
每个配置选项的具体含义，如下表格:
<table>
  <thead align=left>
    <tr>
      <th>
        配置选项
      </th>
      <th>
        配置方式
      </th>
       <th>
        描述
      </th>
    </tr>
  </thead>
  <tbody align=left valign=top>
    <tr>
      <th>
        <code>LogLevel</code>
      </th>
      <td>
        <code>lark.WithLogLevel(logLevel larkcore.LogLevel)</code>
      </td>
      <td>
设置 API Client 的日志输出级别(默认为 Info 级别)，枚举值如下：

- LogLevelDebug
- LogLevelInfo
- LogLevelWarn
- LogLevelError

</td>
</tr>

<tr>
      <th>
        <code>Logger</code>
      </th>
      <td>
        <code>lark.WithLogger(logger larkcore.Logger)</code>
      </td>
      <td>
设置API Client的日志器，默认日志输出到标准输出。

开发者可通过实现下面的 Logger 接口，来设置自定义的日志器:

```go
type Logger interface {
    Debug(context.Context, ...interface{})
    Info(context.Context, ...interface{})
    Warn(context.Context, ...interface{})
    Error(context.Context, ...interface{})
}
```

</td>
</tr>

<tr>
      <th>
        <code>LogReqAtDebug</code>
      </th>
      <td>
        <code>lark.WithLogReqAtDebug(printReqRespLog bool)</code>
      </td>
      <td>
设置是否开启 Http 请求参数和响应参数的日志打印开关；开启后，在 debug 模式下会打印 http 请求和响应的 headers,body 等信息。

在排查问题时，开启该选项，有利于问题的排查。

</td>
</tr>


<tr>
      <th>
        <code>BaseUrl</code>
      </th>
      <td>
        <code>lark.WithOpenBaseUrl(baseUrl string)</code>
      </td>
      <td>
设置飞书域名，默认为FeishuBaseUrl，可用域名列表为：

```go
// 飞书域名
var FeishuBaseUrl = "https://base-api.feishu.cn"

// Lark域名
var LarkBaseUrl = "https://base-api.larksuite.com"
```

<tr>
      <th>
        <code>ReqTimeout</code>
      </th>
      <td>
        <code>lark.WithReqTimeout(time time.Duration)</code>
      </td>
      <td>
设置 SDK 内置的 Http Client 的请求超时时间，默认为0代表永不超时。
</td>
</tr>

<tr>
      <th>
        <code>HttpClient</code>
      </th>
      <td>
        <code>lark.WithHttpClient(httpClient larkcore.HttpClient)</code>
      </td>
      <td>
设置 HttpClient，用于替换 SDK 提供的默认实现。

开发者可通过实现下面的 HttpClient 接口来设置自定义的 HttpClient:

```go
type HttpClient interface {
  Do(*http.Request) (*http.Response, error)
}

```

</td>
</tr>

  </tbody>
</table>

## 基本用法
完整使用方法请参考sample目录
### 如何使用
SDK 提供了语义化的调用方式，只需要提供相关参数创建 client 实例，接着使用其上的语义化方法client.base.[资源].[方法]即可完成 API 调用。例如列出 Base 数据表记录：
```go
package main

import (
	"context"
	"fmt"
	"github.com/larksuite/base-sdk-go/v3"
	"github.com/larksuite/base-sdk-go/v3/core"
	"github.com/larksuite/base-sdk-go/v3/service/base/v1"
)

func main() {
	// 创建 Client
	// 全局baseAppToken,如果builder中有也设置了全局appToken，以build中为准
	client := lark.NewClient("personalBaseToken", "appToken")
	// 创建请求对象
	req := larkbase.NewListAppTableRecordReqBuilder().
		AppToken("bascnCMII2ORej2RItqpZZUNMIe").
		TableId("tblxI2tWaxP5dG7p").
		ViewId("vewqhz51lk").
		Filter("AND(CurrentValue.[身高]>180, CurrentValue.[体重]>150)").
		Sort("").
		FieldNames("").
		TextFieldAsArray(true).
		UserIdType("user_id").
		DisplayFormulaRef(true).
		AutomaticFields(true).
		PageToken("recn0hoyXL").
		PageSize(20).
		Build()
	// 发起请求
	resp, err := client.Base.AppTableRecord.List(context.Background(), req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
}
```

### 附件上传
```go
package main

import (
	"context"
	"fmt"
	"github.com/larksuite/base-sdk-go/v3"
	"github.com/larksuite/base-sdk-go/v3/core"
	"github.com/larksuite/base-sdk-go/v3/service/drive/v1"
	"os"
)

func main() {
	// 创建 Client
	// 全局baseAppToken,如果builder中有也设置了全局appToken，以build中为准
	client := lark.NewClient("personalBaseToken", "appToken")
	file, err := os.Open("filepath")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 创建请求对象
	req := larkdrive.NewUploadAllMediaReqBuilder().
		Body(larkdrive.NewUploadAllMediaReqBodyBuilder().
			FileName("demo.jpeg").
			ParentType("doc_image").
			ParentNode("doccnFivLCfJfblZjGZtxgabcef").
			Size(1024).
			Checksum("12345678").
			Extra("").
			File(file).
			Build()).
		Build()
	// 发起请求
	resp, err := client.Drive.Media.UploadAll(context.Background(), req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
}
```

### 附件下载
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/larksuite/base-sdk-go/v3"
	"github.com/larksuite/base-sdk-go/v3/core"
	"github.com/larksuite/base-sdk-go/v3/service/drive/v1"
)

// 如果附件开启高级权限，需要配置Extra字段。
// Attachments格式：第一个key为字段id，第二个key为记录id，value为文件token数组。详见：
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/drive-v1/media/introduction
type Extra struct {
	BitablePerm struct {
		TableId     string                         `json:"tableId"`
		Attachments map[string]map[string][]string `json:"attachments"`
	} `json:"bitablePerm"`
}

func main() {
	// 创建 Client
	// 全局baseAppToken,如果builder中有也设置了全局appToken，以build中为准
	client := lark.NewClient("personalBaseToken", "appToken")

	// 如果开启高级权限，设置extra
	extra := &Extra{
		BitablePerm: struct {
			TableId     string                         `json:"tableId"`
			Attachments map[string]map[string][]string `json:"attachments"`
		}(struct {
			TableId     string
			Attachments map[string]map[string][]string
		}{TableId: "table_id", Attachments: map[string]map[string][]string{
			"field_id": {
				"record_id": {"file_token"},
			},
		}}),
	}
	bs, err := json.Marshal(extra)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建请求对象
	req := larkdrive.NewDownloadMediaReqBuilder().
		FileToken("appToken").
		Extra(string(bs)).
		Build()
	// 发起请求
	resp, err := client.Drive.Media.Download(context.Background(), req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
}
```
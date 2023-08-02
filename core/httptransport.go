/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package larkcore

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var reqTranslator ReqTranslator

func NewHttpClient(config *Config) {
	if config.HttpClient == nil {
		if config.ReqTimeout == 0 {
			config.HttpClient = http.DefaultClient
		} else {
			config.HttpClient = &http.Client{Timeout: config.ReqTimeout}
		}
	}
}

func NewSerialization(config *Config) {
	if config.Serializable == nil {
		config.Serializable = &DefaultSerialization{}
	}
}

func determineTokenType(option *RequestOption) AccessTokenType {
	return AccessTokenTypePersonal
}

func validate(config *Config, option *RequestOption, accessTokenType AccessTokenType) error {

	if option.Header != nil {
		if option.Header.Get(HttpHeaderKeyRequestId) != "" {
			return &IllegalParamError{msg: fmt.Sprintf("use %s as header key is not allowed", HttpHeaderKeyRequestId)}
		}
		if option.Header.Get(httpHeaderRequestId) != "" {
			return &IllegalParamError{msg: fmt.Sprintf("use %s as header key is not allowed", httpHeaderRequestId)}
		}
		if option.Header.Get(HttpHeaderKeyLogId) != "" {
			return &IllegalParamError{msg: fmt.Sprintf("use %s as header key is not allowed", HttpHeaderKeyLogId)}
		}
	}

	return nil
}

func doSend(ctx context.Context, rawRequest *http.Request, httpClient HttpClient, logger Logger) (*ApiResp, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(rawRequest)
	if err != nil {
		if er, ok := err.(*url.Error); ok {
			if er.Timeout() {
				return nil, &ClientTimeoutError{msg: er.Error()}
			}

			if e, ok := er.Err.(*net.OpError); ok && e.Op == "dial" {
				return nil, &DialFailedError{msg: er.Error()}
			}
		}
		return nil, err
	}

	if resp.StatusCode == http.StatusGatewayTimeout {
		logID := resp.Header.Get(HttpHeaderKeyLogId)
		if logID == "" {
			logID = resp.Header.Get(HttpHeaderKeyRequestId)
		}
		logger.Info(ctx, fmt.Sprintf("req path:%s, server time out,requestId:%s",
			rawRequest.URL.RequestURI(), logID))
		return nil, &ServerTimeoutError{msg: "server time out error"}
	}
	body, err := readResponse(resp)
	if err != nil {
		return nil, err
	}

	return &ApiResp{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		RawBody:    body,
	}, nil
}

func Request(ctx context.Context, req *ApiReq, config *Config, options ...RequestOptionFunc) (*ApiResp, error) {
	option := &RequestOption{}
	for _, optionFunc := range options {
		optionFunc(option)
	}

	accessTokenType := determineTokenType(option)
	err := validate(config, option, accessTokenType)
	if err != nil {
		return nil, err
	}

	return doRequest(ctx, req, accessTokenType, config, option)

}

func doRequest(ctx context.Context, httpReq *ApiReq, accessTokenType AccessTokenType, config *Config, option *RequestOption) (*ApiResp, error) {
	var rawResp *ApiResp
	var errResult error
	for i := 0; i < 2; i++ {
		req, err := reqTranslator.translate(ctx, httpReq, accessTokenType, config, option)
		if err != nil {
			return nil, err
		}

		if config.LogReqAtDebug {
			config.Logger.Debug(ctx, fmt.Sprintf("req:%v", req))
		} else {
			config.Logger.Debug(ctx, fmt.Sprintf("req:%s,%s", httpReq.HttpMethod, httpReq.ApiPath))
		}
		rawResp, err = doSend(ctx, req, config.HttpClient, config.Logger)
		if config.LogReqAtDebug {
			config.Logger.Debug(ctx, fmt.Sprintf("resp:%v", rawResp))
		}
		_, isDialError := err.(*DialFailedError)
		if err != nil && !isDialError {
			return nil, err
		}
		errResult = err
		if isDialError {
			continue
		}

		fileDownloadSuccess := option.FileDownload && rawResp.StatusCode == http.StatusOK
		if fileDownloadSuccess || !strings.Contains(rawResp.Header.Get(contentTypeHeader), contentTypeJson) {
			break
		}

		codeError := &CodeError{}
		err = config.Serializable.Deserialize(rawResp.RawBody, codeError)
		if err != nil {
			return nil, err
		}
		break
	}

	if errResult != nil {
		return nil, errResult
	}
	return rawResp, nil
}

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
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewHttpClient(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test_new_http_client",
			args: args{
				config: &Config{
					ReqTimeout: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewHttpClient(tt.args.config)
		})
	}
}

func TestNewSerialization(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test_new_serialization",
			args: args{
				config: &Config{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewSerialization(tt.args.config)
		})
	}
}

type MockHttpClient struct {
}

func (client *MockHttpClient) Do(*http.Request) (*http.Response, error) {
	return &http.Response{
		Status:           "200",
		StatusCode:       200,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             io.NopCloser(strings.NewReader("Hello, world!")),
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}, nil
}

func TestRequest(t *testing.T) {

	mockClient := MockHttpClient{}

	type args struct {
		ctx     context.Context
		req     *ApiReq
		config  *Config
		options []RequestOptionFunc
	}
	tests := []struct {
		name    string
		args    args
		want    *ApiResp
		wantErr bool
	}{
		{
			name: "test_request",
			args: args{
				ctx: context.Background(),
				req: &ApiReq{
					HttpMethod:                "POST",
					ApiPath:                   "",
					Body:                      nil,
					QueryParams:               QueryParams{},
					PathParams:                PathParams{},
					SupportedAccessTokenTypes: []AccessTokenType{AccessTokenTypePersonal},
				},
				config: &Config{
					Logger: newLoggerProxy(LogLevelInfo, defaultLogger{
						logger: log.New(os.Stdout, "", log.LstdFlags),
					}),
					HttpClient: &mockClient,
				},
			},
			want: &ApiResp{
				StatusCode: 200,
				RawBody:    []byte("Hello, world!"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Request(tt.args.ctx, tt.args.req, tt.args.config, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_determineTokenType(t *testing.T) {
	type args struct {
		option *RequestOption
	}
	tests := []struct {
		name string
		args args
		want AccessTokenType
	}{
		{
			name: "test_determineTokenType",
			args: args{},
			want: AccessTokenTypePersonal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineTokenType(tt.args.option); got != tt.want {
				t.Errorf("determineTokenType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_doRequest(t *testing.T) {
	mockClient := MockHttpClient{}
	type args struct {
		ctx             context.Context
		httpReq         *ApiReq
		accessTokenType AccessTokenType
		config          *Config
		option          *RequestOption
	}
	tests := []struct {
		name    string
		args    args
		want    *ApiResp
		wantErr bool
	}{
		{
			name: "test_doRequest",
			args: args{
				ctx:             context.Background(),
				httpReq:         &ApiReq{},
				accessTokenType: AccessTokenTypePersonal,
				config: &Config{
					HttpClient: &mockClient,
					Logger: newLoggerProxy(LogLevelInfo, defaultLogger{
						logger: log.New(os.Stdout, "", log.LstdFlags),
					}),
				},
				option: &RequestOption{},
			},
			want: &ApiResp{
				StatusCode: 200,
				RawBody:    []byte("Hello, world!"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := doRequest(tt.args.ctx, tt.args.httpReq, tt.args.accessTokenType, tt.args.config, tt.args.option)
			if (err != nil) != tt.wantErr {
				t.Errorf("doRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_doSend(t *testing.T) {

	mockClient := MockHttpClient{}

	type args struct {
		ctx        context.Context
		rawRequest *http.Request
		httpClient HttpClient
		logger     Logger
	}
	tests := []struct {
		name    string
		args    args
		want    *ApiResp
		wantErr bool
	}{
		{
			name: "test_do_send",
			args: args{
				httpClient: &mockClient,
			},
			want: &ApiResp{
				StatusCode: 200,
				RawBody:    []byte("Hello, world!"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := doSend(tt.args.ctx, tt.args.rawRequest, tt.args.httpClient, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("doSend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doSend() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	type args struct {
		config          *Config
		option          *RequestOption
		accessTokenType AccessTokenType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test_validate",
			args: args{
				option: &RequestOption{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.args.config, tt.args.option, tt.args.accessTokenType); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

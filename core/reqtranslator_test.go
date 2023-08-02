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
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestFormdata_AddField(t *testing.T) {
	type fields struct {
		fields map[string]interface{}
		data   *struct {
			content     []byte
			contentType string
		}
	}
	type args struct {
		field string
		val   interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Formdata
	}{
		{
			name: "test_add_field",
			fields: fields{
				fields: map[string]interface{}{
					"test": nil,
				},
			},
			args: args{
				field: "test",
				val:   nil,
			},
			want: &Formdata{
				fields: map[string]interface{}{
					"test": nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fd := &Formdata{
				fields: tt.fields.fields,
				data:   tt.fields.data,
			}
			if got := fd.AddField(tt.args.field, tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormdata_AddFile(t *testing.T) {
	type fields struct {
		fields map[string]interface{}
		data   *struct {
			content     []byte
			contentType string
		}
	}
	type args struct {
		field string
		r     io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Formdata
	}{
		{
			name: "test_add_file",
			fields: fields{
				fields: map[string]interface{}{},
			},
			args: args{},
			want: &Formdata{
				fields: map[string]interface{}{
					"": nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fd := &Formdata{
				fields: tt.fields.fields,
				data:   tt.fields.data,
			}
			if got := fd.AddFile(tt.args.field, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFormdata(t *testing.T) {
	tests := []struct {
		name string
		want *Formdata
	}{
		{
			name: "test_new_form_data",
			want: &Formdata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFormdata(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFormdata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReqTranslator_getFullReqUrl(t *testing.T) {
	type args struct {
		domain   string
		httpPath string
		pathVars map[string]interface{}
		queries  map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test_fet_full_req_url",
			args: args{
				domain:   "test",
				httpPath: "test1",
			},
			want:    "testtest1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := &ReqTranslator{}
			got, err := translator.getFullReqUrl(tt.args.domain, tt.args.httpPath, tt.args.pathVars, tt.args.queries)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFullReqUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFullReqUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReqTranslator_newHTTPRequest(t *testing.T) {

	request, _ := http.NewRequestWithContext(context.Background(), "POST", "test", bytes.NewBuffer([]byte("test")))
	request.Header.Set(userAgentHeader, userAgent())
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ""))

	type args struct {
		ctx             context.Context
		httpMethod      string
		url             string
		contentType     string
		body            []byte
		accessTokenType AccessTokenType
		option          *RequestOption
		config          *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name: "test_new_http_request",
			args: args{
				ctx:             context.Background(),
				httpMethod:      "POST",
				url:             "test",
				accessTokenType: AccessTokenTypePersonal,
				option:          &RequestOption{},
				config:          &Config{},
				body:            []byte("test"),
			},
			want:    request,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := &ReqTranslator{}
			got, err := translator.newHTTPRequest(tt.args.ctx, tt.args.httpMethod, tt.args.url, tt.args.contentType, tt.args.body, tt.args.accessTokenType, tt.args.option, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("newHTTPRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 比较header
			if !reflect.DeepEqual(got.Header, tt.want.Header) {
				t.Errorf("newHTTPRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReqTranslator_parseInput(t *testing.T) {
	type args struct {
		input  interface{}
		option *RequestOption
	}
	tests := []struct {
		name  string
		args  args
		want  map[string]interface{}
		want1 map[string]interface{}
		want2 interface{}
	}{
		{
			name: "test_nil",
			args: args{
				input: nil,
			},
			want:  nil,
			want1: nil,
			want2: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translator := &ReqTranslator{}
			got, got1, got2 := translator.parseInput(tt.args.input, tt.args.option)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseInput() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseInput() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("parseInput() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

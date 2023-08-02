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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiResp_JSONUnmarshalBody(t *testing.T) {

	type Value struct {
		val string `json:"val"`
	}
	str := &Value{
		val: "test",
	}
	br, _ := json.Marshal(str)

	var result Value

	type fields struct {
		StatusCode int
		Header     http.Header
		RawBody    []byte
	}
	type args struct {
		val    interface{}
		config *Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test_json_unmarshal_body",
			fields: fields{
				StatusCode: 200,
				Header:     http.Header{contentTypeHeader: []string{contentTypeJson}},
				RawBody:    br,
			},
			args: args{
				val: &result,
				config: &Config{
					Serializable: &DefaultSerialization{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ApiResp{
				StatusCode: tt.fields.StatusCode,
				Header:     tt.fields.Header,
				RawBody:    tt.fields.RawBody,
			}
			if err := resp.JSONUnmarshalBody(tt.args.val, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("JSONUnmarshalBody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApiResp_RequestId(t *testing.T) {
	type fields struct {
		StatusCode int
		Header     http.Header
		RawBody    []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test_get_request_id",
			fields: fields{
				StatusCode: 200,
				Header: map[string][]string{
					HttpHeaderKeyLogId: []string{"12345"},
				},
			},
			want: "12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ApiResp{
				StatusCode: tt.fields.StatusCode,
				Header:     tt.fields.Header,
				RawBody:    tt.fields.RawBody,
			}
			if got := resp.RequestId(); got != tt.want {
				t.Errorf("RequestId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApiResp_String(t *testing.T) {
	type fields struct {
		StatusCode int
		Header     http.Header
		RawBody    []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test_string",
			fields: fields{
				StatusCode: 200,
				Header: map[string][]string{
					contentTypeHeader: []string{"123"},
				},
			},
			want: "StatusCode: 200, Header:map[Content-Type:[123]], Content-Type: 123, Body: <binary> len 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ApiResp{
				StatusCode: tt.fields.StatusCode,
				Header:     tt.fields.Header,
				RawBody:    tt.fields.RawBody,
			}
			if got := resp.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApiResp_Write(t *testing.T) {

	type fields struct {
		StatusCode int
		Header     http.Header
		RawBody    []byte
	}
	type args struct {
		writer http.ResponseWriter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test_write",
			fields: fields{
				StatusCode: 200,
				Header: map[string][]string{
					contentTypeHeader: {"123"},
				},
			},
			args: args{
				writer: httptest.NewRecorder(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ApiResp{
				StatusCode: tt.fields.StatusCode,
				Header:     tt.fields.Header,
				RawBody:    tt.fields.RawBody,
			}
			resp.Write(tt.args.writer)
		})
	}
}

func TestCodeError_Error(t *testing.T) {
	type fields struct {
		Code int
		Msg  string
		Err  *struct {
			Details              []*CodeErrorDetail              `json:"details,omitempty"`
			PermissionViolations []*CodeErrorPermissionViolation `json:"permission_violations,omitempty"`
			FieldViolations      []*CodeErrorFieldViolation      `json:"field_violations,omitempty"`
		}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test_code_error",
			fields: fields{
				Code: 200,
				Msg:  "error",
			},
			want: "msg:error,code:200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := CodeError{
				Code: tt.fields.Code,
				Msg:  tt.fields.Msg,
				Err:  tt.fields.Err,
			}
			if got := ce.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCodeError_String(t *testing.T) {
	type fields struct {
		Code int
		Msg  string
		Err  *struct {
			Details              []*CodeErrorDetail              `json:"details,omitempty"`
			PermissionViolations []*CodeErrorPermissionViolation `json:"permission_violations,omitempty"`
			FieldViolations      []*CodeErrorFieldViolation      `json:"field_violations,omitempty"`
		}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test_code_error",
			fields: fields{
				Code: 200,
				Msg:  "error",
			},
			want: "msg:error,code:200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ce := CodeError{
				Code: tt.fields.Code,
				Msg:  tt.fields.Msg,
				Err:  tt.fields.Err,
			}
			if got := ce.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileNameByHeader(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test_file_name",
			args: args{
				header: map[string][]string{
					"Content-Disposition": []string{},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileNameByHeader(tt.args.header); got != tt.want {
				t.Errorf("FileNameByHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

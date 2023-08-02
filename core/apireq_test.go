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
	"net/url"
	"testing"
)

func TestPathParams_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		u    PathParams
		args args
		want string
	}{
		{
			name: "test_get_path_param_key",
			args: args{
				key: "test_key",
			},
			u: PathParams{
				"test_key": "test_value",
			},
			want: "test_value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Get(tt.args.key); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathParams_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		u    PathParams
		args args
	}{
		{
			name: "test_set_path_param_key",
			args: args{
				key: "test_key",
			},
			u: PathParams{
				"test_key": "test_value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.Set(tt.args.key, tt.args.value)
		})
	}
}

func TestQueryParams_Add(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		u    QueryParams
		args args
	}{
		{
			name: "test_add_path_param_key",
			args: args{
				key: "test_key",
			},
			u: QueryParams{
				"test_key": []string{"test_value"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.Add(tt.args.key, tt.args.value)
		})
	}
}

func TestQueryParams_Encode(t *testing.T) {
	tests := []struct {
		name string
		u    QueryParams
		want string
	}{
		{
			name: "test_query_param_encode",
			u: QueryParams{
				"test_key": []string{"test_value"},
			},
			want: url.Values(map[string][]string{
				"test_key": []string{"test_value"},
			}).Encode(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Encode(); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryParams_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		u    QueryParams
		args args
		want string
	}{
		{
			name: "test_query_params_get",
			u: QueryParams{
				"test_key": []string{"test_value"},
			},
			args: args{
				key: "test_key",
			},
			want: "test_value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.Get(tt.args.key); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryParams_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		u    QueryParams
		args args
	}{
		{
			name: "test_query_params_set",
			u: QueryParams{
				"test_key": []string{"test_value"},
			},
			args: args{
				key: "test_key",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.Set(tt.args.key, tt.args.value)
		})
	}
}

package main

import (
	"encoding/json"
	"testing"
)

type Scenario struct {
	Name       string
	Param      string
	Want       bool
	ErrMessage string
}

func TestLoadConfig(t *testing.T) {
	config, err := loadConfig("conf.json")
	if err == nil {
		t.Fatalf(err.Error())
	}

	config, err = loadConfig("config.json")
	if err != nil {
		t.Fatalf(err.Error())
	}

	// check success config
	if len(config.Services) == 0 {
		t.Fatalf(err.Error())
	}
}

func TestInStringSlice(t *testing.T) {
	type Scenario struct {
		Name  string
		Param string
		Want  bool
	}
	tests := []Scenario{
		{Name: "Success Case", Param: "c", Want: true},
		{Name: "Error Case", Param: "d", Want: false},
	}
	s := []string{"a", "b", "c"}
	for _, tc := range tests {
		r := inStringSlice(tc.Param, s)
		if r != tc.Want {
			t.Fatalf(`%v: inStringSlice(c, %v) result=%v, want=%v`, tc.Name, s, r, tc.Want)
		}
	}
}

func TestGetFilePath(t *testing.T) {
	type Scenario struct {
		Name        string
		Want        string
		BaseDir     string
		FilePrefix  string
		Separator   string
		ServiceName string
	}
	tests := []Scenario{
		{Name: "test1", Want: "/base/dir/krakend-service.json", BaseDir: "/base/dir/", FilePrefix: "krakend", Separator: "-", ServiceName: "service"},
		{Name: "test2", Want: "/base/dir/service.json", BaseDir: "/base/dir/", FilePrefix: "", Separator: "", ServiceName: "service"},
	}
	for _, tc := range tests {
		r := getFilePath(tc.BaseDir, tc.FilePrefix, tc.Separator, tc.ServiceName)
		if r != tc.Want {
			t.Fatalf(`%v: getFilePath: result=%v, want=%v`, tc.Name, r, tc.Want)
		}
	}

}

/*
func TestDoRequest(t *testing.T) {
	var response interface{}
	err := doRequest("http://test.url", &response)
	if err == nil {
		t.Fatalf(err.Error())
	}

	err = doRequest("https://petstore.swagger.io/v2/swagger.json", &response)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if response == nil {
		t.Fatalf("Couldn't set response")
	}
}
*/

func TestTransformSwaggerResponse(t *testing.T) {
	configPayload := []byte(`{
     "name": "test-service",
     "internal_address": "http://localhost:8080",
     "output_encoding": "no-op",
     "backend_encoding": "no-op",
	}`)
	serviceConfig := Service{}
	json.Unmarshal(configPayload, &serviceConfig)

	type Want struct {
		Urls []string
	}
	type Scenario struct {
		Name           string
		Want           Want
		SwaggerPayload []byte
	}

	tests := []Scenario{
		{
			Name:           "test1",
			Want:           Want{Urls: []string{"/user/{username}", "/user/login", "/user/logout"}},
			SwaggerPayload: []byte(`{ "paths": { "/user/{username}": { "get": { "parameters": [ { "name": "username", "in": "path" } ] }, "put": { "parameters": [ { "name": "username", "in": "path" }, { "in": "body", "name": "body" } ] }, "delete": { "parameters": [ { "name": "username", "in": "path" } ] } }, "/user/login": { "get": { "parameters": [ { "name": "username", "in": "query" }, { "name": "password", "in": "query" } ] } }, "/user/logout": { "get": { "tags": [ "user" ], "parameters": [] } } } }`),
		},
	}
	for _, tc := range tests {
		swgResponse := SwaggerResponse{}
		json.Unmarshal(tc.SwaggerPayload, &swgResponse)
		paths := transformSwaggerResponse(serviceConfig, swgResponse)
		for _, p := range paths {
			if !inStringSlice(p.Endpoint, tc.Want.Urls) {
				t.Fatalf(`%v: transformSwaggerResponse: result=%v, want=%v`, tc.Name, p.Endpoint, tc.Want.Urls)
			}
		}

	}

}

func TestCrossMapping(t *testing.T) {
	payloadOne := []byte(`    {
      "name": "test-service",
      "output_encoding": "no-op",
      "backend_encoding": "no-op",
      "cross_mapping": [
        {"from": "/<service_name>/swagger", "to": "/swagger", "method": "GET", "internal_address":  "http://localhost:8080", "headers_to_pass": []}
      ]
    }`)

	payloadTwo := []byte(`    {
      "name": "test-service",
      "output_encoding": "no-op",
      "backend_encoding": "no-op",
      "cross_mapping": [
        {"from": "/swagger", "to": "/swagger", "method": "GET", "internal_address":  "http://localhost:8080", "headers_to_pass": []}
      ]
    }`)

	type Scenario struct {
		Name    string
		Payload []byte
		Want    string
	}
	tests := []Scenario{
		{Name: "test1", Payload: payloadOne, Want: "/test-service/swagger"},
		{Name: "test2", Payload: payloadTwo, Want: "/swagger"},
	}

	for _, tc := range tests {
		serviceConfig := Service{}
		json.Unmarshal(tc.Payload, &serviceConfig)
		paths := crossMapping(serviceConfig)
		r := paths[0].Endpoint
		if r != tc.Want {
			t.Fatalf(`%v: crossMapping: result=%v, want=%v`, tc.Name, r, tc.Want)

		}
	}
}

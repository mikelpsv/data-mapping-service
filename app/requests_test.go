package app

import (
	"net/http"
	"testing"
)

func TestGetSimpleValue(t *testing.T) {
	request, _ := http.NewRequest(
		"GET",
		"http://localhost:8989?test_param=93",
		nil)
	val, exist := GetSimpleValue(request, "test_param")
	if !exist {
		t.Error(!exist, "param `test_param` not found")
	}
	if val != "93" {
		t.Error(val == "93")
	}

}

func TestGetSimpleValueAsInt(t *testing.T) {
	request, _ := http.NewRequest(
		"GET",
		"http://localhost:8989?test_param=93",
		nil)
	val, exist, err := GetSimpleValueAsInt(request, "test_param")
	if !exist {
		t.Error(!exist, "param `test_param` not found")
	}
	if val != 93 {
		t.Error(val == 93, "value is incorrect")
	}
	if err != nil {
		t.Error(err != nil, err)
	}

	request, _ = http.NewRequest(
		"GET",
		"http://localhost:8989?test_param=string_value",
		nil)
	val, exist, err = GetSimpleValueAsInt(request, "test_param")
	if val != 0 {
		t.Error(val != 0)
	}
	if !exist {
		t.Error(!exist, "param `test_param` not found")
	}
	if err == nil {
		t.Error(err == nil)
	}

}

package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestGetLimitOffset(t *testing.T) {
	var limitIn int64 = 64
	var offsetIn int64 = 128

	val := url.Values{}
	val.Add("limit", strconv.FormatInt(limitIn, 10))
	val.Add("offset", strconv.FormatInt(offsetIn, 10))

	u, err := url.Parse(fmt.Sprintf("https://example.com/test?%s", val.Encode()))
	if err != nil {
		t.Errorf("Error while parsing url: %s", err)
	}

	limitOut, offsetOut := getLimitOffset(&http.Request{URL: u})
	if limitOut != limitIn {
		t.Errorf("Getting limit and offset failed: limit \"%v\" and \n%v\n unequal", limitOut, limitIn)
	}
	if offsetOut != offsetIn {
		t.Errorf("Getting limit and offset failed: offset \"%v\" and \n%v\n unequal", offsetOut, offsetIn)
	}
}

func TestGetSort(t *testing.T) {
	var defaultSort = "-test"

	val := url.Values{}
	val.Add("sort", "testSort")

	u, err := url.Parse(fmt.Sprintf("https://example.com/test?%s", val.Encode()))
	if err != nil {
		t.Errorf("Error while parsing url: %s", err)
	}

	sort := getSort(defaultSort, &http.Request{URL: u})
	if v, ok := sort["testSort"]; ok {
		if v != 1 {
			t.Errorf("Getting sort failed: value \"1\" expected but \"%v\" received", v)
		}
	} else {
		t.Error("Getting sort failed: expected key \"testSort\" doesn't exist")
	}

	val.Set("sort", "-testSort")

	u, err = url.Parse(fmt.Sprintf("https://example.com/test?%s", val.Encode()))
	if err != nil {
		t.Errorf("Error while parsing url: %s", err)
	}

	sort = getSort(defaultSort, &http.Request{URL: u})
	if v, ok := sort["testSort"]; ok {
		if v != -1 {
			t.Errorf("Getting sort failed: value \"-1\" expected but \"%v\" received", v)
		}
	} else {
		t.Error("Getting sort failed: expected key \"testSort\" doesn't exist")
	}
}

func TestIsSupportedMediaType(t *testing.T) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")

	if !isSupportedMediaType(&http.Request{Header: header}) {
		t.Error("Checking supported media type failed: valid media type invalid")
	}

	header.Set("Content-Type", "application/xml")

	if isSupportedMediaType(&http.Request{Header: header}) {
		t.Error("Checking supported media type failed: invalid media type valid")
	}
}

type testObject struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func TestParseJSONBody(t *testing.T) {
	testIn := testObject{}

	b, err := json.Marshal(testIn)
	if err != nil {
		t.Errorf("Error during encoding JSON: %s", err)
	}

	testOut := testObject{}

	err = parseJSONBody(ioutil.NopCloser(bytes.NewReader(b)), &testOut)
	if err != nil {
		t.Errorf("Parsing JSON body failed: %s", err)
	}

	if !reflect.DeepEqual(testOut, testIn) {
		t.Errorf("Parsing JSON body failed: object %v and %v unequal", testOut, testIn)
	}
}

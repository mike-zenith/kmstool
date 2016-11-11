package kmstool

import (
	"reflect"
	"strings"
	"testing"
)

var content = `
parameters:
	secret: kms:topsecret
`

func TestNewKMSItem(t *testing.T) {
	i := NewKMSItem(content)
	if i.c != content {
		t.Error("Content was not set")
	}
}

func TestKMSItem_Content(t *testing.T) {
	i := NewKMSItem(content)
	if i.Content() != content {
		t.Fail()
	}
}

func TestKMSItem_Parse(t *testing.T) {
	var expected = map[string]string{
		"secret": "kms:topsecret",
	}
	i := NewKMSItem(content)
	i.Parse(func(content string) (map[string]string, error) {
		return expected, nil
	})

	if !reflect.DeepEqual(expected, i.m) {
		t.Error("Failed asserting that parser run")
	}
}

func TestKMSItem_Replace(t *testing.T) {
	expected := "Contents.FOO.BAR"

	i := NewKMSItem("Contents.1.2")
	i.m = map[string]string{
		"1": "1",
		"2": "2",
	}

	replaceMap := map[string]string{
		"1": "FOO",
		"2": "BAR",
	}

	i.Replace(func(content string, hit string, key string) (string, error) {
		return strings.Replace(content, hit, replaceMap[hit], -1), nil
	})

	if !reflect.DeepEqual(expected, i.c) {
		t.Error("Failed asserting that replacer run")
	}
}

func TestBulkFile_RunOnParsed(t *testing.T) {
	parsed := map[string]string{
		"foo": "oof",
		"bar": "rab",
	}

	callMap := map[string]string{
		"foo": "bar",
		"bar": "foo",
	}

	i := NewKMSItem(content)
	i.m = parsed

	i.RunOnParsed(func(key string, content string) (string, error) {
		return callMap[key], nil
	})

	if !reflect.DeepEqual(i.m, callMap) {
		t.Error("Failed assertint that func run on parsed elements")
	}
}

package kmstool

import (
	"reflect"
	"testing"
)

var baseReg = `(?P<hit>kms:(?P<key>.+))(?:\n|"|'|$)`

var successDataProvider = map[string]struct {
	reg      string
	in       string
	expected map[string]string
}{
	"simple string with hit": {
		reg: baseReg,
		in:  `kms:mysql_staging_in`,
		expected: map[string]string{
			"kms:mysql_staging_in": "mysql_staging_in",
		},
	},
	"multiple string hit": {
		reg: baseReg,
		in: `
		parameters:
			mysql_password: kms:mysql_staging_password
			app_secret: kms:app_secret
			`,
		expected: map[string]string{
			"kms:mysql_staging_password": "mysql_staging_password",
			"kms:app_secret":             "app_secret",
		},
	},
	"overwritten multiple string hit": {
		reg: baseReg,
		in: `
		parameters:
			mysql_password: kms:mysql_password
			another_password: kms:mysql_password
			# comment: kms:mysql_password

		`,
		expected: map[string]string{
			"kms:mysql_password": "mysql_password",
		},
	},
}

func TestNewRegexpStringParser_success(t *testing.T) {
	for key, m := range successDataProvider {

		parser := NewRegexpStringParser(m.reg)
		result, err := parser(m.in)

		if err != nil {
			t.Errorf("DataSet '%s': File parser error: %s", key, err)
		}
		if !reflect.DeepEqual(result, m.expected) {
			t.Errorf("Dataset '%s': Result does not match %s , %s", key, result, m.expected)
		}
	}
}

var errorDataProvider = map[string]struct {
	reg      string
	in       string
	expected bool
}{
	"no hit single result, empty map": {
		reg:      baseReg,
		in:       `kms-mysql_staging_in`,
		expected: false,
	},
	"reg match but no `hit` found": {
		reg:      `(?P<pikachu>kms:(?P<key>.+))(?:\n|"|'|$)`,
		in:       `kms:mysql_password`,
		expected: true,
	},
	"reg match but no `key` found": {
		reg:      `(?P<hit>kms:(?P<pikachu>.+))(?:\n|"|'|$)`,
		in:       `kms:mysql_password`,
		expected: true,
	},
}

func TestNewRegexpStringParser_returnsError(t *testing.T) {
	for key, m := range errorDataProvider {

		parser := NewRegexpStringParser(m.reg)
		_, err := parser(m.in)

		if (err == nil && m.expected) || (err != nil && !m.expected) {
			t.Errorf(
				"DataSet '%s': Failed asserting that expected error returned. Expected: %b, Got: %v",
				key,
				m.expected,
				err,
			)
		}
	}
}

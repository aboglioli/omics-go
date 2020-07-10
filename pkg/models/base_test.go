package models

import (
	"strconv"
	"testing"
)

func TestIDToStr(t *testing.T) {
	tests := []struct {
		id  ID
		res string
	}{{
		25,
		"25",
	}, {
		-2,
		"-2",
	}, {
		0,
		"0",
	}, {
		789,
		"789",
	}}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res := test.id.String()

			if res != test.res {
				t.Errorf("\nExp:%v\nAct:%v", test.res, res)
			}
		})
	}
}

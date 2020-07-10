package errors

import (
	"reflect"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		name     string
		create   func() error
		expected error
	}{{
		"default",
		func() error {
			return New(INTERNAL, "CODE")
		},
		Error{
			kind:    INTERNAL,
			code:    "CODE",
			context: make(Context),
		},
	}, {
		"set status",
		func() error {
			return New(APPLICATION, "CODE").Status(401)
		},
		Error{
			kind:    APPLICATION,
			code:    "CODE",
			status:  401,
			context: make(Context),
		},
	}, {
		"set message",
		func() error {
			return New(APPLICATION, "CODE").Message("%d;%.2f;%s", 1, 1.2, "Hi")
		},
		Error{
			kind:    APPLICATION,
			code:    "CODE",
			message: "1;1.20;Hi",
			context: make(Context),
		},
	}, {
		"context",
		func() error {
			return New(APPLICATION, "CODE").
				Context(Context{
					"Prop1": "Value1",
					"Prop2": 2,
				}).
				AddContext("Prop3", []int{5, 6})
		},
		Error{
			kind: APPLICATION,
			code: "CODE",
			context: Context{
				"Prop1": "Value1",
				"Prop2": 2,
				"Prop3": []int{5, 6},
			},
		},
	}, {
		"custom",
		func() error {
			return New(APPLICATION, "custom.code").
				AddContext("Prop1", "Value1").
				AddContext("Prop2", 2.5).
				Context(Context{
					"Prop3": []int{5, 6},
				}).
				AddContext("Prop1", "Changed1")
		},
		Error{
			kind: APPLICATION,
			code: "custom.code",
			context: map[string]interface{}{
				"Prop1": "Changed1",
				"Prop2": 2.5,
				"Prop3": []int{5, 6},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.create()
			if !reflect.DeepEqual(err, test.expected) {
				t.Errorf("%s\n[EXP]: %v\n[ACT]: %v", test.name, test.expected, err)
			}
		})
	}
}

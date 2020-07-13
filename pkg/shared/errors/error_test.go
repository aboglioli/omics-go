package errors

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		name string
		err1 Error
		err2 Error
	}{{
		"default",
		Internal.New().Code("CODE"),
		Error{
			kind:    Internal,
			code:    "CODE",
			context: make(Context),
		},
	}, {
		"set status",
		App.New().Code("CODE").Status(401),
		Error{
			kind:    App,
			code:    "CODE",
			status:  401,
			context: make(Context),
		},
	}, {
		"set message",
		App.New().Code("CODE").Message("%d;%.2f;%s", 1, 1.2, "Hi"),
		Error{
			kind:    App,
			code:    "CODE",
			message: "1;1.20;Hi",
			context: make(Context),
		},
	}, {
		"context",
		App.New().Code("CODE").
			Context(Context{
				"Prop1": "Value1",
				"Prop2": 2,
			}).
			AddContext("Prop3", []int{5, 6}),
		Error{
			kind: App,
			code: "CODE",
			context: Context{
				"Prop1": "Value1",
				"Prop2": 2,
				"Prop3": []int{5, 6},
			},
		},
	}, {
		"custom",
		App.New().Code("custom.code").
			AddContext("Prop1", "Value1").
			AddContext("Prop2", 2.5).
			Context(Context{
				"Prop3": []int{5, 6},
			}).
			AddContext("Prop1", "Changed1"),
		Error{
			kind: App,
			code: "custom.code",
			context: Context{
				"Prop1": "Changed1",
				"Prop2": 2.5,
				"Prop3": []int{5, 6},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !errors.Is(test.err1, test.err2) {
				t.Errorf("\nExp: %s\nAct: %s", test.err2, test.err1)
			}
		})
	}
}

func TestNestedErrores(t *testing.T) {
	tests := []struct {
		name  string
		err1  error
		err2  error
		equal bool
	}{{
		"same error",
		App.New().Code("code1"),
		App.New().Code("code1"),
		true,
	}, {
		"wrapped error",
		App.New().Code("code2").Wrap(App.New().Code("code1")),
		App.New().Code("code1"),
		true,
	}, {
		"different errors",
		App.New().Code("code2"),
		App.New().Code("code1"),
		false,
	}, {
		"different nested errors",
		App.New().Code("code2").Wrap(App.New().Code("code1")),
		App.New().Code("code3").Wrap(App.New().Code("code4")),
		false,
	}, {
		"wrapped error",
		App.New().Code("code1"),
		App.New().Code("code2").Wrap(App.New().Code("code1")),
		false,
	}, {
		"wrapped internal error",
		App.New().Code("code1").Wrap(
			App.New().Code("code2").Wrap(
				Internal.New().Code("code3"),
			),
		),
		Internal.New().Code("code3"),
		true,
	}, {
		"middle error",
		App.New().Code("code1").Wrap(
			App.New().Code("code2").Wrap(
				Internal.New().Code("code3"),
			),
		),
		Internal.New().Code("code2"),
		true,
	}}

	for _, test := range tests {
		if errors.Is(test.err1, test.err2) != test.equal {
			t.Errorf("\nExp: %s\nAct: %s", test.err2, test.err1)
		}
	}

}

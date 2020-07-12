package errors

import (
	"errors"
	"reflect"
	"testing"
)

func TestDisplay(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		includeInternal bool
		expected        error
	}{{
		"raw error",
		errors.New("code1"),
		true,
		&DisplayError{
			Kind:    Raw,
			Message: "code1",
		},
	}, {
		"simple internal error",
		Internal.New().Code("code1").Message("msg %d", 2),
		true,
		&DisplayError{
			Kind:    Internal,
			Code:    "code1",
			Message: "msg 2",
		},
	}, {
		"simple application error",
		App.New().Code("code1").Context(Context{"key1": "value1"}),
		true,
		&DisplayError{
			Kind: App,
			Code: "code1",
			Context: Context{
				"key1": "value1",
			},
		},
	}, {
		"neste error (includeInternal)",
		App.New().Code("code1").Path("errors.code1").Wrap(
			App.New().Code("code2").Wrap(
				Internal.New().Code("code3").AddContext("key1", 123).Wrap(
					Internal.New().Code("code4").Wrap(
						errors.New("code5"),
					),
				),
			),
		),
		true,
		&DisplayError{
			Kind: App,
			Code: "code1",
			Path: "errors.code1",
			Cause: &DisplayError{
				Kind: App,
				Code: "code2",
				Cause: &DisplayError{
					Kind:    Internal,
					Code:    "code3",
					Context: Context{"key1": 123},
					Cause: &DisplayError{
						Kind: Internal,
						Code: "code4",
						Cause: &DisplayError{
							Kind:    Raw,
							Message: "code5",
						},
					},
				},
			},
		},
	}, {
		"app error inside internal",
		App.New().Code("code1").Wrap(
			Internal.New().Code("code2").Wrap(
				App.New().Code("code3"),
			),
		),
		false,
		&DisplayError{
			Kind: App,
			Code: "code1",
		},
	}, {
		"raw error inside app with includeInternal",
		App.New().Code("code1").Status(404).Wrap(
			errors.New("code2"),
		),
		true,
		&DisplayError{
			Kind:   App,
			Code:   "code1",
			Status: 404,
			Cause: &DisplayError{
				Kind:    Raw,
				Message: "code2",
			},
		},
	}, {
		"raw error inside app without includeInternal",
		App.New().Code("code1").Status(404).Wrap(
			errors.New("code2"),
		),
		false,
		&DisplayError{
			Kind:   App,
			Code:   "code1",
			Status: 404,
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dispErr := Display(test.err, test.includeInternal)

			if !reflect.DeepEqual(test.expected, dispErr) {
				t.Errorf("\nExp: %s\nAct: %s", test.expected, dispErr)
			}
		})
	}
}

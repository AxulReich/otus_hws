package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:10"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:3"`
		Marks  []int    `validate:"in:3,4,5"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	InvalidTagStruct struct {
		Val1 int     `validate:"mama:9"`
		Val2 int     `validate:"min:a"`
		Val3 string  `validate:"int:one,two"`
		Val4 float64 `validate:"max:10"`
		Val5 [][]int `validate:"max:10"`
	}
)

func (a App) String() string {
	return a.Version
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "",
			expectedErr: ErrNotStruct,
		},
		{
			in:          0,
			expectedErr: ErrNotStruct,
		},
		{
			in:          []struct{}{},
			expectedErr: ErrNotStruct,
		},
		{
			in:          nil,
			expectedErr: ErrNotStruct,
		},
		{
			in:          fmt.Stringer(nil),
			expectedErr: ErrNotStruct,
		},
		{
			in:          fmt.Stringer(App{Version: "12345"}),
			expectedErr: nil,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1",
			},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Version",
				Err:   ErrInvalidStringLength,
			}},
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 502,
				Body: "Bad Gateway",
			},
			expectedErr: ValidationErrors{ValidationError{
				Field: "Code",
				Err:   ErrInvalidIntValue,
			}},
		},
		{
			in: User{
				ID:     "Привет",
				Name:   "",
				Age:    12,
				Email:  "sobaka.ru",
				Role:   "bigboy",
				Phones: []string{"0", "1"},
				Marks:  []int{2, 2},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrInvalidStringLength,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrInvalidIntMin,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrInvalidStringSignature,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrInvalidStringValue,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrInvalidStringLength,
				},
				ValidationError{
					Field: "Marks",
					Err:   ErrInvalidIntValue,
				},
			},
		},
		{
			in: User{
				ID:     "Привет Чел",
				Name:   "Чел",
				Age:    19,
				Email:  "yoyo@mail.ru",
				Role:   "stuff",
				Phones: []string{"123", "321"},
				Marks:  []int{4, 5},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in: InvalidTagStruct{},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Val1",
					Err:   ErrUnexpectedRule,
				},
				ValidationError{
					Field: "Val2",
					Err:   ErrCastRuleValue,
				},
				ValidationError{
					Field: "Val3",
					Err:   ErrUnexpectedRule,
				},
				ValidationError{
					Field: "Val4",
					Err:   ErrUnsupportedType,
				},
				ValidationError{
					Field: "Val5",
					Err:   ErrUnsupportedType,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Equal(t, err, tt.expectedErr)
		})
	}
}

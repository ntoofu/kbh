package lib

import "testing"

func TestIsIncludedByStrSlice(t *testing.T) {
	type argsAndExpect struct {
		arg1 []string
		arg2 []string
		expected bool
	}

	testcases := []argsAndExpect{
		argsAndExpect{
			[]string{},
			[]string{},
			true,
		},
		argsAndExpect{
			[]string{"bar"},
			[]string{"foo", "bar", "baz"},
			true,
		},
		argsAndExpect{
			[]string{"bar", "BAZ"},
			[]string{"foo", "bar", "baz"},
			false,
		},
		argsAndExpect{
			[]string{"f"},
			[]string{"foo", "bar", "baz"},
			false,
		},
	}

	for _, tc := range testcases {
		if isIncludedByStrSlice(tc.arg1, tc.arg2) != tc.expected {
			t.Errorf("isIncludedByStrSlice(%v, %v) does not return %v", tc.arg1, tc.arg2, tc.expected)
		}
	}
}

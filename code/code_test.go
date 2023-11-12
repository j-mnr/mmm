package code_test

import (
	"mmm/code"
	"mmm/is"
	"testing"
)

func TestMake(t *testing.T) {
	t.Parallel()
	for name, tc := range map[string]struct{
		op code.Opcode
		operands []int
		want code.Instructions
	} {
		"Constant Ints": {
			op: code.OpConstant,
			operands: []int{0xFFFE},
			want: code.Instructions{byte(code.OpConstant), 0xFF, 0xFE},
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := code.Make(tc.op, tc.operands)
			is.Equal(t, len(tc.want), len(got))
		})
	}
}

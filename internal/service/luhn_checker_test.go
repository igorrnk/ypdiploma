package service

import "testing"

func TestLuhnChecker_Check(t *testing.T) {
	type args struct {
		order string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "good",
			args: args{
				"12345678903",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := LuhnChecker{}
			if got := lc.Check(tt.args.order); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

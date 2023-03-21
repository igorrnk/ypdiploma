package service

import "testing"

func TestGenHash(t *testing.T) {
	t.Logf("user1: %v", GenHash("user1", "97675768"))

}

func TestGenSalt(t *testing.T) {
	t.Logf("salt: %v", GenSalt(8))
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenSalt(tt.args.n)
			if got != tt.want {
				t.Errorf("GenSalt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package domain

import "testing"

func Test_GetSHA256(t *testing.T) {
	type args struct {
		phrase string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "positive 1",
			args: args{phrase: "001656709200Initial Block in the Chain"},
			want: "aed11b71d6c952f7d45756eae9c951f50d2f8902f736662421eb139855f87edd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSHA256([]byte(tt.args.phrase)); got != tt.want {
				t.Errorf("GetSHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

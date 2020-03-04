package main

import "testing"

func TestNormal(t *testing.T) {
	type args struct {
		mean float64
		sd   float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"", args{5.0, 1.0}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normal(tt.args.mean, tt.args.sd); got != tt.want {
				t.Errorf("Normal() = %v, want %v", got, tt.want)
			}
		})
	}
}

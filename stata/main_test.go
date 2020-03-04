package main

import (
	"testing"

	"github.com/drgo/sim/stata"
)

func Test_genWeibull(t *testing.T) {
	type args struct {
		k      float64
		lambda float64
		size   int64
	}
	tests := []struct {
		name string
		args args
		want args
	}{
		{"test1", args{k: 3, lambda: 5, size: 10e4}, args{k: 2, lambda: 5, size: 10e4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := genWeibull(tt.args.k, tt.args.lambda, tt.args.size)
			sf := stata.NewFile()
			sf.AddField("w", "weibull", data)
			if err := sf.WriteFile("weibull.dta"); err != nil {
				t.Fatal("error writing weibull file to disk")
				//t.Errorf("genWeibull() = %v, want %v", got, tt.want)
			}
		})
	}
}

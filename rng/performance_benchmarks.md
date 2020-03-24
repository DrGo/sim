# performance of Alias method

## original implementation
pkg: github.com/drgo/abm/rng
BenchmarkScalingAlias/n-levels=2-12         	50000000	        28.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingAlias/n-levels=4-12         	50000000	        33.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingAlias/n-levels=8-12         	50000000	        27.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingAlias/n-levels=16-12        	50000000	        27.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingAlias/n-levels=32-12        	50000000	        29.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingAlias/n-levels=64-12        	50000000	        35.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingAlias/n-levels=128-12       	50000000	        28.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingAlias/n-levels=256-12       	50000000	        27.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingAlias/n-levels=512-12       	50000000	        34.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingAlias/n-levels=1024-12      	50000000	        28.8 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/drgo/abm/rng	15.367s

## after removing size to constructor
BenchmarkScalingFreqDistributionSampler/n-levels=2-12         	50000000	        27.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=4-12         	50000000	        32.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=8-12         	50000000	        27.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=16-12        	50000000	        27.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=32-12        	50000000	        29.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=64-12        	50000000	        34.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=128-12       	50000000	        28.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=256-12       	50000000	        27.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=512-12       	50000000	        34.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=1024-12      	50000000	        28.6 ns/op	       0 B/op	       0 allocs/op
PASS

## after using to float32
pkg: github.com/drgo/abm/rng
BenchmarkScalingFreqDistributionSampler/n-levels=2-12         	50000000	        32.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=4-12         	50000000	        35.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=8-12         	50000000	        32.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=16-12        	50000000	        32.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=32-12        	50000000	        32.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=64-12        	50000000	        36.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=128-12       	50000000	        31.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=256-12       	50000000	        31.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=512-12       	50000000	        36.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=1024-12      	50000000	        31.8 ns/op	       0 B/op	       0 allocs/op

## using pcg
BenchmarkScalingFreqDistributionSampler/n-levels=2-12         	100000000	        13.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=4-12         	100000000	        17.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=8-12         	100000000	        13.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=16-12        	100000000	        13.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=32-12        	100000000	        14.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=64-12        	100000000	        19.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=128-12       	100000000	        14.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=256-12       	100000000	        13.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=512-12       	100000000	        18.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n-levels=1024-12      	100000000	        14.0 ns/op	       0 B/op	       0 allocs/op

## using pcg + FastBound optimization
BenchmarkScalingFreqDistributionSampler/n_levels=2-12         	100000000	        11.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=4-12         	100000000	        14.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=8-12         	100000000	        11.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=16-12        	100000000	        11.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=32-12        	100000000	        12.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=64-12        	100000000	        16.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=128-12       	100000000	        12.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=256-12       	100000000	        11.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=512-12       	100000000	        15.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=1024-12      	100000000	        12.3 ns/op	       0 B/op	       0 allocs/op

## using pcg + FastBound optimization + Go 1.12.0
enchmarkScalingFreqDistributionSampler/n_levels=2-12         	100000000	        11.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=4-12         	100000000	        14.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=8-12         	100000000	        11.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=16-12        	100000000	        11.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=32-12        	100000000	        11.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=64-12        	100000000	        15.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=128-12       	100000000	        11.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=256-12       	100000000	        11.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=512-12       	100000000	        15.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=1024-12      	100000000	        12.0 ns/op	       0 B/op	       0 allocs/op

## using pcg + FastBound optimization + Go 1.12.0 + minor tweaking
BenchmarkScalingFreqDistributionSampler/n_levels=2-12         	100000000	        11.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=4-12         	100000000	        14.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=8-12         	100000000	        11.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=16-12        	100000000	        11.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=32-12        	100000000	        11.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=64-12        	100000000	        15.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=128-12       	100000000	        11.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=256-12       	100000000	        11.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=512-12       	100000000	        15.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=1024-12      	100000000	        11.9 ns/op	       0 B/op	       0 

# using Xoroshiro as rand.Source
BenchmarkScalingFreqDistributionSamplerXoroshiro/n_levels=2-12         	100000000	        19.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSamplerXoroshiro/n_levels=4-12         	100000000	        22.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSamplerXoroshiro/n_levels=8-12         	100000000	        18.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSamplerXoroshiro/n_levels=16-12        	100000000	        19.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSamplerXoroshiro/n_levels=32-12        	100000000	        18.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSamplerXoroshiro/n_levels=64-12        	100000000	        18.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSamplerXoroshiro/n_levels=128-12       	100000000	        18.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSamplerXoroshiro/n_levels=256-12       	100000000	        18.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSamplerXoroshiro/n_levels=512-12       	100000000	        25.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSamplerXoroshiro/n_levels=1024-12      	100000000	        18.9 ns/op	       0 B/op	       0 allocs/op


## using pcg FastBound optimization + Go 1.12.0 + sfrand
BenchmarkScalingFreqDistributionSampler/n_levels=2-12         	100000000	        14.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=4-12         	100000000	        14.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=8-12         	100000000	        14.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=16-12        	100000000	        13.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=32-12        	100000000	        14.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=64-12        	100000000	        14.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=128-12       	100000000	        14.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=256-12       	100000000	        14.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=512-12       	100000000	        14.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkScalingFreqDistributionSampler/n_levels=1024-12      	100000000	        13.9 ns/op	       0 B/op	       0 allocs/op
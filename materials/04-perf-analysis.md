# Analyzing Performance

Now that we've got virtual nodes in place, let's add some benchmarks and analysis tools we can use to verify it's
working well. We're going to leverage some statistical techniques to analyze the dispersion of keys across the servers
(and virtual nodes).

If inclined, you can read about CV [here](https://en.wikipedia.org/wiki/Coefficient_of_variation). This is not necessary
for this lab, but you may be interested in it.

## Steps

1. Run `task test:bench` to see current benchmarks.
1. Take a look at hashing_bench_test.go.
1. Add the remaining bench mark tests (see TODO comments).
1. Run `task test:bench` to see the results
1. Take a look at examples/performance/main.go to see how we can use the new `AnalyzePerformance` method.
1. Run `task demo:performance` (can provide flags `task demo:performance -- -keys 100000`)

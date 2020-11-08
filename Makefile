test:
	go test

bench:
	go test -bench . -count 5

.PHONY: test bench

test:
	go test

bench:
	go test -bench .

clean:
	rm -rf *.log

.PHONY: test bench clean

default: test

fmt: 
	go fmt ./...

coverage: fmt
	go test ./ -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

test: fmt 
	go vet ./...
	go test ./...

pprof:
	go test -c
	./gorbac.test -test.cpuprofile cpu.out -test.bench .
	go tool pprof gorbac.test cpu.out
	rm cpu.out gorbac.test

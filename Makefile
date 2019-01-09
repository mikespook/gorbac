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
	./gorbac.test -test.cpuprofile cpu.prof -test.bench .
	go tool pprof gorbac.test cpu.prof
	rm cpu.prof gorbac.test

flamegraph:
	go test -c
	./gorbac.test -test.cpuprofile cpu.prof -test.bench .
	go-torch ./gorbac.test cpu.prof
	xdg-open torch.svg
	sleep 5
	rm cpu.prof gorbac.test torch.svg

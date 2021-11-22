cli:
	@make cli-lookup

cli-lookup:
	go build -mod vendor -o bin/lookup cmd/lookup/main.go

cli-complex:
	go build -mod vendor --tags json1 -o bin/current-complex cmd/current-complex/main.go

compile:
	@make compile-gates
	@make compile-galleries
	@make cli-lookup

compile-gates:
	go run -mod vendor cmd/compile-gates-data/main.go

compile-galleries:
	go run -mod vendor cmd/compile-galleries-data/main.go

GOMOD=vendor

cli:
	@make cli-lookup
	go build -mod $(GOMOD) -o bin/supersede-gallery cmd/supersede-gallery/main.go

cli-lookup:
	go build -mod $(GOMOD) -o bin/lookup cmd/lookup/main.go

cli-complex:
	go build -mod $(GOMOD) --tags json1 -o bin/current-complex cmd/current-complex/main.go

compile:
	@make compile-gates
	@make compile-galleries
	@make compile-terminals
	@make cli-lookup

compile-gates:
	go run -mod $(GOMOD) cmd/compile-gates-data/main.go

compile-terminals:
	go run -mod $(GOMOD) cmd/compile-terminals-data/main.go

compile-galleries:
	go run -mod $(GOMOD) cmd/compile-galleries-data/main.go

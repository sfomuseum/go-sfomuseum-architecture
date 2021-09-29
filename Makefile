compile-all:
	@make compile-gates
	@make compile-galleries

compile-gates:
	go run -mod vendor cmd/compile-gates-data/main.go

compile-galleries:
	go run -mod vendor cmd/compile-galleries-data/main.go

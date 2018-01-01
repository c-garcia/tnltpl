.PHONY: see clean test

tnltpl: tnltpl.go
	@go build tnltpl.go

test: tnltpl
	@go test

clean:
	@rm -f tnltpl
	@rm -f *~
see:
	@curl -s -L -H "Content-type: application/json" http://localhost:4040/api/tunnels

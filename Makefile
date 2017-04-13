install:
	@go install

test:
	@go test -v

clean:
	@go clean
	@rm -v -f liverebuild

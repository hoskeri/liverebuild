install:
	@go install

test:
	@go build
	@cd ./testdata/ && ../liverebuild

clean:
	@go clean
	@rm -v -f liverebuild

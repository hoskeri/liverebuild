install:
	@go build -o ./liverebuild

clean:
	@go clean
	@rm -v -f liverebuild

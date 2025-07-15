.PHONY: build
build: 
	@cd backend && go build -o ../dist/ -ldflags "-s -w" -trimpath ./...
	@cp -r frontend dist/templates

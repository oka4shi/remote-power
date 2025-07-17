.PHONY: build
build: 
	@cd backend && go build -o ../dist/ -ldflags "-s -w" -trimpath ./...
	@cd frontend && npm run build && cp -r dist ../dist/template

.PHONY: clean
clean:
	@rm -rf dist

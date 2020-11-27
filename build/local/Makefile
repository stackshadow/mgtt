# Build with make -f ./build/local/Makefile

include .env
export

.PHONY: help 
help: currentversion
	@echo "make docker: Build docker-image"
	@echo "make images: Will build all images from mermaid-files"
	@echo "make tests: Run all tests"
	@echo "make clean: clean all files which are created during build"

docker:
	docker-compose build
docker-run:
	docker-compose run mgtt
docker-down:
	docker-compose down --volumes --remove-orphans

currentversion:
	@echo "package cli" > internal/mgtt/cli/currentversion.go
	@echo "" >> internal/mgtt/cli/currentversion.go
	@echo "var Version string = \"$${VERSION}\"" >> internal/mgtt/cli/currentversion.go

mgtt: mgtt/mgtt
mgtt/mgtt:
	cd mgtt && \
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o mgtt .

docs/assets/node_modules/.bin/mmdc:
	cd docs && \
	yarn add -s @mermaid-js/mermaid-cli

docs/assets/QoS2.png: mmdc
	cd docs/assets && \
	./node_modules/.bin/mmdc -p puppeteer-config.json -i QoS2.mmd -o QoS2.png

mmdc: docs/assets/node_modules/.bin/mmdc
images: docs/assets/QoS2.png
docs: images



tests:
	cd mgtt ;\
	go test -timeout 30s -parallel 1 -coverprofile=coverage.out  ./... ;\
	go tool cover -html=coverage.out -o coverage.html

gocyclo: ~/go/bin/gocyclo
~/go/bin/gocyclo:
	go get github.com/fzipp/gocyclo/cmd/gocyclo
cyclic: gocyclo
		

clean:
	@rm -fv coverage.*
	@rm -fv *.pem
	@rm -fv mgtt/mgtt
	@find . -name "*.db" -print -delete
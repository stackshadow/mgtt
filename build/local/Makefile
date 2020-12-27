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

# command line mmdc
mmdc: docs/assets/node_modules/.bin/mmdc
docs/assets/node_modules/.bin/mmdc:
	cd docs && \
	yarn add -s @mermaid-js/mermaid-cli

# images
docs/assets/QoS2.png: mmdc
	cd docs/assets && \
	./node_modules/.bin/mmdc -p puppeteer-config.json -i QoS2.mmd -o QoS2.png


assets: docs/assets/QoS2.png
docs: assets



tests:
	go test -timeout 30s -parallel 1 -coverprofile=coverage.out  ./... ;\
	go tool cover -html=coverage.out -o coverage.html

anybadge: ${HOME}/.local/bin/anybadge
${HOME}/.local/bin/anybadge:
	pip install --user anybadge

gocyclo: ~/go/bin/gocyclo
~/go/bin/gocyclo:
	go get github.com/fzipp/gocyclo/cmd/gocyclo
cyclic: gocyclo anybadge
	VALUE=$$(~/go/bin/gocyclo -avg -ignore "_test|Godeps|vendor/" . | grep Average |  cut -d ' ' -f 2); \
	${HOME}/.local/bin/anybadge -l "Cyclomatic complexity" -v $$VALUE -f gocyclo.svg  10=green 15=yellow 18=orange 20=red


clean:
	@rm -fv coverage.*
	@rm -fv *.pem
	@rm -fv mgtt/mgtt
	@find . -name "*.db" -print -delete
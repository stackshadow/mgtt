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

${HOME}/.local/bin/anybadge:
	pip install --user anybadge
${HOME}/go/bin/gocyclo: ${HOME}/.local/bin/anybadge
	go get github.com/fzipp/gocyclo/cmd/gocyclo

###################### Badge: Coverage ######################
./coverage.txt:
	go tool cover -func coverage.out | grep total: | awk '{print $$3}' | sed 's/%//g' > $@
./coverage.svg: ./coverage.txt
	${HOME}/.local/bin/anybadge -l "Coverage" -v $$(cat ./coverage.txt) -f $@ 70=green 50=yellow 30=orange 0=red

###################### Badge: Cyclomatic complexity ######################
./gocyclo.txt: ${HOME}/go/bin/gocyclo
	${HOME}/go/bin/gocyclo -avg -ignore "_test|Godeps|vendor/" . | grep Average |  cut -d ' ' -f 2 > $@
./gocyclo.svg: ./gocyclo.txt ${HOME}/.local/bin/anybadge
	${HOME}/.local/bin/anybadge -l "Cyclomatic complexity" -v $$(cat ./gocyclo.txt) -f $@  10=green 15=yellow 18=orange 20=red

###################### Badge: Last build ######################
./lastbuild.svg: ${HOME}/.local/bin/anybadge
	${HOME}/.local/bin/anybadge -l "Last build" -v $$(date +%d.%m.%Y) -f $@ -c gray

badges: ./gocyclo.svg ./lastbuild.svg

gocyto: ~/go/bin/gocyto
~/go/bin/gocyto:
	go get github.com/protolambda/gocyto
callgraph: gocyto
	MACHINE=$$(uname -m); \
	if [ $$MACHINE == "x86_64" ]; then \
		export GOARCH="amd64"; \
	fi; \
	echo $$GOARCH ;\
	~/go/bin/gocyto ...

clean:
	@rm -fv coverage.*
	@rm -fv *.pem
	@rm -fv mgtt/mgtt
	@rm -fv gocyclo.txt gocyclo.svg lastbuild.svg
	@find . -name "*.db" -print -delete
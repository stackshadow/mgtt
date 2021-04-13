# Build with make -f ./build/local/Makefile

include .env
export

HELPTEXT="Targets:\n"

.PHONY: help 
help: currentversion
	@echo -e $(HELPTEXT)

HELPTEXT+="\nmake currentversion\n Create the current currentversion.go\n"
currentversion:
	@echo "package cli" > internal/mgtt/cli/currentversion.go
	@echo "" >> internal/mgtt/cli/currentversion.go
	@echo "var Version string = \"$${VERSION}\"" >> internal/mgtt/cli/currentversion.go

HELPTEXT+="\nmake mgtt\n Create mgtt binary. Need a go-compiler\n"
mgtt: ./mgtt
./mgtt:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o mgtt ./cmd/mgtt

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

HELPTEXT+="\nmake tests\n Run all tests\n"
tests: ./coverage.html
	go tool cover -func coverage.out | grep total: | awk '{print $$3}' | sed 's/%//g' 

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

HELPTEXT+="\nmake clean\n Remove all build-files\n"
clean:
	@rm -fv coverage.*
	@rm -fv *.pem
	@rm -fv mgtt
	@rm -fv gocyclo.txt gocyclo.svg 
	@rm -fv gosec.txt gosec.svg
	@rm -fv lastbuild.svg
	@find . -name "*.db" -print -delete
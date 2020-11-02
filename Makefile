
.PHONY: help
help:
	@echo "make docker: Build docker-image"
	@echo "make images: Will build all images from mermaid-files"
	@echo "make tests: Run all tests"
	@echo "make clean: clean all files which are created during build"

docker:
	docker-compose build


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
	
clean:
	@rm -fv coverage.*
	@rm -fv *.pem
	@find . -name "*.db" -print -delete
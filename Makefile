
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
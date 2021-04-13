.PHONY: help
help:
	@echo -e $(HELPTEXT)

HELPTEXT="Targets:\n"
BADGETARGETS=
CLEANTARGETS=

${HOME}/.local/bin/anybadge:
	pip install --user anybadge
${HOME}/go/bin/gocyclo: ${HOME}/.local/bin/anybadge
	go get github.com/fzipp/gocyclo/cmd/gocyclo
${HOME}/go/bin/gosec:
	go get github.com/securego/gosec/cmd/gosec


###################### Badge: Coverage ######################
./coverage.out:
	go test -v -timeout 60s -parallel 1 -coverprofile=coverage.out  $$(go list ./... | grep -v /mocked)
./coverage.txt: ./coverage.out
	go tool cover -func coverage.out | grep total: | awk '{print $$3}' | sed 's/%//g' > $@
./coverage.svg: ./coverage.txt ${HOME}/.local/bin/anybadge
	${HOME}/.local/bin/anybadge -l "Coverage" -v $$(cat ./coverage.txt) -f $@ 70=green 50=yellow 30=orange 0=red
./coverage.clean:
	@rm -fv ./coverage.out
	@rm -fv ./coverage.txt
	@rm -fv ./coverage.svg

BADGETARGETS+=./coverage.svg
CLEANTARGETS+=./coverage.clean

###################### Badge: Cyclomatic complexity ######################
./gocyclo.txt: ${HOME}/go/bin/gocyclo
	${HOME}/go/bin/gocyclo -avg -ignore "_test|Godeps|vendor/" . | grep Average |  cut -d ' ' -f 2 > $@
./gocyclo.svg: ./gocyclo.txt ${HOME}/.local/bin/anybadge
	${HOME}/.local/bin/anybadge -l "Cyclomatic complexity" -v $$(cat ./gocyclo.txt) -f $@  10=green 15=yellow 18=orange 20=red
./gocyclo.clean:
	@rm -fv ./gocyclo.txt
	@rm -fv ./gocyclo.svg

BADGETARGETS+=./gocyclo.svg
CLEANTARGETS+=./gocyclo.clean

###################### Badge: Last build ######################
./lastbuild.svg: ${HOME}/.local/bin/anybadge
	${HOME}/.local/bin/anybadge -l "Last build" -v $$(date +%d.%m.%Y) -f $@ -c gray
./lastbuild.clean:
	@rm -fv ./lastbuild.svg

BADGETARGETS+=./lastbuild.svg
CLEANTARGETS+=./lastbuild.clean

###################### Badge: gosec ######################
./gosec.txt: ${HOME}/go/bin/gosec
	${HOME}/go/bin/gosec -no-fail -severity medium ./... > $@
./gosec.svg: ./gosec.txt ${HOME}/.local/bin/anybadge
	value=$$(cat gosec.txt | grep Issues | cut -d':' -f2 | sed -e 's/\s*//') &&\
	${HOME}/.local/bin/anybadge -l "gosec" -v $$value -f $@ 10=green 50=yellow 70=orange 100=red
./gosec.clean:
	@rm -fv ./gosec.txt
	@rm -fv ./gosec.svg

BADGETARGETS+=./gosec.svg
CLEANTARGETS+=./gosec.clean

HELPTEXT+="\nmake badges\n Create all badges\n"
badges: $(BADGETARGETS)

HELPTEXT+="\nmake clean\n Remove all badged and temporary data\n"
clean: $(CLEANTARGETS)

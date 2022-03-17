all: spacesvr spacectl

.PHONY: spacesvr
spacesvr:
	@./scripts/build.sh spacesvr

.PHONY: spacectl
spacectl:
	@./scripts/build.sh spacectl

test: all
	@./scripts/test.sh

package:
	@rm -rf packages/
	@scripts/package.sh Windows windows amd64
	@scripts/package.sh Mac     darwin  amd64
	@scripts/package.sh Linux   linux   amd64
	@scripts/package.sh FreeBSD freebsd amd64
	@scripts/package.sh ARM     linux   arm
	@scripts/package.sh ARM64   linux   arm64

clean:
	rm -rf spacesvr spacectl 

distclean: clean
	rm -rf packages/

install: all
	cp spacesvr /usr/local/bin
	cp spacectl /usr/local/bin

uninstall: 
	rm -f /usr/local/bin/spacesvr
	rm -f /usr/local/bin/spacectl

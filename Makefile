.PHONY: setup

setup: bin/atlas

bin/atlas:
	curl -sSf https://atlasgo.sh | sh -s -- -y --no-install -o $(CURDIR)/bin/atlas

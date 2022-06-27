GO_PACKAGES=. ./cli ./clierr ./clivrsn ./validator ./version
GO_MODULE_NAME=github.com/mcdonaldseanp/clibuild
GO_BIN_NAME=clibuild

# Make the build dir, and remove any go bins already there
setup:
	mkdir -p output/
	cd output && \
	rm -f $(GO_BIN_NAME)

# Actually build the thing
build: setup
	go mod tidy
	go build -o output/ $(GO_MODULE_NAME)

install:
	go mod tidy
	go install $(GO_MODULE_NAME)

# Build it before publishing to make sure this publication won't be broken
#
# This also ensures that the clibuild command is available for the version
# command
#
# If NEW_VERSION is set by the user, it will set the new clibuild version
# to that value. Otherwise clibuild will bump the Z version
publish: install format
	NEW_VERSION=$$(clibuild update version ./version/version.go "$(NEW_VERSION)") && \
	echo "Tagging and publishing new version $$NEW_VERSION" && \
	git add --all && \
	git commit -m "(release) Update to new version $$NEW_VERSION" && \
	git tag -a $$NEW_VERSION -m "Version $$NEW_VERSION"
	git push
	git push --tags

format:
	go fmt $(GO_PACKAGES)
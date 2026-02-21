out:
	go install github.com/mitchellh/gox@latest
	rm -rf ./_obj
	PATH="$(shell go env GOPATH)/bin:$$PATH" gox -osarch="linux/386 linux/amd64 linux/arm windows/386 windows/amd64 darwin/amd64 darwin/arm64" -output "_obj/{{.OS}}_{{.Arch}}/{{.Dir}}" ./cmd/dajarep
	-test -d ./_obj/darwin_amd64 && mv ./_obj/darwin_amd64 ./_obj/macos_amd64 || true
	-test -d ./_obj/darwin_arm64 && mv ./_obj/darwin_arm64 ./_obj/macos_arm64 || true
	find ./_obj -mindepth 1 -maxdepth 1 -type d -exec sh -c 'zip -r -j "$$1.zip" "$$1" && rm -rf "$$1"' _ {} \;

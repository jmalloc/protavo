REQ += $(shell find src -name "*.pb.go")

-include artifacts/make/go/Makefile

%.pb.go: %.proto
	protoc --go_out=. $(@D)/*.proto

artifacts/make/%/Makefile:
	curl -sf https://jmalloc.github.io/makefiles/fetch | bash /dev/stdin $*

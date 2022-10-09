#
.PHONY: test
test:
	go test -v -count 1 ./...

.PHONY: compatible
compatible:
	docker build -t quicksilver-cli:latest -f Dockerfile_for_compatiable_test .


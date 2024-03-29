all:
	gox -osarch="linux/amd64 darwin/amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" metrik

dev:
	gox -osarch="darwin/amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" metrik

dist:
	gox -osarch="linux/amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" metrik

clean:
	rm -f ./bin/metrik*

rundarwin: dev
	./bin/metrik_darwin_amd64

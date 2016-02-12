#!/bin/bash -e
# The script does automatic checking on a Go package and its sub-packages, including:
# 1. gofmt         (http://golang.org/cmd/gofmt/)
# 2. goimports     (https://github.com/bradfitz/goimports)
# 3. golint        (https://github.com/golang/lint)
# 4. go vet        (http://golang.org/cmd/vet)
# 5. race detector (http://blog.golang.org/race-detector)
# 6. test coverage (http://blog.golang.org/cover)

# Capture what test we should run
TEST_SUITE=$1

if [[ $TEST_SUITE == "unit" ]]; then
	go get github.com/axw/gocov/gocov
	go get -u github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/vet
	go get golang.org/x/tools/cmd/goimports
	go get github.com/smartystreets/goconvey/convey
	go get golang.org/x/tools/cmd/cover
	go get github.com/marpaia/graphite-golang
	
	TEST_DIRS="main.go graphite logHelper "
	VET_DIRS=". ./graphite/... ./logHelper/..."

	set -e

	# Automatic checks
	echo "gofmt"
	test -z "$(gofmt -l -d $TEST_DIRS | tee /dev/stderr)"

	echo "goimports"
	test -z "$(goimports -l -d $TEST_DIRS | tee /dev/stderr)"

	# Useful but should not fail on link per: https://github.com/golang/lint
	# "The suggestions made by golint are exactly that: suggestions. Golint is not perfect,
	# and has both false positives and false negatives. Do not treat its output as a gold standard.
	# We will not be adding pragmas or other knobs to suppress specific warnings, so do not expect
	# or require code to be completely "lint-free". In short, this tool is not, and will never be,
	# trustworthy enough for its suggestions to be enforced automatically, for example as part of
	# a build process"
	# echo "golint"
	# golint ./...

	echo "go vet"
	go vet $VET_DIRS
	# go test -race ./... - Lets disable for now
 
	# Run test coverage on each subdirectories and merge the coverage profile.
	echo "mode: count" > profile.cov
 
	# Standard go tooling behavior is to ignore dirs with leading underscors
	for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -not -path './examples/*' -not -path './scripts/*' -not -path './build/*' -not -path './Godeps/*' -type d);
	do
		if ls $dir/*.go &> /dev/null; then
	    		go test --tags=unit -covermode=count -coverprofile=$dir/profile.tmp $dir
	    		if [ -f $dir/profile.tmp ]
	    		then
	        		cat $dir/profile.tmp | tail -n +2 >> profile.cov
	        		rm $dir/profile.tmp
	    		fi
		fi
	done
 
	go tool cover -func profile.cov
	
fi
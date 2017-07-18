medium:
	_docker_ps_id="$(docker run -d -p 80:80 -p 2003:2003 -p 2004:2004 -p 4444:4444/udp -p 8126:8126 -p 8086:8086 hopsoft/graphite-statsd)"
	export TEST_TYPE="${TEST_TYPE:-"medium"}"
	export SNAP_GRAPHITE_HOST="0.0.0.0"
	go test -coverprofile=coverage.out -tags=medium
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	docker kill $_docker_ps_id >/dev/null && docker rm $_docker_ps_id >/dev/null
small:
	go test -coverprofile=coverage.out -tags=small
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out	

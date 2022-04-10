build: test compile

run:
	printf "[ run ]: starting solution... "
	docker-compose up --build --remove-orphans --force-recreate

compile:
	echo "[ compile ]: Building ..."
	GOOS=linux go build -o bin/processor ./cmd/message-db-processor/...
	GOOS=linux go build -o bin/manager ./cmd/user-manager-server/...

test:
	echo "[ test ]: running unit tests..."
	go test -race ./... -count=1

e2e-test:
	echo "[ test ]: running integration tests..."
	docker-compose up --build --remove-orphans --force-recreate -d
	go test -v -tags=e2e ./e2e-tests/...  ; (ret=$$?;  docker-compose down && exit $$ret)

.PHONY: all build validate compile test
all:

start:
	go mod tidy
	go run cmd/main.go

test:
	# should have emulator running
	go clean -testcache
	export FIRESTORE_EMULATOR_HOST=127.0.0.1:8080 && export PUBSUB_EMULATOR_HOST=localhost:8085 && go test -v ./...

emulater:
	# can use kill -9 $(pidof java) to kill them
	# execute 2 commands below parallel
	make -j 2 emulater-pubsub emulater-firestore
emulater-pubsub:
	gcloud beta emulators pubsub start --project=tw-rd-ca-leon-lin --host-port=localhost:8085
emulater-firestore:
	gcloud emulators firestore start --host-port=127.0.0.1:8080

install-emulater:
	# should have java installed
	gcloud components install pubsub-emulator
	gcloud components update

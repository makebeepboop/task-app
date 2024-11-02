run:
	go run ./cmd/app/main.go --config=./config/local.yaml

migration:
	go run ./cmd/migrator/main.go --storage-path=./storage/app.db --migrations-path=./migrations/

rm-storage:
	rm -rf ./storage/app.db

clean-run: rm-storage migration run

build:
	go build -o docker/bin/notifier ./cmd/notifier
	cd docker && docker build -t myriaddreamin/bilibili-notifier:latest .

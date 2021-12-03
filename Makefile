
build:
	go build -o deployment/bin/notifier ./cmd/notifier
	cd deployment && docker build -t myriaddreamin/bilibili-notifier:latest .

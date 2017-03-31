BINDIR = bin
TARGET = $(BINDIR)/dungeonfs

build:
	@echo "Building ..."
	@mkdir -p $(BINDIR)
	@go build -o $(TARGET) ./cmd/dungeonfs

deps:
	@go get -u bazil.org/fuse
	@go get -u golang.org/x/net/context
	@go get -u github.com/spf13/cobra
	@go get -u github.com/spf13/viper

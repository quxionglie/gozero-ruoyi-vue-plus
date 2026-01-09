.PHONY: build run clean help

# 构建项目
build:
	@echo "Building..."
	@go build -o bin/gozero-ruoyi-vue-plus admin.go
	@echo "Build complete!"

# 运行项目
run:
	@echo "Running..."
	@go run admin.go -f etc/admin-api.yaml

# 清理构建文件
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf logs/
	@echo "Clean complete!"

# 安装依赖
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download
	@echo "Dependencies installed!"

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete!"

# 代码检查
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete!"

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  build   - Build the project"
	@echo "  run     - Run the project"
	@echo "  clean   - Clean build files"
	@echo "  deps    - Install dependencies"
	@echo "  fmt     - Format code"
	@echo "  vet     - Run go vet"
	@echo "  help    - Show this help message"
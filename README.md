# PocketBase Starter

This is a PocketBase starter project template that provides out-of-the-box basic
configuration and functionality.

## Main Features

### 1. Docker Containerized Deployment

- Provides complete Docker configuration files for one-click deployment
- Uses Alpine Linux as base image, image size less than 20MB
- Container orchestration and management through docker-compose
- Data persistence storage in root directory
- Port mapping: 9900:9900

### 2. Admin Account Configuration

- Flexibly configure admin account through environment variables:
  - ADMIN_EMAIL: Administrator email
  - ADMIN_PASSWORD: Administrator password
- Default configuration:
  - Email: admin@pocketbase-starter.com
  - Password: 0123456789
- Automatically checks and creates admin account on startup

## Quick Start

1. Clone the project

```bash
git clone https://github.com/a876691666/pocketbase-starter.git
cd pocketbase-starter
```

2. Run locally

```bash
go run main.go serve
```

3. Build

```bash
./build.sh linux    # Build Linux version
./build.sh windows  # Build Windows version  
./build.sh darwin   # Build MacOS version
./build.sh all      # Build all versions
```

4. Build image

```bash
./docker-build.sh    # Default build Linux version
./docker-build.sh -k # Build and keep removed old image
```

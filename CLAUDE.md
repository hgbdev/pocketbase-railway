# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a PocketBase deployment template for Railway with custom Go hooks. PocketBase is an open-source backend-as-a-service that provides a complete backend solution, extended here with custom authentication logic.

## Architecture

- **Custom Go application**: Uses PocketBase as a framework with custom hooks
- **Multi-stage Docker build**: Compiles Go application and runs in Alpine Linux container
- **Railway deployment**: Configured for easy deployment on Railway platform
- **Port configuration**: Exposes port 8080 for HTTP traffic
- **OTP Authentication**: Includes custom hook for automatic user creation during OTP requests

## Key Files

- `main.go`: Custom PocketBase application with OTP authentication hooks
- `go.mod` & `go.sum`: Go module dependencies
- `Dockerfile`: Multi-stage build for Go application compilation and deployment
- `README.md`: Contains deployment instructions and Railway template link
- `LICENSE`: MIT license file

## Development Commands

1. **Local development**: `go run main.go serve`
2. **Build application**: `go build -o pocketbase-custom .`
3. **Docker build**: `docker build -t pocketbase-railway .`
4. **Docker run**: `docker run -p 8080:8080 pocketbase-railway`
5. **Test dependencies**: `go mod tidy && go mod verify`

## Custom Hooks

### OTP User Creation Hook
- **Location**: `main.go:OnRecordAuthRequestOTP("users")`
- **Purpose**: Automatically creates a user record if one doesn't exist during OTP request
- **Behavior**: 
  - Extracts email from request body
  - Creates new user with random password
  - Sets the created user as the OTP target

## Configuration

- PocketBase version is controlled via the `PB_VERSION` ARG in the Dockerfile (currently 0.29.0)
- To update PocketBase version, modify the `PB_VERSION` value in the Dockerfile

## Deployment

The project is designed for Railway deployment using the template button in the README. No additional configuration files or build steps are required.
<img width="988" alt="image" src="https://github.com/user-attachments/assets/2a568f3b-4bbd-4cb4-afff-6e93c182176b" />

# Rate Limiter Project

A full-stack rate limiting system with a React frontend and Go backend that allows dynamic configuration of API rate limits.

## Features

- Real-time rate limit configuration through a modern UI
- Dynamic rate limiting for different API endpoints
- Distributed rate limiting using Redis and MongoDB
- Default and custom rate limits per endpoint
- Real-time monitoring and updates
- Dark mode UI with Material-UI components

## Architecture

### Frontend
- Built with React and TypeScript
- Material-UI (MUI) for component styling
- Responsive design with a modern dark theme
- Real-time updates for rate limit configurations

### Backend
- Go server with HTTP endpoints
- Redis for rate limit caching
- MongoDB for persistent storage
- RabbitMQ integration for message queuing
- Middleware-based rate limiting implementation

## API Endpoints

- `GET /listroutes` - List all registered routes and their rate limits
- `POST /ratelimit/changelimit` - Update rate limit for a specific endpoint
- `POST /ratelimit/setlimit` - Set default rate limit for an endpoint

## Getting Started

1. Clone the repository
2. Set up environment variables
3. Start the backend:
```bash
cd backend
go run main.go
```
4. Start the frontend:
```bash
cd frontend/ratelimiter
npm install
npm run dev
```

## Environment Variables

Backend requires:
- `URI` - MongoDB connection string
- `MONGO_DB` - MongoDB database name
- Redis configuration
- RabbitMQ configuration

## License

MIT

# Conference Room Booking Display

A Go web application that displays Microsoft Bookings events for a single conference room. Shows current day events with welcome messages for groups and upcoming events.

## Overview

This application provides a clean, web-based display interface for conference room bookings, making it easy to see:
- Current day's scheduled events
- Welcome messages for meeting groups
- Upcoming events at a glance

Perfect for mounting on a tablet or screen outside your conference room.

## Features

- 📅 **Daily Event Display** - Shows all bookings for the current day
- 👋 **Welcome Messages** - Displays personalized greetings for meeting groups
- ⏰ **Upcoming Events** - Keeps track of ongoing and upcoming events
- 🎨 **Clean Interface** - Simple, easy-to-read display format
- 🚀 **Lightweight** - Fast Go-based web server with embedded assets

## Requirements

- Go 1.24 or higher
- Access to Microsoft Bookings API (configuration required)
- Web browser for display

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/jmidd2/spark-bookings.git
   cd spark-bookings
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

## Configuration

The application uses environment variables for configuration:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `APP_ENV` | Application environment | `development` | No |
| `PORT` | Server port | `8080` (dev), `80` (prod) | No |
| `DOMAIN` | Server domain | `localhost` (dev) | Yes (prod) |
| `HTTP_SECURE` | Enable HTTPS | `false` | No |

### Example .env file:
```
env
APP_ENV=development
PORT=8080
DOMAIN=localhost
HTTP_SECURE=false
```
## Usage

1. **Start the server:**
   ```bash
   go run main.go
   ```

2. **Access the application:**
    - Development: http://localhost:8080
    - Production: Your configured domain

3. **Display on device:**
    - Open the URL in a web browser
    - For kiosk mode, use browser full-screen mode
    - Consider using a dedicated kiosk browser for public displays

## Project Structure
```
spark-bookings/
├── main.go              # Main application server
├── templates/           # HTML templates
│   └── index.html      # Main display template
├── assets/             # Static assets (CSS, JS, images)
│   └── stylesheets/    # CSS
├── go.mod              # Go module file
├── .env.example        # Environment variables template
└── README.md           # This file
```
## Microsoft Bookings Integration

This application is designed to integrate with Microsoft Bookings API to fetch:
- Room booking information
- Meeting titles and organizers
- Start and end times
- Attendee information for welcome messages

*Note: Microsoft Bookings API integration setup instructions will be added as the integration is implemented.*

## Development

### Local Development
```
bash
# Set development environment
export APP_ENV=development

# Run the server
go run main.go

# Server will start on http://localhost:8080
```
### Building for Production
```
bash
# Build binary
go build -o booking-display main.go

# Run binary
./booking-display
```
## Deployment

The application is designed to be easily deployable:

1. **Docker** (recommended for production)
2. **Direct binary deployment**
3. **Cloud platforms** (supports standard Go deployment)

### Environment Setup for Production
- Set `APP_ENV=production`
- Configure `DOMAIN` to your server's domain
- Set `HTTP_SECURE=true` if using HTTPS
- Configure appropriate `PORT` if not using standard ports

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request
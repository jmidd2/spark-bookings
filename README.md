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

## Development

### Local Development
1. **Clone the repository:**
   ```bash
   git clone https://github.com/jmidd2/spark-bookings.git
   cd spark-bookings
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up [environment variables](#configuration):**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start dev server**
    ```bash
    just dev
    # Server will start on http://localhost:8090
    ```
   
### Configuration

The application uses environment variables for configuration:

| Variable                 | Description                                            | Default                             | Required   |
|--------------------------|--------------------------------------------------------|-------------------------------------|------------|
| `APP_ENV`                | Application environment                                | `development`                       | No         |
| `PORT`                   | Server port                                            | `8080` (dev), `80` (prod)           | No         |
| `DOMAIN`                 | Server domain                                          | `localhost` (dev)                   | Yes (prod) |
| `HTTP_SECURE`            | Enable HTTPS                                           | `false`                             | No         |
| `BOOKINGS_CLIENT_ID`     | Microsoft Bookings Client ID                           | -                                   | Yes        |
| `BOOKINGS_CLIENT_SECRET` | Microsoft Bookings Client Secret                       | -                                   | Yes        |
| `BOOKINGS_TENANT_ID`     | Microsoft Bookings Tenant ID                           | -                                   | Yes        |
| `BOOKINGS_BUSINESS_ID`   | Microsoft Bookings Business ID                         | -                                   | Yes        |
| `BOOKINGS_AUTH_URL`      | Microsoft Bookings Authority URL without the Tenant ID | `https://login.microsoftonline.com` | No         |

### Configuring Bookings API

The application must first be registered with Microsoft Graph API in [Microsoft Identity Platform](https://learn.microsoft.com/en-us/graph/auth-register-app-v2). After registration, you will be provided with a client ID and client secret.

The application must be granted the following permissions:
- `Bookings.Read.All`
- `Bookings.ReadWrite.All`
- `User.Read`

The application uses the [Microsoft Bookings API](https://learn.microsoft.com/en-us/graph/api/resources/booking-api-overview?view=graph-rest-1.0) to fetch room booking information. To use the API, you must first register an application with Microsoft.

### Example .env file:
```
env
APP_ENV=development
PORT=8080
DOMAIN=localhost
HTTP_SECURE=false
BOOKINGS_CLIENT_ID=clientId
BOOKINGS_CLIENT_SECRET=clientSecret
BOOKINGS_TENANT_ID=tenantId
BOOKINGS_BUSINESS_ID=businessId
BOOKINGS_AUTH_URL=authorityUrl
```

## Usage

1. **Start the server:**
   ```bash
   just build
   ./bin/booking-display
   ```

2. **Access the application:**
    - Development: http://localhost:8080
    - Production: Your configured domain

3. **Display on device:**
    - Open the URL in a web browser
    - For kiosk mode, use browser full-screen mode
    - Consider using a dedicated kiosk browser for public displays

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

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request
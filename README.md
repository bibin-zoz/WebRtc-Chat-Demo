Here's a README file for your chat application:

---

# Chat Application

This is a real-time chat application developed using Go, Gin, Gorilla WebSocket, and PostgreSQL. The application allows users to join chat rooms and exchange messages in real time.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Database Setup](#database-setup)
  - [Running the Application](#running-the-application)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- Real-time messaging using WebSockets
- User authentication (future implementation)
- Persistent message storage with PostgreSQL
- Scalable architecture

## Getting Started

### Prerequisites

- Go (version 1.15 or later)
- PostgreSQL
- Git

### Installation

Clone the repository:

```sh
git clone https://github.com/your-username/chat-app.git
cd chat-app
```

Install dependencies:

```sh
go mod tidy
```

### Configuration

Create a `.env` file in the project root directory and add the following environment variables:

```env
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=chat_app
```

### Database Setup

1. Install and set up PostgreSQL on your system.
2. Create a PostgreSQL database:

```sql
CREATE DATABASE chat_app;
```

3. Update the database connection details in the `.env` file.

4. Run database migrations to set up the necessary tables:

```sh
go run migration.go
```

### Running the Application

Build and run the application:

```sh
go build
./chat-app
```

The server will start running on `http://localhost:8080`.

## Usage

### WebSocket Endpoint

To connect to the WebSocket:

```sh
ws://localhost:8080/ws?group_id=<GROUP_ID>
```

Replace `<GROUP_ID>` with the ID of the group you want to join.

### Sending Messages

To send a message, the client needs to send a JSON object through the WebSocket connection:

```json
{
  "sender_id": 1,
  "receiver_id": 2,
  "message": "Hello, World!",
  "time": "2024-07-11T10:00:00Z"
}
```




### WebSocket

- **Endpoint**: `/ws`
- **Description**: Handles WebSocket connections for real-time messaging.
- **Parameters**: 
  - `group_id`: ID of the group to join (query parameter)


## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License



---


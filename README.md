# WebRtc-Chat-Demo
This is a real-time chat application developed using Go, Gin, Gorilla WebSocket, and PostgreSQL. The application allows users to join chat rooms and exchange messages in real time.

Table of Contents
Features
Getting Started
Prerequisites
Installation
Configuration
Database Setup
Running the Application
Usage
Project Structure
API Endpoints
Contributing
License
Features
Real-time messaging using WebSockets
User authentication (future implementation)
Persistent message storage with PostgreSQL
Scalable architecture
Getting Started
Prerequisites
Go (version 1.15 or later)
PostgreSQL
Git
Installation
Clone the repository:

sh
Copy code
git clone https://github.com/your-username/chat-app.git
cd chat-app
Install dependencies:

sh
Copy code
go mod tidy
Configuration
Create a .env file in the project root directory and add the following environment variables:

env
Copy code
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=chat_app
Database Setup
Install and set up PostgreSQL on your system.
Create a PostgreSQL database:
sql
Copy code
CREATE DATABASE chat_app;
Update the database connection details in the .env file.

Run database migrations to set up the necessary tables:

sh
Copy code
go run migration.go
Running the Application
Build and run the application:

sh
Copy code
go build
./chat-app
The server will start running on http://localhost:8080.

Usage
WebSocket Endpoint
To connect to the WebSocket:

sh
Copy code
ws://localhost:8080/ws?group_id=<GROUP_ID>
Replace <GROUP_ID> with the ID of the group you want to join.

Sending Messages
To send a message, the client needs to send a JSON object through the WebSocket connection:

json
Copy code
{
  "sender_id": 1,
  "receiver_id": 2,
  "message": "Hello, World!",
  "time": "2024-07-11T10:00:00Z"
}

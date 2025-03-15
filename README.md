# Chat App API (Go)

A simple chat application built with **Go**, implementing the **Clean Architecture** design pattern. This project includes WebSocket communication, RabbitMQ for message brokering, Redis caching, and JWT-based authentication.

---

## 🚀 Features

- **User Registration & Login** with JWT Authentication
- **WebSocket Support** for real-time chat functionality
- **RabbitMQ Integration** for message queuing
- **Redis Cache** to store recent messages
- **Rate Limiting Middleware** for security enhancement
- **Clean Architecture Design Pattern** for scalability and maintainability

---

## 📂 Project Structure

```
/chat-app
 ├── /cmd
 │   └── main.go
 ├── /config
 │   └── db.go
 ├── /internal
 │   ├── /domain
 │   ├── /usecase
 │   ├── /repository
 │   ├── /delivery
 │   └── /middleware
 ├── /pkg
 │   ├── /broker
 │   ├── /cache
 │   └── /load_balancer
 ├── .env
 ├── .gitignore
 ├── go.mod
 ├── go.sum
 └── README.md
```

---

## 🛠️ Installation

1. **Clone the repository**

```bash
git clone https://github.com/your-username/chat-app.git
cd chat-app
```

2. **Install dependencies**

```bash
go mod tidy
```

3. **Setup Environment Variables** Create a `.env` file and configure it with:

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=chat_db
JWT_SECRET=your_secret_key
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
REDIS_ADDR=localhost:6379
```

4. **Run Database Migration**

```bash
mysql -u root -p < db/schema.sql
```

5. **Run the Application**

```bash
go run cmd/main.go
```

---

## 📡 API Endpoints

| Method | Endpoint    | Description        |
| ------ | ----------- | ------------------ |
| `POST` | `/register` | Register new user  |
| `POST` | `/login`    | Login and get JWT  |
| `GET`  | `/ws`       | WebSocket endpoint |

---

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -m 'Add some feature'`)
4. Push to your branch (`git push origin feature/your-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 💬 Contact

For any inquiries or issues, feel free to contact me at: [**your.email@example.com**](mailto\:your.email@example.com)

---

## 🧱 .gitignore

```
# Go-specific files
*.log
*.test
*.out

# Environment variables
.env

# IDE / Editor settings
.idea/
.vscode/
*.swp

# Dependency directories
vendor/

# OS generated files
.DS_Store
Thumbs.db
```


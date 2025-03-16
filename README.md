# Golang Chat API dengan WebSocket dan RabbitMQ

Proyek ini adalah aplikasi chat sederhana yang dibangun menggunakan **Golang**, dilengkapi dengan fitur:
- **WebSocket** untuk komunikasi real-time
- **RabbitMQ** sebagai message broker
- **MySQL** untuk penyimpanan data
- **Redis** untuk caching data
- **JWT (JSON Web Token)** untuk autentikasi pengguna

---

## 🚀 Fitur Utama
✅ Kirim dan terima pesan secara real-time dengan WebSocket  
✅ Autentikasi berbasis JWT  
✅ Penyimpanan pesan ke database menggunakan MySQL  
✅ Pengiriman pesan ke pengguna yang sedang online  
✅ Penggunaan **RabbitMQ** untuk distribusi pesan agar scalable  
✅ Redis Cache untuk mempercepat akses data chat

---

## 📂 Struktur Proyek
```
├── cmd
│   └── main.go                # Entry point aplikasi
├── config
│   └── db.go                  # Konfigurasi koneksi database
├── internal
│   ├── delivery
│   │   ├── ws_handler.go      # WebSocket Handler untuk chat
│   │   └── consumer.go        # Consumer RabbitMQ untuk menerima pesan
│   ├── repository
│   │   ├── chat_repository.go # Repository untuk penyimpanan data chat
│   │   └── user_repository.go # Repository untuk data pengguna
│   ├── usecase
│   │   ├── chat_usecase.go    # Logika bisnis untuk chat
│   │   └── user_usecase.go    # Logika bisnis untuk pengguna
│   └── domain
│       ├── chat.go            # Struct Chat
│       └── user.go            # Struct User
├── pkg
│   ├── broker
│   │   └── rabbitmq.go        # Konfigurasi RabbitMQ
│   ├── cache
│   │   └── redis.go           # Konfigurasi Redis
│   └── middleware
│       └── auth_middleware.go # Middleware untuk autentikasi JWT
└── README.md                  # Dokumentasi proyek ini
```

---

## ⚙️ Instalasi dan Konfigurasi
### 1. Clone Repository
```bash
git clone https://github.com/DimasAriyanto/golang-chat-api.git
cd golang-chat-api
```

### 2. Buat File Konfigurasi
Buat file `.env` dan isi dengan konfigurasi berikut:
```
DB_USER=root
DB_PASS=your_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=chat_db

REDIS_ADDR=localhost:6379

RABBITMQ_URL=amqp://guest:guest@localhost:5672/

JWT_SECRET=secret_key
```

### 3. Instalasi Dependency
```bash
go mod tidy
```

### 4. Menjalankan RabbitMQ dan Redis
Pastikan RabbitMQ dan Redis sudah berjalan. Jika menggunakan Docker:
```bash
docker-compose up -d
```

### 5. Jalankan Aplikasi
```bash
go run cmd/main.go
```

---

## 🗄️ Struktur Database
Gunakan skrip berikut untuk membuat database dan tabel:

```sql
CREATE DATABASE chat_db;

USE chat_db;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chat (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT,         -- NULL jika pesan dikirim ke grup
    group_id INT,            -- NULL jika pesan dikirim ke individu
    message TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 🔎 Penggunaan API Endpoint
### 1. **Register User**
**Endpoint:** `POST /register`
```json
{
    "username": "john_doe",
    "password": "securepassword"
}
```

### 2. **Login User**
**Endpoint:** `POST /login`
```json
{
    "username": "john_doe",
    "password": "securepassword"
}
```
**Respon:**
```json
{
    "token": "<JWT_TOKEN>"
}
```

### 3. **WebSocket Chat**
**Endpoint WebSocket:** `ws://localhost:8080/ws?token=<JWT_TOKEN>`

**Kirim Pesan:**
```json
{
    "receiver_id": 2,
    "message": "Halo, Klien 2!"
}
```

**Respon:**
```json
{
    "sender_id": 1,
    "receiver_id": 2,
    "message": "Halo, Klien 2!",
    "timestamp": "2025-03-16T14:21:58+07:00"
}
```

---

## 🧪 Pengujian
1. **Login User** untuk mendapatkan token JWT.  
2. Buka **dua WebSocket** melalui Postman atau `wscat`.  
3. Kirim pesan dari User 1 ke User 2 melalui WebSocket.  
4. Verifikasi bahwa User 2 menerima pesan dan data berhasil disimpan di database.

---

## 📄 Dokumentasi API
📌 [Postman Documentation](https://documenter.getpostman.com/view/22351148/2sAYkBt22a)

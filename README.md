# Books Rendezvous 📚

A modern book management and social reading platform built with Nuxt.js (frontend) and Go Fiber (backend).

## 🚀 Features

- **Book Management**: Add, edit, and organize your book collection
- **Social Reading**: Connect with other readers and share recommendations
- **Achievement System**: Track your reading progress with achievements
- **Modern UI**: Beautiful and responsive interface built with Nuxt.js and Tailwind CSS
- **Fast API**: High-performance backend built with Go Fiber
- **Database**: PostgreSQL for reliable data storage
- **Caching**: Optional Redis integration for improved performance

## 🏗️ Architecture

- **Frontend**: Nuxt.js 3 with Vue.js, Tailwind CSS, and Pinia for state management
- **Backend**: Go with Fiber framework for REST API
- **Database**: PostgreSQL 15
- **Cache**: Redis (optional)
- **Containerization**: Docker and Docker Compose for easy deployment

## 📋 Prerequisites

- Docker and Docker Compose
- Git

## 🚀 Quick Start

1. **Clone the repository**

   ```bash
   git clone <your-repo-url>
   cd booksrendezvous
   ```

2. **Configure environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your preferred settings
   ```

3. **Start the application**

   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:6050
   - Database: localhost:5432

## ⚙️ Configuration

### Environment Variables

Copy `.env.example` to `.env` and customize the following variables:

#### Database Configuration

- `DB_NAME`: Database name (default: bibliotheque)
- `DB_USER`: Database user (default: bibliotheque_user)
- `DB_PASSWORD`: Database password (default: bibliotheque_password)
- `DB_PORT`: Database port (default: 5432)

#### Backend Configuration

- `BACKEND_PORT`: Backend server port (default: 6050)
- `LOG_LEVEL`: Logging level (default: info)

#### Frontend Configuration

- `FRONTEND_PORT`: Frontend server port (default: 3000)
- `FRONTEND_URL`: Frontend URL for CORS (default: http://localhost:3000)
- `BACKEND_URL`: Backend API URL (default: http://localhost:6050)

#### Registration Authorization

- `AUTHORIZED_EMAILS`: Comma-separated list of email addresses authorized to create accounts (e.g., "admin@example.com,user1@example.com")

#### Optional Services

- `REDIS_PORT`: Redis port (default: 6379)

## 🐳 Docker Deployment

### Development

```bash
# Start all services
docker-compose up

# Start in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production

```bash
# Build and start in production mode
docker-compose -f docker-compose.yml up -d

# Update and restart
docker-compose pull
docker-compose up -d --build
```

## 🛠️ Development

### Local Development Setup

#### Backend Development

```bash
cd backend
go mod download
go run main.go
```

#### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

### Database Setup

The database will be automatically initialized with the required schema when you first run the application using Docker Compose.

## 📁 Project Structure

```
booksrendezvous/
├── frontend/                 # Nuxt.js frontend application
│   ├── components/          # Vue components
│   ├── pages/              # Application pages
│   ├── plugins/            # Nuxt plugins
│   ├── stores/             # Pinia stores
│   ├── package.json        # Frontend dependencies
│   ├── nuxt.config.ts      # Nuxt configuration
│   └── Dockerfile          # Frontend container
├── backend/                 # Go backend application
│   ├── controllers/        # API controllers
│   ├── database/           # Database configuration
│   ├── models/             # Data models
│   ├── routes/             # API routes
│   ├── services/           # Business logic
│   ├── utils/              # Utility functions
│   ├── data/               # Static data files
│   ├── main.go             # Application entry point
│   ├── go.mod              # Go dependencies
│   └── Dockerfile          # Backend container
├── docker-compose.yml       # Docker services configuration
├── .env.example            # Environment variables template
└── README.md               # This file
```

## 🔧 API Documentation

The backend provides a REST API with the following main endpoints:

- `GET /health` - Health check
- `GET /api/books` - Get all books
- `POST /api/books` - Create a new book
- `PUT /api/books/:id` - Update a book
- `DELETE /api/books/:id` - Delete a book
- `GET /api/achievements` - Get achievements

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🙏 Credits

This project was inspired by and built upon the work of [@yoanbernabeu](https://github.com/yoanbernabeu), particularly the [lectures repository](https://github.com/yoanbernabeu/lectures). The overall architecture and approach were influenced by his excellent work on book and reading management applications.

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [Issues](../../issues) page
2. Create a new issue with detailed information
3. Include logs and configuration details

## 🔄 Updates

To update your deployment:

```bash
# Pull latest changes
git pull origin main

# Rebuild and restart containers
docker-compose down
docker-compose up -d --build
```

## 🏷️ Version

Current version: 1.0.0

---

Made with ❤️ for book lovers everywhere!

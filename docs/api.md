# Backend Architecture

The Python API follows a layered architecture to separate HTTP handling, business logic, and data access.

Request Flow:

HTTP Request
→ Router
→ Service
→ Repository
→ SQLAlchemy
→ PostgreSQL

## Layers

### Router Layer

Responsibilities:

* Define API endpoints
* Validate requests
* Return responses

Routers should not contain business logic or database queries.

Example:

Router
→ MovieService

---

### Service Layer

Responsibilities:

* Business rules
* Validation
* Workflow orchestration
* Cross-domain operations

Services coordinate repositories and other services.

Example:

MovieService
→ MovieRepository

---

### Repository Layer

Responsibilities:

* Database access
* ORM queries
* Persistence operations

Repositories should not contain business rules.

Example:

MovieRepository
→ PostgreSQL

---

# Project Structure

api/

├── app/
│
├── config.py
├── database.py
├── dependencies.py
│
├── models/
│   └── movie.py
│
├── schemas/
│   └── movie.py
│
├── repositories/
│   └── movie_repository.py
│
├── services/
│   └── movie_service.py
│
├── routers/
│   └── movie_router.py
│
└── main.py

---

# Data Ownership

## Python API

Owns:

* Movie metadata
* Users (future)
* Watch history (future)
* Search (future)

Persists data in PostgreSQL.

## Go Streaming Service

Owns:

* Video delivery
* HTTP range requests
* Media file access

Does not access PostgreSQL.

Video files are read directly from the shared media storage.

---

# Database Strategy

PostgreSQL stores metadata only.

Examples:

* Movie title
* Description
* Filename
* Duration (future)
* Thumbnail path (future)

Video files are not stored in PostgreSQL.

Video files are stored on the filesystem under:

media/

Example:

media/
├── interstellar.mp4
├── batman.mp4
└── matrix.mp4

The database stores only the filename reference.

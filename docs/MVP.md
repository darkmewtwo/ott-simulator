# OTT Simulator - MVP Specification

## Objective

Build a minimal OTT platform that allows users to browse available movies and stream video content.

The goal is to learn backend architecture, service decomposition, database design, video streaming concepts, Docker, PostgreSQL, Go, and Python.

The platform will run entirely on a local machine using Docker containers.

---

# MVP Scope

## User Features

### List Movies

Users can view a list of available movies.

### View Movie Details

Users can view metadata about a movie.

### Stream Movie

Users can play a movie from a browser.

---

# Out of Scope

The following features are intentionally excluded from the MVP:

* User Authentication
* Authorization
* Search
* Categories
* Watch History
* Resume Playback
* Recommendations
* CDN
* Adaptive Bitrate Streaming
* Upload API
* Analytics
* Payments
* Subscriptions

---

# Architecture

## Python Metadata Service

Responsibilities:

* Manage movie metadata
* Expose movie APIs
* Communicate with PostgreSQL

Technology:

* FastAPI
* SQLAlchemy
* PostgreSQL

Endpoints:

GET /movies

GET /movies/{id}

---

## Go Streaming Service

Responsibilities:

* Stream video files
* Handle HTTP range requests
* Serve media content

Technology:

* Go
* net/http

Endpoints:

GET /stream/{filename}

---

## PostgreSQL

Stores:

Movie metadata only.

Video files are not stored in PostgreSQL.

---

## Filesystem Storage

Video files are stored on disk.

Example:

media/
├── interstellar.mp4
├── batman.mp4
└── matrix.mp4

---

# Movie Data Model

Movie

* id
* title
* description
* filename
* created_at

Example:

Movie
{
"id": 1,
"title": "Interstellar",
"filename": "interstellar.mp4"
}

---

# Request Flow

Movie Listing

Browser
→ Python API
→ PostgreSQL

Movie Playback

Browser
→ Python API
→ Movie Metadata

Browser
→ Go Streamer
→ Filesystem

---

# MVP Success Criteria

The MVP is complete when:

1. Docker Compose starts all services.
2. PostgreSQL is running.
3. Movie metadata is stored in PostgreSQL.
4. Movie listing API returns data.
5. Browser displays available movies.
6. Browser can play a movie using the Go streaming service.

---

# Future Versions

V2

* User accounts
* Login

V3

* Watch history
* Resume playback

V4

* Upload API
* Thumbnail generation

V5

* Adaptive streaming
* Transcoding service

V6

* CDN simulation
* Multiple client simulation

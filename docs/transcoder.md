# OTT Transcoder Service

## Overview

The Transcoder Service is responsible for processing uploaded video files and converting them into streaming-friendly formats.

This service operates independently from the API and Streaming services.

### Responsibilities

* Poll for movies awaiting processing.
* Extract video metadata.
* Calculate and store video duration.
* Generate HLS playlists and segments.
* Update movie processing status.
* Handle transcoding failures.
* Generate additional media assets in the future.

### Non-Responsibilities

The Transcoder Service does NOT:

* Authenticate users.
* Serve video files.
* Manage movie metadata.
* Handle user watch events.
* Generate recommendations.

These concerns belong to other services in the OTT ecosystem.

---

# Architecture

```text
FastAPI
    │
    │ Upload MP4
    ▼
Shared Storage
    │
    ▼
Transcoder Worker
    │
    ├── ffprobe
    ├── ffmpeg
    ├── duration extraction
    └── HLS generation
    │
    ▼
PostgreSQL

Go Streamer
    │
    ▼
Serve HLS files
```

---

# Movie Lifecycle

```text
PENDING
    │
    ▼
PROCESSING
    │
    ├── Extract duration
    ├── Generate HLS
    │
    ▼
READY
```

Failure path:

```text
PENDING
    │
    ▼
PROCESSING
    │
    ▼
FAILED
```

---

# Project Structure

```text
transcoder/
├── cmd/
│   └── main.go
│
├── internal/
│   ├── db/
│   ├── models/
│   ├── repository/
│   ├── worker/
│   └── ffmpeg/
│
├── go.mod
├── Dockerfile
└── README.md
```

---

# Entry Point

File:

```text
cmd/main.go
```

Package:

```go
package main
```

Responsibilities:

* Load configuration.
* Establish database connection.
* Create worker instance.
* Start polling loop.
* Handle service shutdown.

Business logic should not be implemented inside `main.go`.

---

# Package Responsibilities

## internal/db

Responsible for:

* Database connection.
* Connection pool configuration.
* Database initialization.

Example:

```go
db.Connect()
```

---

## internal/models

Contains database models and enums.

Examples:

```go
Movie
MovieStatus
```

---

## internal/repository

Responsible for database operations.

Examples:

```go
GetPendingMovie()
UpdateMovieStatus()
UpdateMovieDuration()
```

This package isolates database access from business logic.

---

## internal/worker

Contains the transcoding workflow.

Responsibilities:

* Poll pending movies.
* Claim work.
* Execute processing pipeline.
* Handle failures.

Workflow:

```text
PENDING
    ↓
PROCESSING
    ↓
READY
```

---

## internal/ffmpeg

Encapsulates all interactions with FFmpeg and FFprobe.

Examples:

```go
GetDuration()
GenerateHLS()
```

Only this package should directly execute external FFmpeg commands.

---

# Storage Layout

```text
storage/
├── uploads/
│   └── movie.mp4
│
└── hls/
    └── <movie_id>/
        ├── master.m3u8
        ├── segment000.ts
        ├── segment001.ts
        └── ...
```

The movie ID is used as the HLS output directory name.

---

# Future Enhancements

Potential future responsibilities:

* Multi-bitrate transcoding.
* Thumbnail generation.
* Preview clip generation.
* Video quality analysis.
* AI-powered content analysis.
* Distributed worker scaling.
* Queue-based processing.

The current design should support these additions without major restructuring.

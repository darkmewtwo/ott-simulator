# OTT Streaming Platform Simulator

A modular OTT (Over-The-Top) streaming platform simulator that recreates the backend architecture of modern video streaming services. The project combines multiple services written in Go and Python to simulate user authentication, movie management, HLS video streaming, transcoding, watch progress tracking, and autonomous user behavior.

The goal is not just to stream videos, but to model the complete lifecycle of an OTT platform—from video ingestion to playback analytics—using production-inspired architecture.

---

## Features

### User Management

* User registration and authentication
* JWT-based authorization
* Persistent user profiles
* Watch history tracking

### Movie Management

* Upload movies
* Store metadata
* Browse available content
* Retrieve movie details

### Video Processing

* Automatic transcoding using FFmpeg
* HLS playlist generation
* Multi-segment video streaming
* Metadata extraction using FFprobe

### Streaming

* HTTP Live Streaming (HLS)
* Efficient Go-based streaming service
* Range request support
* Playlist and segment delivery

### Watch Progress

* Save playback position
* Resume watching
* Track watch time
* Playback event recording

### Autonomous Client Simulator

* Simulates real users interacting with the platform
* Generates realistic user accounts
* Performs registration and login
* Streams content automatically
* Models user behavior through sessions and capabilities

---

# Architecture

```
                    +-------------------+
                    |   Web / Client    |
                    +---------+---------+
                              |
                              v
                    +-------------------+
                    |   FastAPI Backend |
                    +---------+---------+
                              |
             +----------------+----------------+
             |                                 |
             v                                 v
     PostgreSQL                        File Storage
     Metadata                          Videos / HLS

                              ^
                              |
             +----------------+----------------+
             |                                 |
             |                                 |
     +-------+--------+               +--------+-------+
     | Streaming API  |               | Transcoder     |
     | Go             |               | Go + FFmpeg    |
     +----------------+               +----------------+

                              ^
                              |
                     Simulator Clients
```

---

# Project Structure

```
.
├── backend/           # FastAPI REST API
├── streamer/          # Go streaming service
├── transcoder/        # Go transcoding service
├── simulator/         # Autonomous client simulator
├── frontend/          # Web UI
├── videos/            # Uploaded videos
├── hls/               # Generated HLS assets
├── docker-compose.yml
└── README.md
```

---

# Technology Stack

## Backend

* FastAPI
* SQLAlchemy
* PostgreSQL
* JWT Authentication
* Pydantic

## Streaming

* Go
* HTTP Range Requests
* HLS

## Video Processing

* FFmpeg
* FFprobe

## Simulator

* Go
* Concurrent workers
* Autonomous user sessions
* Modular capabilities

## Infrastructure

* Docker
* Docker Compose

---

# Services

## Backend

Responsible for:

* Authentication
* Movie metadata
* Watch progress
* User management
* REST APIs

---

## Streaming Service

Provides:

* HLS playlists
* Video segments
* Efficient file serving
* Streaming endpoints

---

## Transcoder

Processes uploaded videos into HLS format.

Pipeline:

```
Video Upload
      │
      ▼
FFprobe Metadata
      │
      ▼
FFmpeg Transcoding
      │
      ▼
HLS Playlist
      │
      ▼
Streaming Service
```

---

## Simulator

The simulator behaves like thousands of independent OTT users.

Each simulated user has:

* Identity
* Session
* Authentication capability
* Streaming capability
* Independent behavior

Future work includes personality-driven viewing habits, recommendation interactions, and realistic daily activity patterns.

---

# Running the Project

```bash
git clone <repository>

cd ott-simulator

docker compose up --build
```

Services will be available after startup.

---

# Future Improvements

* Recommendation engine
* Signed HLS URLs
* Adaptive bitrate streaming
* Search service
* Analytics dashboard
* User personalities
* Recommendation algorithms
* Subscription plans
* CDN simulation
* Distributed streaming nodes
* Load testing scenarios

---

# Purpose

This project is designed as a learning and experimentation platform for backend engineering, distributed systems, media streaming, and autonomous client simulation.

It demonstrates how multiple independent services can collaborate to provide the functionality of a modern OTT platform while remaining modular and scalable.

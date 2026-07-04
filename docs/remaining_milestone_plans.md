# OTT Platform Roadmap

This document outlines the next major milestones for evolving the OTT
platform from a functional prototype into a production-grade streaming
system.

------------------------------------------------------------------------

# Current Architecture

``` text
                Upload
                   │
                   ▼
             FastAPI Backend
                   │
                   ▼
             PostgreSQL Database
                   │
                   ▼
        Go Transcoding Worker
                   │
        ┌──────────┴──────────┐
        ▼                     ▼
    FFprobe              FFmpeg HLS
        │                     │
        └──────────┬──────────┘
                   ▼
          /media/hls/{movieId}
                   │
                   ▼
          Go Streaming Service
                   │
             CORS Middleware
                   │
                   ▼
             Video.js Player
                   │
                   ▼
      Watch Progress & Analytics
```

# Upcoming Milestones

## Milestone 8 --- Secure Streaming

**Goal**

Prevent unauthenticated users from accessing HLS content.

### Tasks

-   Add JWT authentication middleware to the Go streaming service.
-   Validate Authorization Bearer tokens.
-   Reject expired or invalid JWTs.
-   Return HTTP 401 for unauthorized requests.
-   Ensure Video.js sends Authorization headers with every HLS request.

### Learning Objectives

-   Go middleware
-   JWT validation
-   Secure media delivery
-   HTTP Authorization

------------------------------------------------------------------------

## Milestone 9 --- Reverse Proxy & Unified Gateway

**Goal**

Serve the frontend, API, and streaming service through a single public
origin.

Current:

``` text
localhost:8080
localhost:8100
localhost:8180
```

Target:

``` text
                https://ott.local
                       │
                       ▼
              Reverse Proxy
            (NGINX / Traefik)
                       │
      ┌────────────────┼────────────────┐
      ▼                ▼                ▼
  Frontend         FastAPI API      Go Streamer
```

### Responsibilities

-   `/` → Frontend
-   `/api/*` → FastAPI
-   `/hls/*` → Go Streamer
-   `/posters/*` → Go Streamer

### Features

-   Request routing
-   HTTPS termination
-   Remove CORS in production
-   Static asset serving
-   Compression (gzip/Brotli)
-   Security headers
-   Centralized logging
-   Foundation for load balancing

### Learning Objectives

-   Reverse proxy architecture
-   NGINX / Traefik
-   HTTP routing
-   TLS termination
-   Production deployment

------------------------------------------------------------------------

## Milestone 10 --- Signed URLs

-   Expiring HLS URLs
-   HMAC signature validation
-   Prevent URL sharing
-   CDN-ready secure delivery

------------------------------------------------------------------------

## Milestone 11 --- Adaptive Bitrate Streaming

-   Multi-quality transcoding
-   Master playlist generation
-   Automatic quality switching

------------------------------------------------------------------------

## Milestone 12 --- Subtitle Support

-   WebVTT
-   Multiple subtitle languages
-   HLS subtitle groups

------------------------------------------------------------------------

## Milestone 13 --- Multiple Audio Tracks

-   Multiple audio languages
-   HLS media groups
-   Runtime audio switching

------------------------------------------------------------------------

## Milestone 14 --- Thumbnail Previews

-   Thumbnail sprites
-   Preview VTT
-   Seek preview

------------------------------------------------------------------------

## Milestone 15 --- CDN Integration

-   Edge caching
-   Cache invalidation
-   Object storage integration

------------------------------------------------------------------------

## Milestone 16 --- Production Hardening

-   Structured logging
-   Request IDs
-   Metrics
-   Health checks
-   Graceful shutdown
-   Rate limiting
-   Configuration management
-   Docker Compose production deployment

------------------------------------------------------------------------

# Long-Term Vision

``` text
                Browser
                    │
                    ▼
              Reverse Proxy
                    │
      ┌─────────────┴─────────────┐
      ▼                           ▼
   Frontend                   FastAPI
                                   │
                                   ▼
                            PostgreSQL
                                   │
                                   ▼
                        Go Transcoder Worker
                                   │
                            FFprobe / FFmpeg
                                   │
                                   ▼
                          Object Storage
                                   │
                                   ▼
                            Go Streamer
                                   │
                                   ▼
                                  CDN
                                   │
                                   ▼
                             Video.js Player
```

# Learning Goals

-   Backend API Development
-   Go HTTP Servers
-   Background Workers
-   FFmpeg Transcoding
-   HLS Streaming
-   JWT Authentication
-   Reverse Proxying
-   Secure Media Delivery
-   Adaptive Bitrate Streaming
-   Subtitle & Audio Tracks
-   CDN Integration
-   Production Deployment
-   Streaming Platform Architecture
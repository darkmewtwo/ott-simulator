# OTT Simulator - Autonomous Client Architecture

## Overview

The objective of the OTT Simulator is not merely to generate HTTP traffic, but to simulate realistic users interacting with an OTT platform. Each simulated user behaves autonomously, has persistent identity, individual preferences, and a distinct personality.

The simulator is designed to evolve alongside the OTT platform. As new platform features are introduced, the simulator should require minimal changes while allowing autonomous clients to naturally begin using those new capabilities.

---

# Design Principles

The simulator follows several core principles.

### Autonomous Users

Each client behaves independently after being spawned.

The orchestrator is responsible only for creating and managing clients. It does not control their individual actions during the simulation.

---

### Platform-Driven Simulation

The simulator should adapt to the capabilities of the OTT platform.

Instead of clients having hardcoded knowledge of every platform feature, the orchestrator maintains knowledge of the platform and provides this information to newly created clients.

As the OTT platform evolves, only the orchestrator's platform description and workflows need to be updated.

---

### Separation of Concerns

The simulator separates responsibilities into distinct components.

| Component          | Responsibility                                                      |
| ------------------ | ------------------------------------------------------------------- |
| Orchestrator       | Manages the simulation lifecycle                                    |
| Platform Knowledge | Describes available platform capabilities and interaction workflows |
| Autonomous Client  | Simulates an individual user                                        |
| Personality Model  | Determines user preferences and decision making                     |
| Workflow Library   | Encapsulates platform interactions                                  |
| Metrics Collector  | Aggregates simulation statistics                                    |

---

# High-Level Architecture

```
                    OTT Platform
                          ▲
                          │
                 Platform Workflows
                          ▲
                          │
                  +------------------+
                  |   Orchestrator   |
                  +------------------+
                           │
            Platform Knowledge + Personalities
                           │
      ------------------------------------------------
      │                  │                  │
   Client 1          Client 2          Client N
```

---

# Orchestrator

The orchestrator is the central coordinator of the simulation.

Its responsibilities include:

* Initializing the simulation.
* Maintaining knowledge of platform capabilities.
* Registering available workflows.
* Spawning autonomous clients.
* Generating user identities.
* Assigning personalities.
* Collecting simulation metrics.
* Monitoring client lifecycle.

The orchestrator **does not** make behavioral decisions for clients.

---

# Platform Knowledge

The orchestrator maintains a description of the platform rather than exposing raw HTTP endpoints.

Examples of supported workflows include:

* User Registration
* Login
* Browse Catalog
* View Movie
* Playback
* Continue Watching

As additional platform functionality is implemented, corresponding workflows are added without modifying client behavior.

---

# Workflow Library

A workflow represents a complete interaction with the platform.

Examples include:

* Registration Workflow
* Login Workflow
* Movie Browsing Workflow
* Playback Workflow
* Search Workflow (future)
* Recommendation Workflow (future)

Each workflow contains the implementation details necessary to interact with the platform.

Examples include:

* HTTP endpoints
* Request payloads
* Authentication handling
* Error recovery
* Response parsing

This keeps clients independent of API implementation details.

---

# Autonomous Client

Each client represents an independent OTT user.

A client owns:

* Identity
* Authentication state
* Personality
* Viewing history
* Preferences
* Random number generator
* Internal state

Once initialized, the client operates independently.

---

# Client Lifecycle

Initial implementation:

```
NEW

↓

REGISTER

↓

LOGIN

↓

READY
```

Future lifecycle:

```
READY

↓

BROWSE

↓

VIEW MOVIE

↓

PLAY

↓

PAUSE

↓

SEEK

↓

STOP

↓

IDLE
```

---

# Identity

Each client is assigned a unique identity.

The identity contains:

* Email
* Password
* Display Name

Future additions may include:

* Region
* Subscription Plan
* Age Group
* Device Type

---

# Personality

Personality determines how a user behaves.

Examples include:

* Casual Viewer
* Binge Watcher
* Explorer
* Documentary Fan
* Action Enthusiast
* Kids Profile
* Weekend Viewer

Personality influences decision making but does not contain platform implementation details.

---

# Client Decision Making

Clients make decisions based on two independent inputs.

**Personality**

Determines what the user would like to do.

Examples:

* Browse movies
* Watch another episode
* Leave the platform
* Resume playback

**Platform Knowledge**

Determines whether the desired action is supported.

If a feature does not exist, the client naturally chooses an alternative behavior.

---

# Interaction Strategy

The simulator will primarily interact with the OTT platform through backend APIs.

Advantages include:

* High scalability
* Low resource consumption
* Deterministic execution
* Large numbers of concurrent users

A secondary browser-based validation layer may be introduced using Playwright.

Browser clients will verify:

* User interface functionality
* JavaScript execution
* Page navigation
* Video playback
* Cookie handling
* Frontend regressions

API-driven clients remain the primary simulation mechanism.

---

# Simulation Flow

```
Initialize Orchestrator

↓

Load Platform Knowledge

↓

Register Available Workflows

↓

Spawn Clients

↓

Generate Identity

↓

Assign Personality

↓

Client Executes Registration Workflow

↓

Client Executes Login Workflow

↓

Client Enters Autonomous Execution
```

---

# Initial Milestone

The first milestone focuses exclusively on identity creation.

Objectives:

* Implement the orchestrator.
* Generate unique user identities.
* Spawn autonomous clients.
* Execute user registration.
* Execute login.
* Persist authentication state.
* Report metrics.

No browsing or playback functionality is included in this milestone.

---

# Long-Term Vision

The simulator should evolve together with the OTT platform.

Whenever a new platform feature is implemented:

1. Add a corresponding workflow.
2. Register it with the orchestrator.
3. Update platform knowledge.
4. Existing personalities naturally begin using the new capability when appropriate.

This architecture minimizes changes to client logic while allowing the simulator to grow alongside the platform, providing increasingly realistic autonomous user behavior over time.

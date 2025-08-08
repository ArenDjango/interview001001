**Event Model**
**Each event should include:**

`id (UUID, auto-generated)
title (string, required, max 100 characters)
description (string, optional)
start_time and end_time (timestamps; start_time must be before end_time)
created_at (timestamp, auto-set on creation)`

**API Endpoints:**

**POST /events**

Accepts JSON input
Validates required fields and time constraints
Saves to database and returns created event with HTTP 201

**GET /events**

Returns all events ordered by start_time ascending

**GET /events/{id}**

Returns a
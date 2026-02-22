# API Testing Guide

This document contains precise `curl` commands to simulate a complete end-to-end flow of the Secure Online Quiz Backend. This allows you to test Role-Based Access Control, session management, and the anti-cheating mechanisms.

Ensure your server is running (e.g. `docker-compose up --build` or `make run`) on `http://localhost:8080`.

## Full E2E Simulation Flow

### Step 1: Register and Login as a Teacher

```bash
# Register a teacher
curl -X POST http://localhost:8080/api/auth/register \
-H "Content-Type: application/json" \
-d '{"username": "prof_smith", "password": "password123", "role": "teacher"}'

# Login to get the JWT
curl -X POST http://localhost:8080/api/auth/login \
-H "Content-Type: application/json" \
-d '{"username": "prof_smith", "password": "password123"}'
```
*Note: Copy the `token` from the response. This is your `<TEACHER_TOKEN>`.*

### Step 2: Create a Quiz (Teacher)

```bash
# Create a quiz using the teacher token
curl -X POST http://localhost:8080/api/quizzes \
-H "Authorization: Bearer <TEACHER_TOKEN>" \
-H "Content-Type: application/json" \
-d '{
  "title": "Midterm Exam",
  "start_time": "2026-02-21T00:00:00Z",
  "end_time": "2026-12-31T23:59:59Z",
  "published": true
}'
```
*Note: Let's assume the created quiz has ID `1`.*

### Step 3: Register and Login as a Student

```bash
# Register a student
curl -X POST http://localhost:8080/api/auth/register \
-H "Content-Type: application/json" \
-d '{"username": "alice", "password": "password123", "role": "student"}'

# Login to get the student JWT
curl -X POST http://localhost:8080/api/auth/login \
-H "Content-Type: application/json" \
-d '{"username": "alice", "password": "password123"}'
```
*Note: Copy the new `token` from the response. This is your `<STUDENT_TOKEN>`.*

### Step 4: Start the Quiz Session (Student)

```bash
# Start a session for Quiz ID 1
curl -X POST http://localhost:8080/api/sessions/start \
-H "Authorization: Bearer <STUDENT_TOKEN>" \
-H "Content-Type: application/json" \
-d '{"quiz_id": 1}'
```
*Note: This returns a `session_id`. Let's assume it is `1`.*

### Step 5: Simulate Cheating (Tab-Switching)

Simulate the student leaving the browser tab multiple times. Run this command **6 times** to exceed the suspicious score threshold (which is set to 5):

```bash
curl -X POST http://localhost:8080/api/sessions/1/tab-switch \
-H "Authorization: Bearer <STUDENT_TOKEN>"
```

### Step 6: Submit the Quiz

```bash
curl -X POST http://localhost:8080/api/sessions/1/submit \
-H "Authorization: Bearer <STUDENT_TOKEN>"
```

**Expected Result:** Because the tab-switch endpoint was hit 6 times, the submit response will enforce the anti-cheating policy: The score will be deducted, and instead of returning `state: "submitted"`, the API will return `state: "under_review"`. This successfully demonstrates the strict server-side state enforcement.

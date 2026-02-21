# Secure Online Quiz Backend

This repository contains the backend for a Secure Online Quiz platform, built in Go using the Gin framework and PostgreSQL.

## Architecture Decisions

The system is built using **Clean Architecture** principles, prioritizing separation of concerns, testability, and maintainability.

- **cmd/**: Contains the main application entry point and wiring (`main.go`).
- **internal/domain/**: Houses core business entities (Models) such as User, Quiz, Session, and SessionQuestion, along with constants like `Role` and `SessionState`. These models have no dependencies on outer layers.
- **internal/repository/**: Implementations of data access interfaces using `gorm`. Interactions with PostgreSQL are isolated here (e.g. `UserRepository`, `QuizRepository`, `SessionRepository`).
- **internal/service/**: Contains the core business logic. Services coordinate between repositories. Includes authentication, quiz management, session tracking, and the anti-cheating state logic.
- **internal/handler/**: HTTP transport layer using Gin. Handles deserialization, invokes services, and formats JSON responses.
- **internal/middleware/**: Reusable HTTP middleware for verifying JWT tokens and enforcing Role-Based Access Control (RBAC).

## Anti-Cheating Design Choices

The platform implements realistic constraints to deter and detect cheating.
1. **Limited Attempts & Expiration Verification**: A student can only have one session per quiz. A strict state-machine (`not_started` -> `active` -> `submitted` -> `expired`/`under_review`) prevents invalid lifecycle transitions or duplicate attempts.
2. **Server-Side Scoring & Question Shuffling**: To prevent students from sharing answers based on question order, questions and their multiple-choice options are randomly shuffled per-session and mapped to a specific `SessionQuestion` entity. Scoring is exclusively evaluated server-side.
3. **Frontend Tab-Switch Reporting**: The API exposes a `POST /sessions/tab-switch` endpoint that the frontend calls whenever the browser loses focus (e.g., via `visibilitychange` event). 
4. **Scoring Deductions & Flagging**: Each registered tab switch increments the session's `TabSwitchCount` and `SuspiciousScore`. When the quiz is submitted:
   - 1 mark is deducted from the final score per tab switch (enforcing a minimum score of 0).
   - If `SuspiciousScore` exceeds a hardcoded threshold (e.g., 5), the session is placed in an `under_review` state rather than `submitted`, requiring teacher intervention.

## Database Design
The schema uses strong foreign keys and constraints. E.g., `QuizSession` has a unique constraint on `(quiz_id, student_id)` to enforce the single-attempt rule. Complex operations, such as submitting a quiz, are executed within Database Transactions using GORM to ensure atomicity and consistency.

## Running the Application

Prerequisites: Docker and docker-compose.

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8080`.

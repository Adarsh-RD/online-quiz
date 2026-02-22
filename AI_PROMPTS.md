# AI Prompt Engineering Log

The following provides a log of the key architectural directives and structured prompts I used alongside the AI coding assistant to generate the `online-quiz` backend. 

The prompts were explicitly engineered to enforce strict enterprise standards, Golang best practices, and a clear Clean Architecture structure, ensuring the AI acted as a code generator for my predetermined system design.

### 1. Initial Architecture & Project Scaffolding
> "Initialize a new Go module named `online-quiz`. Scaffold the project strictly using Clean Architecture principles. Create the directory layers: `cmd/server/`, `internal/domain/`, `internal/repository/`, `internal/service/`, `internal/handler/`, and `internal/middleware/`. We will use the `gin-gonic/gin` framework for routing and `gorm.io/gorm` with PostgreSQL for the database implementation. Do not write the implementation yet; just set up the module hierarchy and directory tree."

### 2. Domain Modeling & RBAC Design
> "Let's define our core generic domain models in `internal/domain/`. Create a `User` struct with standard fields (ID, username, password_hash) and an enum for `Role` containing `student` and `teacher`. Next, create a `Quiz` model referencing `teacher_id` with `start_time` and `end_time` temporal bounds. Ensure GORM tags are correctly applied for primary keys and foreign key constraints."

### 3. State Machine & Anti-Cheat Session Modeling
> "I need to securely track student attempts. Define a `QuizSession` domain model with a composite unique-index on `(quiz_id, student_id)` to enforce a strict one-attempt limit at the database level. It needs a state machine enum: `not_started`, `active`, `submitted`, `expired`, and `under_review`. Also, include anti-cheating tracking fields: `tab_switch_count` (int) and `suspicious_score` (int)."

### 4. JWT & Middleware Implementation
> "Write the JWT generation logic in `internal/service/jwt.go` utilizing `golang-jwt`. The claims must embed the user's ID, username, and Role. Afterward, implement two Gin middlewares in `internal/middleware/auth.go`: one to authenticate the JWT and attach the claims to the Gin context, and another factory middleware `RequireRole(roles ...domain.Role)` to enforce Role-Based Access Control iteratively on our upcoming endpoints."

### 5. Repository Interfaces & Context Propagation
> "Implement the Repository layer interfaces. I want `UserRepository`, `QuizRepository`, and `SessionRepository`. It is critically important that every single repository method accepts a `context.Context` as its first argument. This is to allow for upstream HTTP request cancellation and tracing down to the database connection. Use GORM's `.WithContext(ctx)` method exclusively."

### 6. Shuffling Strategy & Isolation Logic
> "In `internal/service/session.go`, implement a `StartSession` function. When a student starts a session, pull the quiz's questions, shuffle the specific question order, and shuffle each question's multiple-choice options using a pseudo-random seed locally. Persist this specific mapping into `session_questions` and `session_options` proxy tables so each student gets a unique, deterministic version of the test. This prevents parallel answer sharing."

### 7. Absolute Temporal Enforcement & Transactions
> "Implement the `SubmitQuiz` service logic. It needs to calculate the score entirely server-side based only on the `option_id` submitted by the client. Before allowing the submission, validate that the current absolute UTC server time falls exactly between the `quiz.start_time` and `quiz.end_time`. If the session's `suspicious_score` is >= 5, transition the state to `under_review` instead of `submitted`. Wrap this entire scoring and submission process inside a GORM SQL transaction for atomicity."

### 8. Handling Browser Events (Anti-Cheat)
> "Create a Gin handler and service endpoint for `POST /api/sessions/:id/tab-switch`. This will be called asynchronously by the frontend's `visibilitychange` event listener. When hit, increment the `tab_switch_count` and `suspicious_score` by 1. Keep the context propagation strict throughout this call."

### 9. Structured Logging & Graceful Shutdown
> "Let's elevate this codebase to handle enterprise DPI standards. Replace all standard library `log` implementations with Uber Zap (`go.uber.org/zap`). Initialize a global structured JSON logger. Then, in `cmd/server/main.go`, implement a graceful shutdown mechanism catching `syscall.SIGINT` and `syscall.SIGTERM`. Spin the Gin server off into a dedicated goroutine and use `srv.Shutdown(ctx)` with a 5-second timeout on exit to prevent abruptly dropping HTTP connections."

### 10. QA Tooling & Mermaid Diagrams
> "Finally, ensure our documentation is perfect for the repository reviewers. Add a `Makefile` with targets for `build`, `run`, `test`, `lint`, and Docker commands. Ensure we have a strict `.golangci.yml` configuration turning on cyclomatic complexity and unhandled error checks. In the `README.md`, map out the Clean Architecture boundaries, the ER Database models, and the QuizSession State Machine using `mermaid` syntax."

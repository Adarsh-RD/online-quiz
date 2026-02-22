# AI Assistant Prompts Used

This document outlines the main prompts I used to build this project with the help of an AI coding assistant. As a student learning Go, I understood the core components needed for a quiz backend (like databases, APIs, and basic security) and used the AI to help me quickly implement the Go code and structure it properly.

### 1. Initial Setup and Architecture
> "I need to build a secure online quiz backend in Go for my class project. I know the basics of Go and how to use the Gin framework, but I want to use 'Clean Architecture'. Can you help me set up the initial project folders like handlers, services, and repositories?"

### 2. Database Models
> "We need to store data in a PostgreSQL database. Let's use GORM. Can you help me write the Go structs/models for a User (could be a student or a teacher), a Quiz (with start and end times), Questions, and a QuizSession to track student attempts?"

### 3. Implementing Authentication
> "I need a way for users to register and log in. I've heard JWT is a good way to handle this. Can you help me write a Gin middleware that uses JWT to authenticate students and teachers, so only teachers can create quizzes?"

### 4. Preventing Multiple Attempts
> "A student should only be able to take a specific quiz one time. How can we make sure of this at the database level? Can we use a unique constraint on the student ID and quiz ID in the QuizSession table?"

### 5. Shuffling Questions
> "To stop students from cheating by sitting next to each other, I want every student to get the questions and options in a random order. Can you write a function in the session service that shuffles the questions when they start the quiz and saves that order to the database?"

### 6. Server-Side Timer Logic
> "How do we handle the quiz timer securely? If we just use a Javascript timer, a student could hack it. Can we make it so the Go server checks the current UTC time against the quiz's `end_time` when they try to submit?"

### 7. Tab-Switching Anti-Cheat
> "I want to add another anti-cheating feature. If the frontend detects that a student switched browser tabs (like to Google an answer), it should send a request to the server. Can you make an endpoint that deducts 1 mark from their score every time this happens?"

### 8. Database Transactions
> "When a student submits a quiz, we have to update their score and change the session state to 'submitted'. If something crashes during this, the data might be corrupted. Can you show me how to wrap the submission logic in a GORM transaction?"

### 9. Adding Good Engineering Practices
> "My professors want to see 'enterprise-level' practices in the code. I read that passing `context.Context` everywhere is standard in Go. Can we update our repository methods to use contexts for database queries? Also, can we replace the basic log prints with the `zap` structured logger?"

### 10. Docker Setup
> "Finally, I need to make sure this is easy for anyone to run without installing Go and Postgres manually. Can you write a Dockerfile for the Go app and a docker-compose.yml file that sets up both the app and the Postgres database together?"

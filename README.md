# API Overview (Video)


https://github.com/user-attachments/assets/5a7d81e5-9a2e-4928-946c-e2662ec7acca


# Setup Process

#### Clone the repo
```bash
git clone https://github.com/rudransh-shrivastava/rudransh-backend-task.git
```

#### Firebase Config
Ensure you have a `key.json` file in the root of the project.
The `key.json` is the Firebase service account key which can be obtained from the firebase console

#### Build & Run docker image
```bash
sudo docker-compose up --build
```

# Project Structure Overview

- ├── Dockerfile           (Contains the instructions to dockerise the api server)
- ├── Makefile             (Define commands like make run)
- ├── bin                 (Static Binaries are compiled here)
- │   └── server    
- ├── cmd                 (The main entry point of our application)
- │   └── server
- │       └── main.go
- ├── db.sqlite3          (The database file)
- ├── go.mod
- ├── go.sum
- ├── internal            
- │   ├── api             (The api handlers, and the api server)
- │   │   ├── api.go
- │   │   ├── auth.go
- │   │   ├── handler.go
- │   │   ├── handler_test.go
- │   │   ├── middleware.go
- │   │   └── mockquiz.go
- │   ├── config              (Handles the environment variables if any)
- │   │   └── config.go
- │   ├── db                 (The configuration of our database)
- │   │   └── db.go
- │   ├── schema             (The schemas of our database is defined here)
- │   │   └── schema.go
- │   ├── store              (Store handles communication with the database)
- │   │   ├── course.go
- │   │   ├── quiz.go
- │   │   ├── store.go
- │   │   └── user.go
- │   └── utils
- │       ├── logger
- │       │   └── logger.go   (The logger configurations)
- │       └── utils.go        (Reusable utilities that we use in our code)
- └── key.json                (Firebase credentials obtained from firebase)

# Database Schema and API Overview

## 1. User Registration

**Endpoint:** `POST /api/v1/register`  
**Description:** Registers a new user.  

### Request Body (JSON):
```json
{
  "email": "user@example.com",
  "password": "your_password",
  "name": "Full Name"
}
```

### Response:
Returns the created user object (excluding sensitive details).

**NOTE:**
Please make sure to use the `Authorization` Header for the following requests with value set to `Bearer <token>`
To obtain the token, use the endpoint provided by google.
```
https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=<YOUR_API_KEY>
```
with the request body 
```json
{
    "email": "user@example.com",
    "password": "your_password",
     "returnSecureToken": true
}
```
and copy the `idToken`

---

## 2. Get Courses

**Endpoint:** `GET /api/v1/courses`  
**Roles Allowed:** `STUDENT`, `EDUCATOR`, `ADMIN`  
**Description:** Retrieves a list of courses.  

### Query Parameters:
- `limit` (optional): Number of courses to return (default is 10).
- `offset` (optional): Starting position in the list.

### Response (JSON):
```json
[
  {
    "id": 1,
    "title": "Course Title",
    "user": { ... },
    "user_id": 123,
    "created_at": "2023-03-15T10:00:00Z"
  }
]
```

---

## 3. Post (Create) Course

**Endpoint:** `POST /api/v1/courses`  
**Roles Allowed:** `EDUCATOR`, `ADMIN`  
**Description:** Creates a new course.  

### Request Body (JSON):
```json
{
  "title": "New Course Title"
}
```

### Notes:
- The server extracts the authenticated user (via the context and user store) to set the course creator.

### Response:
Returns the created course object.

---

## 4. Delete Course

**Endpoint:** `DELETE /api/v1/courses`  
**Roles Allowed:** `EDUCATOR`, `ADMIN`  
**Description:** Deletes an existing course.  

### Request Body (JSON):
```json
{
  "id": 1
}
```

### Response:
Returns a confirmation message or the deleted course object.

---

## 5. Generate Quiz

**Endpoint:** `POST /api/v1/quiz/generate`  
**Roles Allowed:** `EDUCATOR`, `ADMIN`  
**Description:** Generates a quiz for a course. The quiz questions are chosen from a predefined pool.  

### Request Body (JSON):
```json
{
  "course_id": "1",
  "number": "3"
}
```

### Response:
Returns the created quiz object with a JSON-encoded string of questions.  

#### Example:
```json
{
  "id": 5,
  "questions": "[{\"question\":\"What is the capital of France?\", \"options\":[\"Berlin\",\"Madrid\",\"Paris\",\"Rome\"], \"answer\":\"Paris\"}]",
  "course": { ... },
  "course_id": 1,
  "created_at": "2023-03-15T10:00:00Z"
}
```

---

## 6. Get Quiz

**Endpoint:** `GET /api/v1/quiz`  
**Roles Allowed:** `EDUCATOR`, `ADMIN` *(as per your code; you might adjust if students should also take quizzes)*  
**Description:** Retrieves a specific quiz by its ID for a given course.  

### Query Parameters:
- `course_id` (required): ID of the course.
- `quiz_id` (required): ID of the quiz.

### Response:
Returns the quiz object. Additionally, the endpoint records that the user has taken the quiz (via the `QuizzesTaken` record).

# How to run tests?
To run the tests, please run the following command.
```bash
go test ./... -v
```

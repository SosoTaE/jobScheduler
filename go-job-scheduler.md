# **Go Job Scheduler**

A backend job scheduling system built with Go, Fiber, and GORM. This application provides a RESTful API to create, manage, and monitor jobs that can execute shell commands based on a flexible, user-defined schedule. It features a background worker pool for concurrent job processing, user authentication with role-based access, and detailed execution history.

## **Key Features**

* **User Authentication**: Secure login/logout functionality with session management.  
* **Role-Based Access Control**: Differentiates between regular users and admin users, with admins having special privileges like registering new users.  
* **CRUD for Jobs**: Full Create, Read, Update, and Delete operations for jobs via the API.  
* **Flexible Scheduling**: A powerful scheduling system allowing jobs to be run at specific times and on specific dates, including:  
  * Years  
  * Months  
  * Days of the month  
  * Days of the week (e.g., Monday, Tuesday)  
* **Concurrent Job Execution**: A robust background worker pool processes jobs from a queue, ensuring non-blocking and efficient execution.  
* **Execution History**: Automatically records the outcome (success/failure), output, and timing of every job run.  
* **Paginated API**: List endpoints for jobs and executions are paginated for efficient data handling.  
* **Configuration via .env**: Easy setup and configuration using environment variables.  
* **Structured Logging**: All events are logged to app.log in JSON format for easy parsing and monitoring.

## **Technologies Used**

* **Backend**: Go  
* **Web Framework**: [Fiber](https://gofiber.io/)  
* **ORM**: [GORM](https://gorm.io/)  
* **Database**: [SQLite](https://www.sqlite.org/index.html)  
* **Configuration**: [godotenv](https://github.com/joho/godotenv)  
* **Authentication**: bcrypt for password hashing, Fiber's session middleware.

## **Getting Started**

Follow these instructions to get the project up and running on your local machine.

### **Prerequisites**

* [Go](https://go.dev/doc/install) (version 1.18 or higher)

### **Installation**

1. **Clone the repository:**  
   git clone https://github.com/SosoTaE/jobScheduler.git  
   cd jobScheduler

2. Install dependencies:  
   This project uses Go Modules. Dependencies will be automatically downloaded when you build or run the project. You can also install them manually:  
   go mod tidy

3. Configure Environment Variables:  
   Create a .env file in the root of the project directory. This file is used to configure the application.  
   touch .env

   Open the .env file and add the following required variables. These will be used to create the initial admin user on the first run.  
   \# Admin User Credentials (Required)  
   ADMIN\_USERNAME=admin  
   ADMIN\_PASSWORD=your-secure-password

   \# Worker Configuration (Optional \- Defaults are used if not set)  
   WORKERS=5  
   QUEUE\_SIZE=100  
   **Note:** The ADMIN\_PASSWORD has a typo in the provided source code (os.Getenv("ADMIN\_PASSWORD") is used for both username and password). For it to work as intended, the .env should be:  
   ADMIN\_PASSWORD=your-secure-password

   The username will be hardcoded as admin by the SeedAdminUser function.  
4. **Run the application:**  
   go run main.go

   The server will start on http://localhost:3000. You will see log messages in your console and in the app.log file indicating that the database connection was successful and the worker pool has started.

## **API Endpoints**

All endpoints are prefixed with /api. An authentication session is required for all routes except /api/login.

| Endpoint | Method | Description | Authentication | Admin Only |
| :---- | :---- | :---- | :---- | :---- |
| /login | POST | Authenticates a user and creates a session. | No | No |
| /logout | POST | Logs out the user and destroys the session. | Yes | No |
| /register | POST | Registers a new user. | Yes | **Yes** |
| /profile | GET | Retrieves the current user's profile. | Yes | No |
| /users | GET | Lists all registered users. | Yes | **Yes** |
| /create/job | POST | Creates a new job. | Yes | No |
| /update/job | PUT | Updates an existing job by id. | Yes | No |
| /delete/job | DELETE | Deletes a job by id. | Yes | No |
| /jobs | GET | Lists all jobs with pagination. Can be filtered by userID. | Yes | No |
| /job/:id | GET | Retrieves the details of a single job. | Yes | No |
| /job/:id/history | GET | Lists the execution history for a specific job. | Yes | No |
| /executions | GET | Lists all job executions across all jobs. | Yes | No |

### **Example API Usage**

#### **Login**

curl \-X POST http://localhost:3000/api/login \\  
\-H "Content-Type: application/json" \\  
\-d '{"username": "admin", "password": "your-secure-password"}' \\  
\-c cookie.txt

#### **Create a Job**

This example creates a job that runs echo "Hello World" twice a day.

curl \-X POST http://localhost:3000/api/create/job \\  
\-H "Content-Type: application/json" \\  
\-b cookie.txt \\  
\-d '{  
    "name": "Hello World Job",  
    "command": "echo \\"Hello World from Job\\"",  
    "schedule": {  
        "times": \[{"hour": 0, "minute": 0}, {"hour": 12, "minute": 30}\]  
    }  
}'

#### **schedule Object Structure**

The schedule object is highly flexible. Here's a more complex example for a job that runs at 8:00 AM and 6:00 PM on Mondays and Fridays in October 2025\.

{  
  "name": "Complex Scheduled Job",  
  "command": "ls \-la",  
  "schedule": {  
    "years": \[2025\],  
    "months": \[10\],  
    "weekdays": \[1, 5\],  
    "times": \[  
      { "hour": 8, "minute": 0 },  
      { "hour": 18, "minute": 0 }  
    \]  
  }  
}

* weekdays: 0 \= Sunday, 1 \= Monday, ..., 6 \= Saturday.  
* months: 1 \= January, ..., 12 \= December.  
* Omitting a field (e.g., daysOfMonth) means the schedule applies to all values for that field.

## **Project Structure**

/  
├── config/           \# Environment variable loading and configuration structs.  
├── handlers/         \# Fiber handlers for authentication and user management.  
│   ├── adminHandler.go \# Logic for seeding the admin user.  
│   └── authHandler.go  \# Logic for login, logout, registration, and auth middleware.  
├── logger/           \# Application-wide structured logger setup.  
├── models/           \# GORM data models for Job and User.  
│   ├── job.go  
│   └── user.go  
├── routes/           \# Fiber handlers for all API endpoints, organized by resource.  
│   ├── createJob.go  
│   ├── deleteJob.go  
│   ├── executionList.go  
│   ├── jobDetail.go  
│   ├── jobHistory.go  
│   ├── jobs.go  
│   ├── profile.go  
│   ├── updateJob.go  
│   └── users.go  
├── scheduler/        \# Core logic to determine if a job is due to run.  
│   └── checker.go  
├── structs/          \# Shared data structures for API requests and responses.  
│   ├── loginRequest.go  
│   └── response.go  
├── worker/           \# Background worker pool, job queue, and scheduler ticker.  
│   └── worker.go  
├── .env              \# Environment variables file (you must create this).  
├── .gitignore        \# Files and directories to be ignored by Git.  
├── app.log           \# JSON log output file (created on first run).  
├── dispatch.db       \# SQLite database file (created on first run).  
├── go.mod            \# Go module definition file.  
└── main.go           \# Application entry point.  

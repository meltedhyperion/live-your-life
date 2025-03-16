# Globetrotter Challenge

Globetrotter Challenge is a full-stack travel guessing game where users receive cryptic clues about famous destinations and must guess the correct location to unlock fun facts and trivia. The application leverages a Golang based backend (using chi) and a React-based frontend and Supabase for authentication and PostgreSQL DB.

## Table of Contents

0. [Tech Stack](#tech-stack)
1. [User Workflow](#user-workflow)
2. [Create Player & Get Player Internals](#create-player--get-player-internals)
3. [Invite System](#invite-system)
4. [Get Leaderboard](#get-leaderboard)
5. [Get Questions](#get-questions)
6. [Check Answer](#check-answer)
7. [Local Setup](#local-setup)
8. [API Documentation](#api-documentation)
9. [How I utilized AI](#how-i-utilized-ai)

---
## 0. Tech Stack

Golang (chi) for backend apis
PostgreSQL (through supabase) for database
Supabase Auth for User Authentication
React (Typescript) and Tailwind for frontend

Using go was a personal choice. I have been using golang for more than 2 years now and it has grown on me. Plus since its a compiled language, its super fast). Chi framework was a design choice as its minimal, mostly built upon the base http package, letting me work with lightweight packages, also a good choice for production.

Using supabase was because of their free teir, and, (whom to lie to), their user auth integration. made it so easy for me to authenticate users. No need to manage JWT creation or storing passowrds and other user sensitive details. Authentication is further carried out by sending email to the registering mail id. Meaning, more auth.
Postgresql because i had very structured queries to make which should be fast enough. And for storing array as values, it mad emore sense.

I am going for a backend specific role and dont have as much experience in developing frontend as i have for backend. Reach and Next are the ones I have worked with and here React felt, simple, easy to use and fast.

---

## 1. User Workflow 

![image](https://github.com/user-attachments/assets/386c830b-274b-4632-af48-33b8c8ede28a)

1. **Authentication:**  
   - Users can sign up or log in using Supabase email authentication.
   - On first login, the backend creates a player profile using the Supabase user ID.
2. **Gameplay:**  
   - Users receive batches of 5 questions with destination clues and four multiple-choice options.
   - When an answer is submitted, the backend checks the answer, updates the player’s statistics (including correct/total attempts and score), and returns fun facts and trivia about that destination.
3. **Social Features:**  
   - Users can share a unique invite link (e.g., `https://globetrotter.aryansingh.dev/?invite-code=<user_id>`) to add friends.
   - Invited friends automatically become mutual friends, and users can view a leaderboard based on friend performance.
4. **Persistent UI:**  
   - The Navbar (displaying player details and navigation options) remains visible as long as the user is logged in.

---
# API Handler Logic
## 2. Create Player & Get Player Internals

### Create Player API

- **Purpose:**  
  On first login, create a player record in the `players` table using the Supabase user ID. The API uses the Dicebear API to generate an avatar.
- **Internal Workflow:**
  - Uses the Supabase auth user ID as the player `id`.
  - Calls Dicebear to generate an avatar based on the user ID.
  - Inserts a new player record with default score metrics.
    
### Get Player API

- **Purpose:**  
  Retrieves the logged-in player's record.

## 3. Invite System

![image](https://github.com/user-attachments/assets/989f2cc4-ca4e-4c64-90b0-0eaf4057390f)

### Overview

- **Invite Flow:**  
  Users can share a unique invite link with query parameters (e.g., `?invite-code=<user_id>`).  
  When a new user signs up/ logs in using the invite link, the backend automatically creates mutual friendship records between the inviter and the invitee.

## 4. Get Leaderboard

### Overview

- **Functionality:**  
  Displays a leaderboard of the player's friends and their game statistics.
  Fetches all the friends of the player and filter the response along with sorting it in descending score.


## 5. Get Questions

### Overview

- **Functionality:**  
  Fetches a batch of 5 random questions. For each question:
  - **question_id:** The destination ID.
  - **question_hints:** An array of clues.
  - **answer_options:** Four options (one correct, three incorrect) with format `"City, Country"`.

## 6. Check Answer

### Overview
The Wilson score interval method provides a statistically robust way to rank players by incorporating the uncertainty inherent in small sample sizes. Instead of simply calculating the percentage of correct answers, it computes a confidence interval for a player's true success rate, using a z-score (typically 1.96 for 95% confidence). This approach penalizes players with few attempts—preventing someone who gets one out of one correct from being unfairly ranked at 100%—while rewarding consistency as more data is accumulated. In essence, the Wilson score interval offers a fairer and more reliable measure of performance by balancing raw accuracy with the volume of attempts.


<img width="502" alt="image" src="https://github.com/user-attachments/assets/aee32eca-7b3c-4fa1-a22e-7301de61200c" />

- **Functionality:**  
  Accepts a question ID and an answer (formatted as `"City, Country"`).  
  The API:
  - Compares the submitted answer to the correct one.
  - Updates the player’s record (incrementing `total_attempts` and, if correct, `correct_answers`).
  - Recalculates the player's score using the Wilson score interval method.
  - Returns fun facts, trivia, the correct answer, and updated stats.

---

## 7. Local Setup

### Prerequisites

- **Backend:**  
  - Golang (v1.18+ recommended)
  - Supabase Project
- **Frontend:**  
  - Node.js and pnpm

or 

- `docker` and `docker-compose` installed
- configure the .env in /server and /client
- ```
  docker-compose up
  ```
  
### Steps

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/meltedhyperion/globetrotter.git
   cd globetrotter
   ```

2. **Backend Setup:**
   - Configure your environment variables in a `.env` file in /source:
     ```env
      API_PORT='<api port>'
      SUPABASE_URL='<supabase url>'
      SUPABASE_ANON_KEY='<supabase anon key>'
      SUPABASE_JWT_SECRET='<supabase jwt secret>'
     ```

   - Start the backend server:
     ```bash
     
     cd server
     go run cmd/*.go
     ```

3. **Frontend Setup:**
    - Configure your environment variables in a `.env` file in /client:
     ```env
      VITE_SUPABASE_URL='<supabase url>'
      VITE_SUPABASE_ANON_KEY='<supabase anon key>'
      VITE_BACKEND_API='<backend api>'
     ```
   - Navigate to the frontend directory:
     ```bash
     cd client
     pnpm install
     ```
   - Start the development server:
     ```bash
     pnpm dev
     ```
   - The application should be available at `http://localhost:5173` (or as specified).

---

## 8. API Documentation

All API endpoints require an `Authorization` header with a valid JWT token from Supabase.

### Create Player

**Request:**
```http
POST /players/create HTTP/1.1
Host: localhost:5050
Content-Type: application/json
Authorization: Bearer <access_token>

{
  "name": "john_doe"
}
```

**Response:**
```json
{
  "status": 201,
  "message": "Player created successfully",
}
```

### Get Player

**Request:**
```http
GET /players HTTP/1.1
Host: localhost:5050
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "status": 200,
  "message": "Player fetched successfully",
  "data": {
    "id": "9769c5e4-6c17-4ad6-9c50-64635d897847",
    "name": "john_doe",
    "avatar": "https://api.dicebear.com/7.x/pixel-art/svg?seed=9769c5e4-6c17-4ad6-9c50-64635d897847&backgroundColor=ffd5dc&size=128",
    "correct_answers": 2,
    "total_attempts": 3,
    "score": 0.47
  }
}
```

### Get Questions

**Request:**
```http
GET /questions HTTP/1.1
Host: localhost:5050
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "status": 200,
  "message": "Questions fetched successfully",
  "data": [
    {
      "question_id": 5,
      "question_hints": [
        "Known for its ancient ruins and the Colosseum.",
        "This city is a treasure trove of history and art."
      ],
      "answer_options": [
        "Rome, Italy",
        "Barcelona, Spain",
        "New York, USA",
        "Tokyo, Japan"
      ]
    },
    ...
  ]
}
```

### Check Answer

**Request:**
```http
POST /questions/check HTTP/1.1
Host: localhost:5050
Content-Type: application/json
Authorization: Bearer <access_token>

{
  "question_id": 5,
  "answer": "Rome, Italy"
}
```

**Response:**
```json
{
  "status": 200,
  "message": "Answer checked successfully",
  "data": {
    "correct": true,
    "fun_facts": [
      "Rome was founded in 753 BC and has a history of over 2,500 years.",
      "The Vatican City, an independent state, is located within Rome."
    ],
    "trivia": [
      "Rome is also called the Eternal City.",
      "The city has over 900 churches, including St. Peter's Basilica."
    ],
    "correct_answer": "Rome, Italy",
    "correct_answers": 5,
    "total_attempts": 5,
    "score": 0.565508505247919
  }
}
```

### Invite Friend

**Request:**
```http
POST /friends/<invite-code> HTTP/1.1
Host: localhost:5050
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "status": 201,
  "message": "Friend added successfully",
}
```

### Get Leaderboard

**Request:**
```http
GET /players/leaderboard HTTP/1.1
Host: localhost:5050
Authorization: Bearer <access_token>
```

**Response:**
```json
{
  "status": 200,
  "message": "Leaderboard fetched successfully",
  "data": {
    "player_stats": [
     {
          "id": "9769c5e4-6c17-4ad6-9c50-64635d897847",
          "name": "john_doe",
          "avatar": "https://api.dicebear.com/7.x/pixel-art/svg?seed=9769c5e4-6c17-4ad6-9c50-64635d897847&backgroundColor=ffd5dc&size=128",
          "score": 0.85,
          "correct_answers": 10,
          "total_attempts": 12
        },
        {
          "id": "another-user-id",
          "name": "jane_doe",
          "avatar": "https://api.dicebear.com/7.x/pixel-art/svg?seed=9769c5e4-6c17-4ad6-9c50-64635d897847&backgroundColor=ffd5dc&size=128",
          "score": 0.65,
          "correct_answers": 8,
          "total_attempts": 15
        }
    ]
  }
}
```

---
## 9. How I utilized AI

I used various AI tools like ChatGPT, Windsurf, V0.dev and Bolt.new. Each for different purposes.

1) Chatgpt to generate dataset and asking suggestions on what apprach to go for in cases. Ex: Batching of 5 questions in 1 request to limit the number of API calls. Use of Wilson score interval method for calculation of scores for effective rankings. It also helped in writing down this README ;p

2) Windsurf as my primary code editor, for auto code completions for productivity.

3) I used Bolt.new at the very start, prompting it the challenge statement to visualize what kind of application and functionality is expected. As I am not applying for Frontend (also not very fluent in it), I used its code as a starter to include more functionalities like leaderboard, auth components.

4) I used V0.dev for creating several components. Auth, Game page and leaderboard.
---

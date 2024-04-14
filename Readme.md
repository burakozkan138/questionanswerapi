**Question Answer API**

## Overview

This project aims to develop a question-answer API. The API provides users with the ability to ask questions and receive answers. Additionally, API users can perform operations such as saving, updating, and deleting questions.

## Installation
You can run the project locally by following the steps below.

### Step 1: Clone the Repository
Clone the GitHub repository to your local environment:
```bash
git clone https://github.com/burakozkan138/questionanswerapi.git
```

### Step 2: Navigate to the Project Directory
Navigate to the root directory of the project.
```bash
cd questionanswerapi
```

### Step 3: Configure Environment
Copy the required config files from the example files and make the necessary modifications.

```bash
# For Linux
cp config/.env.example config/.env && cp config/.env.example config/.env.test

# For Windows
copy config\.env.example config\.env && copy config\.env.example config\.env.test
```

### Step 4: Deploy with Docker
After configuring the environment, ensure that the database environment section in the docker-compose.yml file is correct, then deploy the project on Docker.
```bash
docker-compose up -d --build
```

## Usage
If the project has been successfully deployed, visit the following URL to access the Swagger interface:

```bash
http://localhost:8080/swagger
```

## Notes
To view the Postman workspace used in the project, you can utilize the following link:
### https://www.postman.com/cinemabookingsystem/workspace/questionanswer
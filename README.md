[Русская версия]()

# IHATEINSTAGRAM
**IHATEINSTAGRAM** is a Golang backend project that showcases my knowledge and skills in web development. It uses Fiber, a fast and lightweight web framework, to handle HTTP requests and responses. It uses pgx, a PostgreSQL driver and toolkit, to connect and interact with the database. It uses sessions for secure authentication and authorization, to protect the endpoints and verify the users. **IHATEINSTAGRAM** is designed with clean architecture principles. It can be easily deployed with Docker, a platform for containerization and orchestration.
 The project provides functionality for profiles, posting content, uploading pictures in posts, liking posts, and following other users. **IHATEINSTAGRAM** is a project that I'm proud of and I hope it will impress potential employers and clients.

## How to run
Before you continue, ensure you have met the following requirements:
- You have installed the latest version of Go.
- You have installed PostgreSQL and created a database for the project.
- You have created .env file, with `POSTGRES_URL` defined.

To run **IHATEINSTAGRAM**, follow these steps:
1. Clone this repository: `git clone https://github.com/indetensai/ihateinstagram.git`
2. Change into the project directory: `cd ihateinstagram`
3. Install the dependencies: `go mod download`
4. Fill `.env` file with the required environment variables.
5. Build the executable: `go build -o ihateinstagram cmd/ihateinstagram/main.go`

## How to run (docker-compose)
Before you continue, ensure you have met the following requirements:
- You have installed the latest version of docker(-desktop).

To run **IHATEINSTAGRAM** using docker-compose, follow these steps:
1. Clone this repository: `git clone https://github.com/indetensai/ihateinstagram.git`
2. Change into the project directory: `cd ihateinstagram`
3. Run `docker compose up`

## Usage
To run **IHATEINSTAGRAM**, follow these steps:
1. Start the executable: `./ihateinstagram`
2. The server will listen on port 8080.
3. To interact with the chat API, you can use any HTTP client of your choice.

The **IHATEINSTAGRAM** API has the following endpoints:
- `POST /user/register`: Create a new user account.
- `POST /user/login`: Login with an existing user account and get a session.
- `DELETE /session/:id`: Delete session.
- `PUT /user/:user_id/followers/:follower_id`: Follow follower_id user on user_id user. Requires authentication.
- `DELETE /user/:user_id/followers/:follower_id`: Unfollow follower_id user on user_id user. Requires authentication.
- `GET /user/:user_id/followers`:Get followers of user_id user. 
- `POST /post"`: Create post. Requires authentication.
- `GET /post/:post_id` Get post by post_id. Requires authentication.
- `PATCH /post/:post_id`: Change post_id post. Requires authentication.
- `PUT /post/:post_id/like`: Like post_id post. Requires authentication.
- `GET /post/:post_id/likes`: Get likes on post_id post.
- `DELETE /post/:post_id/like`: Delete like on post_id post. Requires authentication.
- `DELETE /post/:post_id`: Delete post_id post. Requires authentication.
- `POST /post/:post_id/image`: Upload image in post_id post. Requires authentication.
- `GET /post/:post_id/images`: Get post_id post images.
- `GET /post/:post_id/thumbnails`: Get post_id post thumbnails
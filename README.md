# Web Forum

This is a web forum application that allows communication between users, associating categories to posts, liking and disliking posts and comments, and filtering posts. The project is built using Golang and SQLite.

## Objectives

The objectives of this project are:

- Allow users to communicate between each other through posts and comments.
- Associate categories to posts.
- Allow users to like and dislike posts and comments.
- Implement a filter mechanism to filter posts by categories, created posts, and liked posts.
- Use SQLite to store the data for the application.
- Implement user authentication and sessions using cookies.
- Follow good coding practices and handle all sorts of errors.

## Installation

To run this project, you need to have Docker installed on your machine. Then follow these steps:

1. Clone this repository.
2. Navigate to the project directory.
3. Run the following command to build the Docker image:

```bash
make build
```
4. Run the following command to start the Docker container:

```bash
make run
```
5. Open your web browser and go to http://localhost:8080.

If you want to run without Docker:
```bash
go run /cmd/web/*
```

## Usage
To use the web forum application, follow these steps:
1. Register a new user by providing your email, username, and password.
2. Log in to the application using your email and password.
3. Create a post by providing a title, content, and one or more categories.
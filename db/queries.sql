-- name: ListAllUsers :many
SELECT id, name, email, created_at FROM users;

-- name: CreateList :one
INSERT INTO
    list (title, user_id, status)
VALUES ($1, $2, 'pending') RETURNING *;

-- name: CreateUser :one
INSERT INTO
    users (
        name,
        email,
        hashed_password,
        created_at
    )
VALUES ($1, $2, $3, NOW()) RETURNING id, name, email, created_at;

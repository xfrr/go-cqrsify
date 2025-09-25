#!/bin/bash

docker run --rm -d --name cqrsify-uow-postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_DB=cqrsify_uow_example \
    -p 5432:5432 postgres:latest
echo "Waiting for PostgreSQL to start..."
sleep 5
docker exec -it cqrsify-uow-postgres psql -U postgres -d cqrsify_uow_example -c "CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT NOT NULL);"
docker exec -it cqrsify-uow-postgres psql -U postgres -d cqrsify_uow_example -c "CREATE TABLE orders (id SERIAL PRIMARY KEY, user_id INT NOT NULL, amount DECIMAL NOT NULL, FOREIGN KEY (user_id) REFERENCES users(id));"
echo "PostgreSQL container with database and tables created."
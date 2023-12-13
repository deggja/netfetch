# Stage 1: Build the Vue.js frontend
FROM node:latest as nodebuilder

WORKDIR /app/frontend

# Copy the frontend code - make sure this path matches the location of your Vue.js files in the repo
COPY frontend/dash /app/frontend

# Install dependencies and build the frontend
RUN npm install
RUN npm run build

# Stage 2: Build the Go backend and bundle the frontend
FROM golang:1.21 as gobuilder

WORKDIR /app/backend

# Copy the backend code
COPY backend /app/backend

# Install statik
RUN go install github.com/rakyll/statik@latest

# Bundle the frontend with the backend
COPY --from=nodebuilder /app/frontend/dist /app/frontend/dist
RUN statik -src=/app/frontend/dist
RUN go build -o netfetch

# Stage 3: Final image
FROM debian:latest

WORKDIR /usr/local/bin/

# Copy the compiled Go binary
COPY --from=gobuilder /app/backend/netfetch /usr/local/bin/

# Install any runtime dependencies here (if necessary)

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
CMD ["netfetch", "dash"]

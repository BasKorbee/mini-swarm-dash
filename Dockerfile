# Use golang image to build server
FROM golang:1.24 AS server-build-stage

WORKDIR /app

# Copy everything from server to CWD
COPY ./server .

# Download libraries
RUN go mod download

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /cmd

# Use node image to build client
FROM node:24-slim AS client-build-stage
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

# Enable package managers
RUN corepack enable

WORKDIR /app

# Copy everything from client to CWD
COPY ./client .

# Install prod dependencies
FROM client-build-stage AS client-prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile --force

# Install build dependencies
FROM client-build-stage AS client-build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile --force

# Build client
RUN pnpm run build

# Deploy the application binary into a lean image
FROM alpine:latest

WORKDIR /

# Copy server binary
COPY --from=server-build-stage /cmd /cmd
# Copy client dependencies 
COPY --from=client-prod-deps /app/node_modules /public/node_modules
# Copy client code
COPY --from=client-build /app/dist /public

EXPOSE 8080

ENTRYPOINT ["/cmd"]
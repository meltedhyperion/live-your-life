# Step 1: Build Stage
FROM node:18 AS build

# Set working directory to /app inside the container
WORKDIR /app

# Install pnpm globally
RUN npm install -g pnpm

# Copy the package.json and pnpm-lock.yaml (if present) for dependency installation
COPY package.json pnpm-lock.yaml* ./

# Install dependencies using pnpm
RUN pnpm install

# Copy the entire project into the container
COPY . .

# Build the project (this outputs the build in /client/dist)
RUN pnpm build

# Step 2: Production Stage
FROM nginx:alpine

# Copy the built app from the build stage to the nginx html folder
COPY --from=build /app/dist /usr/share/nginx/html

# Expose port 80 (default for nginx)
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]

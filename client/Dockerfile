# Development container for React Native
FROM node:18-alpine

WORKDIR /app

# Install React Native CLI and development tools
RUN npm install -g @react-native-community/cli expo-cli

# Install system dependencies
RUN apk add --no-cache \
    git \
    python3 \
    make \
    g++ \
    watchman

# Copy package files
COPY package*.json ./
COPY yarn.lock ./

# Install dependencies
RUN yarn install

# Copy source code
COPY . .

# Expose Metro bundler port
EXPOSE 8081

# Start Metro bundler
CMD ["yarn", "start"]

FROM node:22-alpine AS builder
 
WORKDIR /app
 
COPY package.json package-lock.json ./
RUN npm ci
 
COPY tsconfig.json ./
COPY src ./src
 
RUN npm run build
 
 
FROM node:22-alpine AS production-deps
 
WORKDIR /app
 
COPY package.json package-lock.json ./
RUN npm ci --omit=dev
 
 
FROM node:22-alpine
 
WORKDIR /app
 
ENV NODE_ENV=production
 
COPY --from=production-deps /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
COPY package.json ./
 
CMD ["node", "dist/index.js"]

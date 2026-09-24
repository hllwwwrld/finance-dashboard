FROM node:24-alpine AS builder

WORKDIR /app

ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./

# npm 12 requires Node ^22.22.2 || ^24.15.0 || >=26; pin pnpm for allowBuilds support
RUN npm install -g npm@12.1.0 pnpm@10.28.0 && pnpm install --frozen-lockfile

COPY . .

RUN pnpm build

FROM node:24-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV NEXT_TELEMETRY_DISABLED=1

RUN npm install -g npm@12.1.0

COPY --from=builder /app ./

EXPOSE 3000

CMD ["npm", "run", "start"]

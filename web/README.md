# Web

Next.js frontend for `aws-log-practice`.

## Development

Install dependencies:

```bash
corepack pnpm install
```

Start the development server:

```bash
corepack pnpm run dev
```

Open `http://localhost:3000`.

## Scripts

```bash
corepack pnpm run dev        # Start Next.js dev server
corepack pnpm run build      # Build for production
corepack pnpm run start      # Start production server after build
corepack pnpm run lint       # Run Biome checks
corepack pnpm run format     # Format files with Biome
corepack pnpm run proto:gen  # Generate frontend protobuf types
```

## Backend

The frontend is currently a UI scaffold. API calls and React hooks are expected
to be added separately.

When wiring ConnectRPC clients, use generated types under `gen/product/v1`.

---
description: Generate complete project scaffolding for a new Next.js + Firebase app
---

Create complete project scaffolding for: [APP_NAME]

Reference implementation patterns from: `/Users/sloan/code/bildit/gitbot-tester`
Reference deployment system: `/Users/sloan/code/devops-cloud-run`

## Prerequisites
- PRD and Technical Specification have been created
- App name is decided: [WAIT FOR USER INPUT]

## Project Structure to Generate

```
~/code/{app-name}/
├── nextapp/
│   ├── src/
│   │   ├── app/
│   │   ├── components/
│   │   ├── lib/
│   │   └── types/
│   ├── public/
│   ├── Dockerfile
│   ├── .dockerignore
│   ├── package.json
│   ├── tsconfig.json
│   ├── next.config.ts
│   ├── tailwind.config.ts
│   ├── postcss.config.mjs
│   └── .env.local.example
├── scripts/
│   ├── setup-gcp.sh
│   └── create-firestore-indexes.sh
├── docs/
│   ├── prd.md
│   └── technical-spec.md
├── bitbucket-pipelines.yml
├── .gitignore
└── README.md
```

## Files to Generate

### 1. nextapp/package.json
```json
{
  "name": "{app-name}",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev --turbopack",
    "build": "next build",
    "start": "next start",
    "lint": "eslint",
    "type-check": "tsc --noEmit"
  },
  "dependencies": {
    "@google-cloud/secret-manager": "^6.1.0",
    "@tailwindcss/postcss": "^4",
    "firebase": "^10.14.1",
    "firebase-admin": "^12.7.0",
    "jose": "^5.10.0",
    "next": "15.5.4",
    "react": "19.1.0",
    "react-dom": "19.1.0",
    "tailwindcss": "^4",
    "zod": "^4.1.11"
  },
  "devDependencies": {
    "@eslint/eslintrc": "^3",
    "@types/node": "^20",
    "@types/react": "^19",
    "@types/react-dom": "^19",
    "eslint": "^9",
    "eslint-config-next": "15.5.4",
    "typescript": "^5"
  }
}
```

### 2. nextapp/tsconfig.json
```json
{
  "compilerOptions": {
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": ["./src/*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

### 3. nextapp/next.config.ts
```typescript
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  experimental: {
    serverActions: {
      bodySizeLimit: "2mb",
    },
  },
};

export default nextConfig;
```

### 4. nextapp/tailwind.config.ts
```typescript
import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};

export default config;
```

### 5. nextapp/postcss.config.mjs
```javascript
/** @type {import('postcss-load-config').Config} */
const config = {
  plugins: {
    "@tailwindcss/postcss": {},
  },
};

export default config;
```

### 6. nextapp/Dockerfile
Generate a multi-stage Dockerfile based on the gitbot-tester pattern with:
- Node 18 Alpine base
- Build args for Firebase config
- Multi-stage build (deps, builder, runner)
- Non-root user
- Port 3000
- BuildKit cache optimization

### 7. nextapp/.dockerignore
```
node_modules
.next
.env*.local
.git
.gitignore
README.md
npm-debug.log
.DS_Store
```

### 8. nextapp/src/app/layout.tsx
```typescript
import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "{App Name}",
  description: "{App description from PRD}",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
```

### 9. nextapp/src/app/page.tsx
```typescript
import { redirect } from "next/navigation";
import { verifySession } from "@/lib/auth/session";

export default async function Home() {
  const session = await verifySession();

  if (!session) {
    redirect("/login");
  }

  return (
    <main>
      <h1>{App Name} Dashboard</h1>
      {/* Main app content */}
    </main>
  );
}
```

### 10. nextapp/src/app/login/page.tsx
Generate login page with:
- Google Sign In button
- Firebase auth integration
- Session creation on success
- Redirect to dashboard

### 11. nextapp/src/middleware.ts
Generate auth middleware that:
- Checks session cookie on protected routes
- Verifies JWT signature
- Redirects to /login if unauthenticated
- Allows public routes (/login, /api/auth/*, /api/health)

### 12. nextapp/src/lib/auth/firebase-config.ts
```typescript
import { initializeApp, getApps } from "firebase/app";
import { getAuth } from "firebase/auth";

const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY!,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN!,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID!,
};

const app = getApps().length === 0 ? initializeApp(firebaseConfig) : getApps()[0];
export const auth = getAuth(app);
```

### 13. nextapp/src/lib/auth/admin.ts
Generate Firebase Admin initialization for server-side auth

### 14. nextapp/src/lib/auth/session.ts
Generate JWT session management with:
- createSession(uid, email)
- verifySession()
- deleteSession()
- Uses JWT_SECRET from Secret Manager

### 15. nextapp/src/lib/firestore/firestore-client.ts
```typescript
import { Firestore } from "@google-cloud/firestore";

let firestoreInstance: Firestore | null = null;

export function getFirestoreClient(): Firestore {
  if (!firestoreInstance) {
    firestoreInstance = new Firestore({
      projectId: process.env.GOOGLE_CLOUD_PROJECT!,
    });
  }
  return firestoreInstance;
}
```

### 16. nextapp/src/app/api/health/route.ts
```typescript
export async function GET() {
  return Response.json({
    status: "healthy",
    timestamp: new Date().toISOString(),
    service: "{app-name}",
  });
}
```

### 17. nextapp/src/app/api/auth/session/route.ts
Generate session creation endpoint

### 18. nextapp/src/app/api/auth/logout/route.ts
Generate logout endpoint

### 19. nextapp/src/types/index.ts
Generate TypeScript types based on PRD data models

### 20. nextapp/.env.local.example
```bash
# Firebase Configuration (Public)
NEXT_PUBLIC_FIREBASE_API_KEY=your-api-key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id

# GCP Configuration
GOOGLE_CLOUD_PROJECT=your-project-id

# Local Development Secrets (use ADC for cloud)
JWT_SECRET=your-local-jwt-secret

# Local Development Setup:
# 1. gcloud auth application-default login
# 2. gcloud config set project YOUR_PROJECT_ID
```

### 21. scripts/setup-gcp.sh
Generate GCP setup script that:
- Creates service account: {app-name}-app@{PROJECT}.iam.gserviceaccount.com
- Enables APIs: Identity Platform, Firestore, Secret Manager
- Creates Firestore database (Native mode, nam5)
- Configures minimal IAM permissions
- Creates jwt-secret in Secret Manager
- Grants deployment SA access to secrets
- Includes manual Identity Platform setup instructions

Reference: `/Users/sloan/code/bildit/gitbot-tester/scripts/setup-gcp.sh`

### 22. scripts/create-firestore-indexes.sh
Generate script to create composite indexes based on PRD queries

### 23. bitbucket-pipelines.yml
Generate CI/CD configuration that:
- Parallel build/test and deployment steps
- Uses devops-cloud-run repository
- Configures service account: {app-name}-app@${GCP_PROJECT}.iam.gserviceaccount.com
- Sets Firebase build args
- Configures environment variables and secrets
- Includes health checks and testing
- Has branch cleanup pipeline

Reference: `/Users/sloan/code/bildit/gitbot-tester/bitbucket-pipelines.yml`

### 24. .gitignore
```
# Dependencies
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Next.js
.next/
out/
build/
dist/

# Environment
.env*.local
.env.production

# Testing
coverage/

# Misc
.DS_Store
*.pem

# Debug
*.log

# Local
.cache/
```

### 25. README.md
Generate comprehensive README with:
- Project overview
- Tech stack
- Prerequisites
- Setup instructions (GCP, local development)
- Deployment guide
- Environment variables reference
- Architecture documentation

---

## After Generation

1. Create the directory structure: `mkdir -p ~/code/{app-name}`
2. Generate all files listed above
3. Initialize git: `cd ~/code/{app-name} && git init`
4. Create initial commit: `git add . && git commit -m "feat: initial project scaffolding"`
5. Remind user of next steps:
   - Run `scripts/setup-gcp.sh` to create GCP infrastructure
   - Set up Identity Platform in GCP Console
   - Create `.env.local` for local development
   - Run `npm install` in nextapp directory
   - Test with `npm run dev`

Save all files with proper formatting and no placeholders.

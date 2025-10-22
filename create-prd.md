---
description: Create a comprehensive PRD and Technical Specification for an app idea and save to file
---

Create a complete Product Requirements Document (PRD) and Technical Specification for the following app idea:

[USER WILL PROVIDE APP IDEA - WAIT FOR INPUT]

The user may either:
1. Provide an app name directly (e.g., "habit-tracker-pro")
2. Reference an existing idea file (e.g., "use the idea from docs/app-ideas/habit-tracker-pro.md")
3. Describe a new app idea

## Document Structure

### 1. EXECUTIVE SUMMARY
- **App Name**:
- **Tagline**: One-sentence description
- **Target Audience**: Specific demographics and characteristics
- **Core Value Proposition**: What problem does this solve and why would people use it?
- **Success Metrics**: How will we measure if this app is successful?

### 2. USER PERSONAS (2-3)
For each persona include:
- **Name & Role**: (e.g., "Sarah, Freelance Designer")
- **Demographics**: Age, location, profession
- **Goals**: What are they trying to achieve?
- **Pain Points**: What frustrates them currently?
- **Technical Proficiency**: Beginner / Intermediate / Advanced
- **How They'll Use This App**: Specific scenarios

### 3. MVP FEATURE SPECIFICATIONS
List all features for the Minimum Viable Product:

**Must-Have Features (Priority 1)**:
- Feature 1: [Name and detailed description]
  - User Story: As a [persona], I want to [action], so that [benefit]
  - Acceptance Criteria: [List specific requirements]
  - Technical Notes: [Implementation considerations]

[Continue for all must-have features]

**Nice-to-Have Features (Priority 2)**:
[Features for v2 or if time permits]

### 4. TECHNICAL ARCHITECTURE

**Frontend**:
- Framework: Next.js 15+ (App Router)
- Language: TypeScript
- Styling: Tailwind CSS 4
- UI Patterns: Server Components (default), Client Components (where needed)

**Backend Services**:
- Database: Firebase Firestore (Native mode)
- Authentication: Firebase Auth with Google OAuth
- Secrets: GCP Secret Manager
- Storage: (if needed) GCP Cloud Storage or Firebase Storage

**Infrastructure**:
- Container: Docker (multi-stage build)
- Hosting: GCP Cloud Run
- CI/CD: Bitbucket Pipelines
- Deployment: devops-cloud-run repository

### 5. DATA MODELS

For each Firestore collection:

**Collection: `collection-name`**
```typescript
interface CollectionName {
  id: string;
  field1: string;
  field2: number;
  field3: Timestamp;
  // ... all fields with types
}
```
- **Purpose**: What this collection stores
- **Access Patterns**: How data is queried
- **Indexes Needed**: List composite indexes
- **Security Rules**: Who can read/write

[Repeat for all collections]

**Relationships**:
- Diagram showing how collections relate
- Denormalization strategy (if applicable)

### 6. SCREEN INVENTORY

List every page/route in the application:

**Public Routes**:
- `/` - Landing page (if app has one)
- `/login` - Authentication page

**Protected Routes**:
- `/dashboard` - Main application view
- `/[feature]` - Feature-specific pages
- `/settings` - User settings
- `/[other routes]`

**API Routes**:
- `/api/health` - Health check endpoint
- `/api/auth/session` - Session creation
- `/api/auth/logout` - Session termination
- `/api/[feature endpoints]` - Feature-specific APIs

For each screen, describe:
- Purpose
- Key components
- Data displayed
- User actions available
- Navigation to/from this screen

### 7. USER FLOWS

**Authentication Flow**:
```
User lands on app →
Check session cookie →
If authenticated: Redirect to dashboard →
If not: Show landing/login →
Click "Sign in with Google" →
Firebase auth flow →
Create session JWT →
Redirect to dashboard
```

**Primary Feature Flows** (for each major feature):
```
Step 1 → Step 2 → Step 3 → Success state
         ↓ (error handling)
         Error state → Recovery action
```

**Edge Cases**:
- No data state
- Loading states
- Error states
- Offline behavior (if applicable)

### 8. SECURITY & PERMISSIONS

**Service Account**:
- Name: `{app-name}-app@{PROJECT_ID}.iam.gserviceaccount.com`
- Required Permissions:
  - `roles/datastore.user` (Firestore access)
  - `roles/identitytoolkit.viewer` (Auth token verification)
  - `roles/secretmanager.secretAccessor` (per-secret bindings)

**Secrets in Secret Manager**:
- `jwt-secret` - Session signing key
- [Other secrets needed]

**Firestore Security Rules** (outline):
```javascript
match /collection-name/{docId} {
  allow read: if [conditions];
  allow write: if [conditions];
}
```

**Route Protection**:
- List which routes require authentication
- Public API endpoints (if any)
- Admin-only features (if any)

### 9. DEPLOYMENT CONFIGURATION

**Environment Variables (Build-Time)**:
```bash
# These are baked into the client bundle
NEXT_PUBLIC_FIREBASE_API_KEY=...
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=...
NEXT_PUBLIC_FIREBASE_PROJECT_ID=...
```

**Environment Variables (Runtime)**:
```bash
# Plain environment variables
GOOGLE_CLOUD_PROJECT=...
[Other non-secret vars]
```

**Secret Manager References**:
```bash
# Mounted by Cloud Run
JWT_SECRET=<from-secret-manager>
[Other secrets]
```

**Cloud Run Configuration**:
- CPU: 1 vCPU
- Memory: 512Mi (or higher if needed)
- Min Instances: 0 (scale to zero)
- Max Instances: 100
- Concurrency: 80
- Service Account: {app-name}-app

**Docker Build Args**:
- Firebase config (public values)
- Any build-time configuration

### 10. DEVELOPMENT PHASES

Break the implementation into logical phases:

**Phase 1: Foundation**
- [ ] Project scaffolding
- [ ] Authentication system
- [ ] Database client setup
- [ ] Basic routing and layout

**Phase 2: Core Features**
- [ ] Feature 1 implementation
- [ ] Feature 2 implementation
- [ ] Feature 3 implementation

**Phase 3: Polish**
- [ ] Loading states
- [ ] Error handling
- [ ] Empty states
- [ ] Responsive design
- [ ] Animations

**Phase 4: Deployment**
- [ ] GCP infrastructure setup
- [ ] CI/CD configuration
- [ ] Production deployment
- [ ] Monitoring & health checks

### 11. SUCCESS CRITERIA

**Functional Requirements**:
- All MVP features work as specified
- Authentication flow is secure
- Data persists correctly
- API endpoints respond properly

**Performance Requirements**:
- Cold start < 5 seconds
- Time to interactive < 3 seconds
- API response time < 500ms (p95)
- No memory leaks

**Quality Requirements**:
- No TypeScript errors
- No linting errors
- 100% of primary user flows tested
- Works on mobile, tablet, desktop
- Accessible (WCAG AA basics)

### 12. OUT OF SCOPE (for MVP)

List features explicitly NOT included in this version:
- Feature X (save for v2)
- Integration Y (not needed yet)
- Admin panel (build later if needed)

---

Make this document comprehensive enough that a developer could build the entire application from it without needing to make architectural decisions.

## Save Output

Save this PRD to: `/Users/sloan/code/app-ideas/docs/prds/{app-name-kebab-case}-prd.md`

For example:
- `habit-tracker-pro-prd.md`
- `recipe-cost-calculator-prd.md`
- `meeting-cost-analyzer-prd.md`

After saving, confirm the file location and suggest next steps:
1. Review and refine the PRD
2. Create implementation plan (before scaffolding)
3. Use `/scaffold-app` to generate project structure

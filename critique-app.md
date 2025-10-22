---
description: Perform comprehensive technical review and identify all issues in the application
---

You are a senior technical architect and security expert reviewing this Next.js application. Be extremely critical and identify every issue, no matter how small.

## Application to Review

App Name: [WAIT FOR USER INPUT]
Location: [USER WILL PROVIDE PATH]

## Review Categories

### 1. CODE QUALITY

Examine all code files and identify:
- TypeScript usage and type safety issues
- Missing or incorrect type definitions
- Component patterns and organization issues
- Server vs Client Component decisions that could be optimized
- Error handling completeness and quality
- Code duplication
- Unused imports or variables
- Console.logs in production code
- Missing null/undefined checks
- Async/await patterns and error handling
- Variable naming and code clarity

### 2. ARCHITECTURE & PATTERNS

Review the overall architecture:
- Data flow and state management approach
- API design and Server Actions implementation
- Database query efficiency and N+1 problems
- Security implementation (authentication, authorization)
- Separation of concerns (is business logic properly separated?)
- Component reusability and composition
- File organization and structure
- Import paths and module organization
- Use of Next.js features (App Router, Server Components, etc.)

### 3. SECURITY

Critical security review:
- Authentication implementation (session management, cookies, JWT)
- Authorization checks on all protected routes
- Input validation (forms, API endpoints, server actions)
- SQL injection / NoSQL injection risks
- XSS vulnerabilities
- CSRF protection
- Secret management (are secrets in code? environment variables?)
- Firestore security rules (do they match the application logic?)
- Service account permissions (principle of least privilege?)
- API endpoint protection
- Rate limiting considerations
- Session expiry and refresh logic

### 4. USER EXPERIENCE

Evaluate UX comprehensively:
- Navigation clarity and intuitiveness
- Form validation and error feedback quality
- Loading states (are they present on all async operations?)
- Error states (are all error cases handled gracefully?)
- Empty states (what happens when lists are empty?)
- Success feedback (do users know when actions succeed?)
- Mobile responsiveness (test at 375px, 768px, 1024px)
- Touch targets (minimum 44x44px on mobile?)
- Keyboard navigation
- Accessibility basics (alt text, ARIA labels, semantic HTML)
- Consistency in spacing, colors, typography
- Animation and transition quality
- Text readability and hierarchy

### 5. PERFORMANCE

Analyze performance:
- Bundle size and code splitting
- Image optimization (are next/image used correctly?)
- Database query patterns and indexes
- Caching strategy (is data fetched efficiently?)
- Server vs Client Component optimization
- Lazy loading for heavy components
- Cold start impact (unnecessary dependencies?)
- Memory leaks (useEffect cleanup, listeners, intervals)
- Infinite loops or excessive re-renders
- Large dependency bundle sizes
- Build time and optimization

### 6. DATA & STATE

Review data handling:
- Firestore collection structure (is it efficient?)
- Data denormalization strategy (appropriate?)
- Query patterns (are composite indexes needed?)
- Data validation (Zod schemas or similar?)
- State management approach (props drilling? context overuse?)
- Data fetching patterns (Server Components vs client?)
- Optimistic updates (where appropriate?)
- Error recovery and retry logic
- Data consistency across components

### 7. DEPLOYMENT & OPERATIONS

Check deployment configuration:
- Dockerfile optimization (layer caching, size)
- Environment variable handling (build vs runtime)
- Secret management (Secret Manager integration)
- Health check endpoint implementation
- Error logging and monitoring setup
- CI/CD pipeline correctness
- Branch deployment strategy
- Service account configuration
- Cloud Run resource allocation
- Docker build args (are Firebase configs correct?)

### 8. EDGE CASES & ERROR HANDLING

Test edge cases:
- What happens with no data?
- What happens with network errors?
- What happens with invalid input?
- What happens when auth fails?
- What happens with expired sessions?
- What happens with Firestore write failures?
- What happens with Secret Manager access failures?
- What happens when external APIs fail?
- What happens with rate limiting?
- What happens with malformed data?

### 9. TESTING & QUALITY

Assess testing coverage:
- Can the app build successfully? (`npm run build`)
- Are there TypeScript errors? (`npm run type-check`)
- Are there linting errors? (`npm run lint`)
- Are critical user flows testable?
- Is the health check endpoint working?
- Can authentication be tested?
- Are API endpoints testable?

### 10. DOCUMENTATION

Review documentation quality:
- Is the README comprehensive?
- Are setup instructions clear?
- Are environment variables documented?
- Is the architecture explained?
- Are deployment steps documented?
- Are there inline code comments where needed?
- Is the PRD still accurate?

## Output Format

For EACH issue found, provide:

### [SEVERITY] Category: Issue Title

**File**: `path/to/file.ts:123`

**Problem**: Detailed description of what's wrong

**Impact**: Why this matters (security risk, poor UX, performance issue, etc.)

**Solution**: Exact steps to fix this issue

**Code Example** (if applicable):
```typescript
// Before (wrong)
const bad = ...

// After (correct)
const good = ...
```

---

## Severity Levels

- **[CRITICAL]**: Security vulnerability, data loss risk, app crashes
- **[HIGH]**: Major UX issue, performance problem, architectural flaw
- **[MEDIUM]**: Code quality issue, minor UX problem, missing optimization
- **[LOW]**: Style inconsistency, minor refactoring opportunity

## After Review

1. **Count Issues by Severity**:
   - Critical: X
   - High: Y
   - Medium: Z
   - Low: W
   - Total: N issues found

2. **Priority Fix Order**:
   List the top 10 issues that MUST be fixed before production

3. **Estimated Fix Time**:
   - Critical/High fixes: X hours
   - Medium fixes: Y hours
   - Low fixes: Z hours

4. **Overall Assessment**:
   - Is this app production-ready? Yes/No
   - What's the biggest architectural concern?
   - What's the biggest security concern?
   - What's the biggest UX concern?

Be brutally honest. If something isn't production-ready, say so explicitly.

The goal is to ship a high-quality, secure, performant application. Don't hold back on criticism.

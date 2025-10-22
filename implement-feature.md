---
description: Implement a specific feature according to PRD specifications
---

Implement the following feature for: [APP_NAME]

**Feature Name**: [WAIT FOR USER INPUT]

**Reference PRD**: `~/code/{app-name}/docs/prd.md`

## Feature Implementation Guidelines

### Step 1: Review Feature Specification

From the PRD, identify:
- Feature description and purpose
- User stories and acceptance criteria (map to Gherkin scenarios)
- Data models involved
- Business rules to enforce
- UI/UX requirements
- Technical considerations

### Step 2: Plan Implementation Using Team Pattern

Break down the feature into Team Pattern components:

1. **State Models Layer**:
   - [ ] Immutable domain entities (frozen TypeScript interfaces)
   - [ ] State history tracking with `previous_state`
   - [ ] Basic structural validation (types, formats)
   - [ ] Zod schemas for validation
   - [ ] NO methods or behavior—pure data

2. **Business Logic Layer**:
   - [ ] **Investigators**: Pure boolean functions for business rules (e.g., `is_valid_amount()`)
   - [ ] **Enforcers**: Typed exceptions that use Investigators to enforce rules
   - [ ] **Workers**: Single-responsibility state transformations (input state → output state)
   - [ ] **Delegators**: Workflow orchestration that coordinates Workers

3. **Infrastructure Layer**:
   - [ ] Firestore service functions (CRUD operations, used by Workers)
   - [ ] Server Actions that call Delegators and handle exceptions
   - [ ] API Routes (if applicable)
   - [ ] Authentication/authorization checks

4. **UI Layer**:
   - [ ] Page routes (Server Components for data fetching)
   - [ ] Client Components (for interactivity)
   - [ ] Forms that call Server Actions
   - [ ] Loading states
   - [ ] Error states (display typed exceptions)
   - [ ] Empty states

5. **Integration**:
   - [ ] Navigation links
   - [ ] Component composition
   - [ ] Data flow: UI → Server Action → Delegator → Workers → Firestore
   - [ ] Exception handling in Server Actions

### Step 3: Implementation Pattern Using Team Pattern

Follow Team Pattern architecture with Next.js 15 App Router:

#### State Models (Immutable Domain Entities)

**File**: `src/types/states/{feature}-state.ts`

```typescript
import { z } from "zod";

// Base state interface for state history tracking
export interface BaseState {
  previous_state?: BaseState;
}

// Immutable state model (frozen)
export interface FeatureNameState extends BaseState {
  readonly id: string;
  readonly field1: string;
  readonly field2: number;
  readonly createdAt: Date;
  readonly updatedAt: Date;
  readonly userId: string;
  readonly previous_state?: FeatureNameState;
}

// Zod schema for structural validation
export const featureNameSchema = z.object({
  id: z.string().uuid(),
  field1: z.string().min(1, "Field 1 is required"),
  field2: z.number().min(0, "Field 2 must be positive"),
  createdAt: z.date(),
  updatedAt: z.date(),
  userId: z.string(),
});

// Helper to create immutable state
export function createFeatureNameState(
  data: Omit<FeatureNameState, "previous_state">
): Readonly<FeatureNameState> {
  return Object.freeze({ ...data });
}
```

#### Investigators (Business Rule Checks)

**File**: `src/lib/investigators/{feature}-investigator.ts`

```typescript
import { FeatureNameState } from "@/types/states/{feature}-state";

export class FeatureInvestigator {
  /**
   * Checks if field2 is within valid business range
   */
  static is_valid_field2(state: FeatureNameState): boolean {
    return state.field2 > 0 && state.field2 < 10000;
  }

  /**
   * Checks if user can perform this operation
   */
  static can_user_modify(state: FeatureNameState, userId: string): boolean {
    return state.userId === userId;
  }
}
```

#### Enforcers (Rule Enforcement)

**File**: `src/lib/enforcers/{feature}-enforcer.ts`

```typescript
import { FeatureNameState } from "@/types/states/{feature}-state";
import { FeatureInvestigator } from "@/lib/investigators/{feature}-investigator";

// Custom typed exception
export class FeatureError extends Error {
  constructor(
    message: string,
    public readonly error_state: FeatureNameState
  ) {
    super(message);
    this.name = "FeatureError";
  }
}

export class FeatureEnforcer {
  /**
   * Enforces that field2 is valid before Worker proceeds
   */
  static raise_if_invalid_field2(state: FeatureNameState): void {
    if (!FeatureInvestigator.is_valid_field2(state)) {
      throw new FeatureError(
        "Invalid field2 value: must be between 0 and 10000",
        state
      );
    }
  }

  /**
   * Enforces authorization before Worker proceeds
   */
  static raise_if_unauthorized(state: FeatureNameState, userId: string): void {
    if (!FeatureInvestigator.can_user_modify(state, userId)) {
      throw new FeatureError(
        "User not authorized to modify this resource",
        state
      );
    }
  }
}
```

#### Workers (State Transformations)

**File**: `src/lib/workers/{feature}-worker.ts`

```typescript
import { FeatureNameState, createFeatureNameState } from "@/types/states/{feature}-state";
import { FeatureEnforcer } from "@/lib/enforcers/{feature}-enforcer";
import { getFirestoreClient } from "@/lib/firestore/firestore-client";

export class FeatureWorker {
  /**
   * Creates a new feature item (state transformation)
   */
  static async create_feature(
    data: Omit<FeatureNameState, "id" | "createdAt" | "updatedAt" | "previous_state">
  ): Promise<FeatureNameState> {
    const now = new Date();

    // Create initial state
    const initialState = createFeatureNameState({
      ...data,
      id: "", // Will be set after Firestore create
      createdAt: now,
      updatedAt: now,
    });

    // Enforce business rules before proceeding
    FeatureEnforcer.raise_if_invalid_field2(initialState);

    // Call Firestore service (side effect)
    const db = getFirestoreClient();
    const docRef = await db.collection("features").add({
      field1: initialState.field1,
      field2: initialState.field2,
      userId: initialState.userId,
      createdAt: now,
      updatedAt: now,
    });

    // Return new state with ID
    return createFeatureNameState({
      ...initialState,
      id: docRef.id,
    });
  }

  /**
   * Updates an existing feature item (state transformation)
   */
  static async update_feature(
    currentState: FeatureNameState,
    updates: Partial<Pick<FeatureNameState, "field1" | "field2">>,
    userId: string
  ): Promise<FeatureNameState> {
    // Enforce business rules
    FeatureEnforcer.raise_if_unauthorized(currentState, userId);

    const updatedState = createFeatureNameState({
      ...currentState,
      ...updates,
      updatedAt: new Date(),
      previous_state: currentState, // Track state history
    });

    FeatureEnforcer.raise_if_invalid_field2(updatedState);

    // Call Firestore service (side effect)
    const db = getFirestoreClient();
    await db.collection("features").doc(currentState.id).update({
      ...updates,
      updatedAt: updatedState.updatedAt,
    });

    return updatedState;
  }
}
```

#### Delegators (Workflow Orchestration)

**File**: `src/lib/delegators/{feature}-delegator.ts`

```typescript
import { FeatureNameState } from "@/types/states/{feature}-state";
import { FeatureWorker } from "@/lib/workers/{feature}-worker";

export class FeatureDelegator {
  /**
   * Orchestrates the complete workflow for creating a feature
   * Assumes happy path - exceptions propagate to caller
   */
  async process_create(
    data: Omit<FeatureNameState, "id" | "createdAt" | "updatedAt" | "previous_state">
  ): Promise<FeatureNameState> {
    // Simply coordinate Workers
    // No business logic, no error handling
    return await FeatureWorker.create_feature(data);
  }

  /**
   * Orchestrates the complete workflow for updating a feature
   */
  async process_update(
    currentState: FeatureNameState,
    updates: Partial<Pick<FeatureNameState, "field1" | "field2">>,
    userId: string
  ): Promise<FeatureNameState> {
    // Coordinate Workers in correct order
    return await FeatureWorker.update_feature(currentState, updates, userId);
  }
}
```

#### Server Actions (Call Delegators)

**File**: `src/app/actions/{feature}.ts`

```typescript
"use server";

import { revalidatePath } from "next/cache";
import { verifySession } from "@/lib/auth/session";
import { FeatureDelegator } from "@/lib/delegators/{feature}-delegator";
import { FeatureError } from "@/lib/enforcers/{feature}-enforcer";
import { featureNameSchema } from "@/types/states/{feature}-state";

const delegator = new FeatureDelegator();

export async function createFeatureAction(formData: FormData) {
  try {
    // Verify authentication
    const session = await verifySession();
    if (!session) {
      return { success: false, error: "Unauthorized" };
    }

    // Parse and validate input
    const data = {
      field1: formData.get("field1") as string,
      field2: Number(formData.get("field2")),
      userId: session.uid,
    };

    const validated = featureNameSchema.parse(data);

    // Call Delegator (handles business logic via Team Pattern)
    const newState = await delegator.process_create(validated);

    // Revalidate the page
    revalidatePath("/{feature-route}");

    return { success: true, data: newState };
  } catch (error) {
    // Handle typed exceptions from Enforcers
    if (error instanceof FeatureError) {
      return {
        success: false,
        error: error.message,
        error_state: error.error_state,
      };
    }

    // Handle other errors
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error",
    };
  }
}

export async function updateFeatureAction(
  currentStateJson: string,
  formData: FormData
) {
  try {
    const session = await verifySession();
    if (!session) {
      return { success: false, error: "Unauthorized" };
    }

    // Parse current state
    const currentState = JSON.parse(currentStateJson);

    // Parse updates
    const updates = {
      field1: formData.get("field1") as string,
      field2: Number(formData.get("field2")),
    };

    // Call Delegator
    const updatedState = await delegator.process_update(
      currentState,
      updates,
      session.uid
    );

    revalidatePath("/{feature-route}");

    return { success: true, data: updatedState };
  } catch (error) {
    if (error instanceof FeatureError) {
      return {
        success: false,
        error: error.message,
        error_state: error.error_state,
      };
    }

    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error",
    };
  }
}
```

#### Page Route (Server Component)

**File**: `src/app/{feature}/page.tsx`

```typescript
import { verifySession } from "@/lib/auth/session";
import { redirect } from "next/navigation";
import { getFeatures } from "@/lib/firestore/{feature}-service";
import FeatureList from "@/components/FeatureList";
import CreateFeatureForm from "@/components/CreateFeatureForm";

export default async function FeaturePage() {
  const session = await verifySession();
  if (!session) {
    redirect("/login");
  }

  const features = await getFeatures(session.uid);

  return (
    <main className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-6">Feature Name</h1>

      <div className="mb-8">
        <CreateFeatureForm />
      </div>

      {features.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          <p>No items yet. Create your first one!</p>
        </div>
      ) : (
        <FeatureList features={features} />
      )}
    </main>
  );
}
```

#### Client Components (for interactivity)

**File**: `src/components/CreateFeatureForm.tsx`

```typescript
"use client";

import { useState } from "react";
import { createFeatureAction } from "@/app/actions/{feature}";

export default function CreateFeatureForm() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const formData = new FormData(e.currentTarget);
      await createFeatureAction(formData);

      // Reset form
      e.currentTarget.reset();
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setLoading(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="field1" className="block text-sm font-medium mb-1">
          Field 1
        </label>
        <input
          id="field1"
          name="field1"
          type="text"
          required
          className="w-full border rounded px-3 py-2"
        />
      </div>

      <div>
        <label htmlFor="field2" className="block text-sm font-medium mb-1">
          Field 2
        </label>
        <input
          id="field2"
          name="field2"
          type="number"
          required
          className="w-full border rounded px-3 py-2"
        />
      </div>

      {error && (
        <div className="bg-red-50 text-red-600 p-3 rounded">
          {error}
        </div>
      )}

      <button
        type="submit"
        disabled={loading}
        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
      >
        {loading ? "Creating..." : "Create"}
      </button>
    </form>
  );
}
```

### Step 4: Quality Checklist

Before considering the feature complete:

#### Functionality
- [ ] All acceptance criteria from PRD are met
- [ ] Feature works as expected for happy path
- [ ] Edge cases are handled
- [ ] Error cases display appropriate messages

#### Code Quality
- [ ] TypeScript types are complete and accurate
- [ ] No `any` types (unless absolutely necessary)
- [ ] Server Components used by default
- [ ] Client Components only where interactivity needed
- [ ] Proper error boundaries
- [ ] Loading states on all async operations
- [ ] Form validation (client and server)

#### Security
- [ ] Authentication checks on all routes
- [ ] Authorization checks (user can only access their data)
- [ ] Input validation with Zod
- [ ] No SQL/NoSQL injection risks
- [ ] XSS prevention (React handles most of this)

#### UX
- [ ] Loading states with visual feedback
- [ ] Error states with clear messages
- [ ] Empty states with helpful guidance
- [ ] Success feedback after actions
- [ ] Responsive on mobile, tablet, desktop
- [ ] Keyboard navigation works
- [ ] Touch targets are 44x44px minimum

#### Performance
- [ ] No unnecessary client-side JavaScript
- [ ] Images optimized with next/image
- [ ] Data fetched server-side when possible
- [ ] No prop drilling or excessive context
- [ ] Efficient Firestore queries

#### Testing
- [ ] Feature works locally (`npm run dev`)
- [ ] No TypeScript errors (`npm run type-check`)
- [ ] No linting errors (`npm run lint`)
- [ ] Production build succeeds (`npm run build`)
- [ ] Feature tested on deployed branch

### Step 5: Integration

After implementation:

1. **Update Navigation**:
   Add links to the new feature in relevant places:
   - Main navigation
   - Dashboard
   - Related features

2. **Update Documentation**:
   - Add feature to README
   - Update architecture docs if needed
   - Document any new environment variables

3. **Test End-to-End**:
   - Create items
   - Update items
   - Delete items
   - Test error cases
   - Test on mobile

4. **Deploy to Feature Branch**:
   ```bash
   git checkout -b feature/{feature-name}
   git add .
   git commit -m "feat: implement {feature-name}"
   git push origin feature/{feature-name}
   ```

5. **Test in Production Environment**:
   - Wait for deployment
   - Test at feature branch URL
   - Verify Firestore data
   - Check Cloud Run logs

6. **Request Review** (if applicable):
   - Code review
   - UX review
   - Security review

### Step 6: Merge to Master

Once everything is verified:

```bash
git checkout master
git merge feature/{feature-name}
git push origin master
```

## Common Patterns

### Loading State Pattern
```typescript
import { Suspense } from "react";

export default function Page() {
  return (
    <Suspense fallback={<LoadingSkeleton />}>
      <DataComponent />
    </Suspense>
  );
}
```

### Error Handling Pattern
```typescript
try {
  await riskyOperation();
} catch (error) {
  console.error("Operation failed:", error);
  return {
    success: false,
    error: error instanceof Error ? error.message : "Unknown error",
  };
}
```

### Form Validation Pattern
```typescript
"use client";

import { useFormState } from "react-dom";
import { submitAction } from "@/app/actions";

export default function MyForm() {
  const [state, formAction] = useFormState(submitAction, null);

  return (
    <form action={formAction}>
      {/* form fields */}
      {state?.error && <div className="error">{state.error}</div>}
    </form>
  );
}
```

## Reference Implementations

Look at existing features in the codebase for patterns:
- Authentication: `src/app/login`, `src/lib/auth/`
- Data fetching: Other feature pages
- Forms: Existing form components
- Styling: Existing components

Follow the established patterns for consistency.

---

Now, please specify:
1. Which feature from the PRD should I implement?
2. Any specific requirements or constraints?
3. Should I start with data layer, server layer, or UI layer?

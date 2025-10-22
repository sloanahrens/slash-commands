---
description: Guide through GCP infrastructure setup and deployment configuration for an app
---

Set up GCP infrastructure and deployment configuration for: [APP_NAME]

**Prerequisites**:
- App code is complete and tested locally
- PRD and technical spec are available
- You have GCP project access
- gcloud CLI is installed and authenticated

## Phase 1: Environment Configuration

### Step 1: Gather Required Information

Collect the following:
- **GCP Project ID**: [WAIT FOR USER INPUT]
- **App Name**: [WAIT FOR USER INPUT]
- **GCP Region**: [USER INPUT or default: us-central1]
- **Deployment SA Email**: [From devops-cloud-run setup]

### Step 2: Create .env.gcp File

Create `~/code/{app-name}/.env.gcp`:

```bash
# GCP Configuration
export GCP_PROJECT="your-project-id"
export GCP_REGION="us-central1"

# Service Account Configuration
export SERVICE_ACCOUNT_NAME="{app-name}-app"
export DEPLOYMENT_SERVICE_ACCOUNT_EMAIL="deploy-sa@project.iam.gserviceaccount.com"

# Firestore Configuration
export FIRESTORE_LOCATION="nam5"
```

Load environment:
```bash
cd ~/code/{app-name}
source .env.gcp
```

## Phase 2: GCP Infrastructure Setup

### Step 3: Run GCP Setup Script

This script will:
- Create dedicated service account for the app
- Enable required APIs (Identity Platform, Firestore, Secret Manager)
- Create Firestore database
- Configure IAM permissions
- Create secrets in Secret Manager
- Grant deployment SA access to secrets

```bash
cd ~/code/{app-name}
bash scripts/setup-gcp.sh
```

Review the output and confirm:
- ✓ Service account created: {app-name}-app@{PROJECT}.iam.gserviceaccount.com
- ✓ APIs enabled
- ✓ Firestore database created
- ✓ IAM permissions configured
- ✓ Secrets created

### Step 4: Manual Identity Platform Configuration

The script will prompt you to complete these manual steps:

1. **Configure OAuth Consent Screen**:
   - Go to: https://console.cloud.google.com/apis/credentials/consent?project={PROJECT}
   - Select "Internal" (if using Google Workspace) or "External"
   - Fill in app name, support email, developer contact
   - No scopes needed for basic auth
   - Save

2. **Enable Identity Platform and Google Provider**:
   - Go to: https://console.cloud.google.com/customer-identity?project={PROJECT}
   - Click "Enable"
   - Go to "Providers" tab
   - Enable "Google" provider
   - Save

3. **Get Firebase API Key**:
   - Go to: https://console.cloud.google.com/customer-identity/settings?project={PROJECT}
   - Copy the "API Key" shown
   - Save this for the next step

4. **Add Authorized Domains**:
   - In Identity Platform settings, go to "Settings" tab
   - Under "Authorized domains", add:
     - `localhost` (for local development)
     - Your Cloud Run domain (will be available after first deployment)

### Step 5: Create Firestore Indexes

If your app needs composite indexes:

```bash
cd ~/code/{app-name}
bash scripts/create-firestore-indexes.sh
```

Verify indexes are building:
```bash
gcloud firestore indexes composite list --project=$GCP_PROJECT
```

## Phase 3: Bitbucket Configuration

### Step 6: Create Bitbucket Repository

1. Create new repository in Bitbucket:
   - Name: {app-name}
   - Private repository
   - Initialize without README (we have one)

2. Add remote and push:
```bash
cd ~/code/{app-name}
git remote add origin git@bitbucket.org:your-org/{app-name}.git
git push -u origin master
```

### Step 7: Configure Repository Variables

In Bitbucket, go to: Repository Settings → Repository variables

Add these variables:

**Required Variables**:
```
GCP_PROJECT = your-project-id
GCP_SERVICE_ACCOUNT_KEY = [base64-encoded service account key from devops-cloud-run]
```

**Optional Variables**:
```
GCP_REGION = us-central1
```

### Step 8: Configure Firebase Build Args

Update `bitbucket-pipelines.yml` with your Firebase API key:

```yaml
# Line ~133 in bitbucket-pipelines.yml
export DOCKER_BUILD_ARGS="--build-arg NEXT_PUBLIC_FIREBASE_API_KEY=YOUR_API_KEY_HERE --build-arg NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=${GCP_PROJECT}.firebaseapp.com --build-arg NEXT_PUBLIC_FIREBASE_PROJECT_ID=${GCP_PROJECT}"
```

Replace `YOUR_API_KEY_HERE` with the API key from Step 4.

Commit this change:
```bash
git add bitbucket-pipelines.yml
git commit -m "chore: configure Firebase API key for deployment"
git push origin master
```

## Phase 4: First Deployment

### Step 9: Trigger Initial Deployment

The push to master will automatically trigger the Bitbucket Pipeline.

Monitor the pipeline:
- Go to: Bitbucket → Pipelines
- Watch both parallel steps: "Build and Test" and "Deploy to Cloud Run"

Expected timeline:
- Build/test: 2-5 minutes
- Deployment: 5-8 minutes
- Total: ~10-15 minutes

### Step 10: Verify Deployment

Once the pipeline completes:

1. **Get Service URL**:
   Look for output in the pipeline logs:
   ```
   Service URL is https://{app-name}-master-xxxxx-uc.a.run.app
   ```

2. **Test Health Endpoint**:
   ```bash
   curl https://{app-name}-master-xxxxx-uc.a.run.app/api/health
   ```

   Expected response:
   ```json
   {
     "status": "healthy",
     "timestamp": "2024-10-20T...",
     "service": "{app-name}"
   }
   ```

3. **Test in Browser**:
   - Visit the service URL
   - Should redirect to `/login`
   - Try Google Sign In
   - Should redirect to dashboard after auth

### Step 11: Update Authorized Domains

Now that you have the Cloud Run URL, add it to Identity Platform:

1. Go to: https://console.cloud.google.com/customer-identity/settings?project={PROJECT}
2. Under "Authorized domains", click "Add domain"
3. Add: `{app-name}-master-xxxxx-uc.a.run.app` (your actual domain)
4. Save

## Phase 5: Local Development Setup

### Step 12: Configure Local Environment

Create `nextapp/.env.local`:

```bash
# Firebase Configuration (Public)
NEXT_PUBLIC_FIREBASE_API_KEY=your-api-key-here
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN={PROJECT_ID}.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID={PROJECT_ID}

# GCP Configuration
GOOGLE_CLOUD_PROJECT={PROJECT_ID}

# Local Development Secrets
JWT_SECRET=your-local-jwt-secret-for-development

# Note: For Firestore and Secret Manager access, use ADC:
# gcloud auth application-default login
# gcloud config set project {PROJECT_ID}
```

### Step 13: Install Dependencies and Test

```bash
cd ~/code/{app-name}/nextapp
npm install
npm run build          # Verify it builds
npm run type-check     # Verify no TypeScript errors
npm run lint           # Verify no linting errors
npm run dev            # Start dev server
```

Visit http://localhost:3000 and test locally.

## Phase 6: Branch Deployment (Optional)

### Step 14: Test Feature Branch Deployment

Create and deploy a feature branch:

```bash
git checkout -b feature/test-deployment
git push origin feature/test-deployment
```

This creates a separate Cloud Run service:
- Service name: `{app-name}-feature-test-deployment`
- Separate URL
- Isolated from production

Test the feature branch, then clean up:

```bash
# In Bitbucket: Pipelines → Run pipeline
# Select: Custom: branch-cleanup
# Branch: feature/test-deployment
```

## Phase 7: Monitoring & Maintenance

### Step 15: Set Up Monitoring

1. **View Cloud Run Logs**:
   ```bash
   gcloud run services logs read {app-name}-master --project=$GCP_PROJECT --limit=50
   ```

2. **Monitor Firestore Usage**:
   - Go to: https://console.cloud.google.com/firestore/usage?project={PROJECT}
   - Check document reads/writes
   - Monitor storage usage

3. **Monitor Cloud Run Metrics**:
   - Go to: https://console.cloud.google.com/run?project={PROJECT}
   - Click on your service
   - View: Request count, latency, instance count, CPU/memory usage

### Step 16: Set Up Alerts (Optional)

Create alerts for:
- Error rate spikes
- Unusual traffic patterns
- Cost anomalies
- Cold start performance

## Troubleshooting

### Pipeline Fails with "Image not found"
- Check that devops-cloud-run repository is accessible
- Verify GCP_SERVICE_ACCOUNT_KEY is set correctly
- Check that Artifact Registry is created (bootstrap step)

### Authentication Doesn't Work
- Verify Firebase API key in build args
- Check Identity Platform is enabled
- Ensure authorized domains include your Cloud Run domain
- Check service account has `identitytoolkit.viewer` role

### Firestore Queries Fail
- Check service account has `datastore.user` role
- Verify Firestore database is created
- Check composite indexes are built (not "Building")
- Review Firestore security rules

### Secrets Not Accessible
- Verify secrets exist in Secret Manager
- Check service account has `secretAccessor` on each secret
- Verify deployment SA also has access (for Cloud Run config)
- Check secret names in bitbucket-pipelines.yml match Secret Manager

### Health Check Fails
- Check Docker build succeeded
- Verify port 3000 is exposed
- Check that /api/health route exists
- Review Cloud Run logs for startup errors

## Success Checklist

Before considering deployment complete:

- [ ] GCP infrastructure created successfully
- [ ] Identity Platform configured and working
- [ ] Firestore database created and accessible
- [ ] All secrets created in Secret Manager
- [ ] Service account has correct permissions
- [ ] Bitbucket repository configured
- [ ] CI/CD pipeline runs successfully
- [ ] Health check endpoint responds
- [ ] Authentication flow works end-to-end
- [ ] All primary features work in production
- [ ] No errors in Cloud Run logs
- [ ] Authorized domains configured
- [ ] Local development environment working
- [ ] Documentation updated with URLs

## Next Steps

After successful deployment:

1. **Custom Domain** (optional):
   ```bash
   cd ~/code/devops-cloud-run/devops-cloud-run/cloud-run
   export TF_VAR_custom_domain="app.yourdomain.com"
   ./cloud-run.sh apply master
   ```

2. **Monitoring Dashboard**: Set up custom monitoring in GCP Console

3. **Cost Tracking**: Tag resources and set up budget alerts

4. **Backup Strategy**: Configure Firestore backups

5. **Security Review**: Run security scan and review IAM policies

Your app is now deployed and production-ready! 🚀

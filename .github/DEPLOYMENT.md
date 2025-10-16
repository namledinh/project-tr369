# GitHub repository settings and secrets configuration

## Required GitHub Secrets

### For Container Registry (GitHub Container Registry)
- `GITHUB_TOKEN` - Automatically provided by GitHub Actions

### For Production Deployment (if needed)
- `DEPLOY_HOST` - Production server hostname/IP
- `DEPLOY_USER` - SSH username for deployment
- `DEPLOY_KEY` - SSH private key for deployment
- `DATABASE_URL` - Production database connection string

### For External Services (if needed)
- `CODECOV_TOKEN` - For code coverage reporting (optional, get from codecov.io)
- `SLACK_WEBHOOK` - For notifications (optional)

## Environment Configuration

### Staging Environment
- Branch: `develop`
- Auto-deploy on successful CI

### Production Environment  
- Branch: `main`
- Manual approval required
- Protected branch rules recommended

## Branch Protection Rules (Recommended)

### For `main` branch:
- Require pull request reviews before merging
- Require status checks to pass before merging
- Require branches to be up to date before merging
- Include administrators in restrictions

### For `develop` branch:
- Require status checks to pass before merging
- Require branches to be up to date before merging

## Docker Image Registry

Images will be pushed to GitHub Container Registry:
- Registry: `ghcr.io`
- Image name: `ghcr.io/namledinh/project-tr369`
- Tags: branch names, PR numbers, and semantic versions

## Deployment Strategy

1. **Feature branches** → Create PR → Run CI
2. **develop branch** → Auto-deploy to staging after CI passes
3. **main branch** → Manual deploy to production after CI passes
4. **Tags (v*.*)** → Create release with binaries and Docker images
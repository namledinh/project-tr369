#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="https://github.com/namledinh/project-tr369.git"
CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "main")

echo -e "${BLUE}🚀 Setting up CI/CD for USP Management Device API${NC}"
echo "=================================================="

# Function to print colored messages
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "This is not a git repository. Initializing..."
    git init
    print_success "Git repository initialized"
fi

# Add all files to git
print_info "Adding files to git..."
git add .

# Check if there are changes to commit
if git diff --staged --quiet; then
    print_warning "No changes to commit"
else
    # Commit changes
    print_info "Committing changes..."
    git commit -m "feat: add CI/CD configuration with GitHub Actions

- Add comprehensive CI pipeline with testing, security scanning, and Docker builds
- Add release automation with multi-platform binary builds
- Add CodeQL security analysis and dependency review
- Add Dependabot configuration for automated dependency updates
- Add development tools configuration (Air for hot reload)
- Update README with CI/CD badges and deployment instructions
- Add GitHub issue templates and security policy"
    
    print_success "Changes committed successfully"
fi

# Check if remote origin exists
if ! git remote get-url origin > /dev/null 2>&1; then
    print_info "Adding GitHub remote..."
    git remote add origin $GITHUB_REPO
    print_success "Remote origin added: $GITHUB_REPO"
fi

# Push to GitHub
print_info "Pushing to GitHub repository..."
echo -e "${YELLOW}Current branch: $CURRENT_BRANCH${NC}"

# Push the current branch
if git push -u origin $CURRENT_BRANCH; then
    print_success "Code pushed to GitHub successfully!"
else
    print_error "Failed to push to GitHub. Please check your permissions."
    echo ""
    print_info "Manual steps to push:"
    echo "1. Make sure you have write access to the repository"
    echo "2. Check your Git credentials"
    echo "3. Try: git push -u origin $CURRENT_BRANCH"
    exit 1
fi

echo ""
print_success "🎉 CI/CD setup completed!"
echo "=================================================="
print_info "Next steps:"
echo "1. Go to your GitHub repository: https://github.com/namledinh/project-tr369"
echo "2. Check the Actions tab to see the CI pipeline running"
echo "3. Set up branch protection rules for main/develop branches"
echo "4. Configure environments (staging/production) in repository settings"
echo "5. Add required secrets if needed for deployment"

echo ""
print_info "Available GitHub Actions workflows:"
echo "• CI Pipeline (.github/workflows/ci.yml) - Runs on push/PR"
echo "• Release (.github/workflows/release.yml) - Runs on git tags"  
echo "• Security Scan (.github/workflows/codeql.yml) - Weekly & on push"
echo "• Dependency Review (.github/workflows/dependency-review.yml) - On PRs"

echo ""
print_info "To create a release:"
echo "git tag -a v1.0.0 -m 'First release'"
echo "git push origin v1.0.0"

echo ""
print_success "Happy coding! 🚀"
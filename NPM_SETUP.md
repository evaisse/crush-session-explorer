# NPM Publishing Setup

This document explains how to set up NPM publishing for this project.

## Prerequisites

1. An NPM account (create one at https://www.npmjs.com/signup if you don't have one)
2. Admin access to this GitHub repository
3. The package name `crush-session-explorer` should be available on NPM (or you control it)

## Setup Instructions

### 1. Generate NPM Access Token

1. Log in to your NPM account at https://www.npmjs.com
2. Click on your profile picture → "Access Tokens"
3. Click "Generate New Token" → "Classic Token"
4. Select "Automation" type (recommended for CI/CD)
5. Copy the generated token (you won't be able to see it again!)

### 2. Add Token to GitHub Repository

1. Go to your GitHub repository settings
2. Navigate to "Secrets and variables" → "Actions"
3. Click "New repository secret"
4. Name: `NPM_TOKEN`
5. Value: Paste the NPM token you generated
6. Click "Add secret"

## Publishing a New Release

Once the NPM_TOKEN is configured, publishing is automatic:

1. Update the version in `package.json` to match your release version
2. Commit and push the changes
3. Create and push a git tag:
   ```bash
   git tag v1.0.2
   git push origin v1.0.2
   ```
4. The GitHub Actions workflow will:
   - Run tests
   - Build binaries for all platforms
   - Create a GitHub release with binaries
   - Publish the package to NPM

## Verifying the Publication

After the workflow completes, verify the package was published:

```bash
npm view crush-session-explorer
```

Users can then install it with:

```bash
npm install -g crush-session-explorer
```

## Troubleshooting

### Package name already taken
If the package name is already taken on NPM, you can either:
- Use a scoped package name (e.g., `@your-username/crush-session-explorer`)
- Choose a different name

Update the `name` field in `package.json` accordingly.

### NPM_TOKEN not working
- Ensure the token type is "Automation"
- Verify the token has not expired
- Check that the token is properly added to GitHub secrets (no extra spaces)

### Version conflicts
- The workflow uses `--allow-same-version` to prevent version conflicts
- Ensure the version in `package.json` matches the git tag (without the 'v' prefix)

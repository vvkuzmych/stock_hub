# Security Guidelines 🔒

Правила безпеки для Stock Hub monorepo.

---

## ⚠️ НІКОЛИ НЕ КОМІТЬ В GIT:

### ❌ Credentials & Secrets

```bash
# Environment files
.env
*.env.local
*.env.production

# API Keys
*api_key*
*apikey*
*api-key*

# Passwords
*password*
*passwd*

# Tokens
*token*
*jwt*
*secret*

# SSH Keys
id_rsa
id_ed25519
*.pem
*.key

# Database
database.db
*.sqlite
*.sqlite3

# Certificates
*.crt
*.p12
*.pfx
```

---

## ✅ Що МОЖНА коміть:

### Examples & Templates

```bash
✅ .env.example           # Template without real values
✅ .env.template          # Template
✅ config.example.yaml    # Example config
✅ credentials.example    # Example credentials
```

**Правило:** Якщо файл містить `example`, `template`, `sample` - можна.

---

## 🛡️ Захист Credentials

### 1. Environment Variables

```bash
# .env (gitignored)
SMTP_USER=real-user@ethereal.email
SMTP_PASSWORD=real-password
DATABASE_DSN=postgres://user:pass@localhost/db

# Load in app
config.Load() // reads from .env
```

### 2. .env.example

```bash
# .env.example (committed to git)
SMTP_USER=your-username-here
SMTP_PASSWORD=your-password-here
DATABASE_DSN=postgres://user:pass@localhost/db

# Usage:
cp .env.example .env
nano .env  # Edit with real values
```

### 3. Git Hooks (pre-commit)

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Check for potential secrets
if git diff --cached --name-only | grep -E "\.env$|credentials|password|secret"; then
    echo "⚠️  Warning: You're about to commit sensitive files!"
    echo "Are you sure? (y/n)"
    read answer
    if [ "$answer" != "y" ]; then
        exit 1
    fi
fi
```

---

## 📋 Security Checklist

### Before Every Commit

- [ ] `.env` файли в `.gitignore`?
- [ ] Credentials НЕ в коді?
- [ ] Passwords НЕ в commit message?
- [ ] API keys НЕ в конфігах?
- [ ] Database credentials через env vars?

### Before Push

```bash
# Check what you're pushing
git log origin/develop..HEAD --oneline

# Search for secrets
git log -p | grep -E "password|token|api_key|secret"

# If found secrets - STOP!
```

---

## 🚨 Що Робити Якщо Закомітив Secrets?

### Option 1: Amend Last Commit (якщо НЕ push)

```bash
# 1. Remove sensitive file
git rm --cached .env

# 2. Add to .gitignore
echo ".env" >> .gitignore
git add .gitignore

# 3. Amend commit
git commit --amend --no-edit

# 4. Force push (ONLY if you're alone on branch)
git push --force
```

### Option 2: BFG Repo Cleaner (якщо вже push)

```bash
# 1. Install BFG
brew install bfg

# 2. Remove file from history
bfg --delete-files .env

# 3. Clean up
git reflog expire --expire=now --all
git gc --prune=now --aggressive

# 4. Force push
git push --force
```

### Option 3: Rotate Credentials (ЗАВЖДИ!)

**Навіть якщо видалив з git - credentials вже скомпрометовані!**

```bash
# 1. Change all passwords immediately
# 2. Revoke API keys
# 3. Generate new credentials
# 4. Update .env
```

---

## 🔐 Best Practices

### 1. Use Environment Variables

```go
// ❌ BAD
const password = "my-secret-password"

// ✅ GOOD
password := os.Getenv("PASSWORD")
```

### 2. Use Secret Management

```bash
# Vault, AWS Secrets Manager, etc.
vault kv get secret/database
```

### 3. Separate Configs per Environment

```
config/
├── .env.development     # Local dev
├── .env.staging         # Staging
├── .env.production      # Production (never commit!)
└── .env.example         # Template (commit this)
```

### 4. Use Different Credentials per Environment

```bash
# Development
SMTP_USER=dev@ethereal.email

# Staging
SMTP_USER=staging@ethereal.email

# Production
SMTP_USER=noreply@realcompany.com
```

---

## 🎯 Stock Hub Specific

### Email Service

```bash
# ✅ Safe (testing)
SMTP_HOST=smtp.ethereal.email
SMTP_USER=test123@ethereal.email
# These are temporary test credentials

# ❌ NEVER commit production Gmail
SMTP_USER=company@gmail.com
SMTP_PASSWORD=real-app-password
```

### Database

```bash
# ✅ Development
DATABASE_DSN=postgres://postgres:postgres@localhost/stock_hub

# ❌ Production (use env vars)
DATABASE_DSN=postgres://admin:super-secret@prod-db.aws.com/stock_hub
```

### Trading Service

```bash
# ✅ Test mode
MOCK_EMAIL=true

# ⚠️  Production
MOCK_EMAIL=false
SMTP_PASSWORD=$PROD_SMTP_PASSWORD  # From env
```

---

## 📚 Learn More

- **Git Secrets**: https://github.com/awslabs/git-secrets
- **BFG Repo Cleaner**: https://rtyley.github.io/bfg-repo-cleaner/
- **OWASP**: https://owasp.org/www-community/vulnerabilities/

---

## ✅ Current Status: SECURE

- ✅ `.env` файли в `.gitignore`
- ✅ `.env.example` без реальних credentials
- ✅ Ethereal credentials тимчасові (24h TTL)
- ✅ No secrets in committed code
- ✅ Environment-based configuration

**Keep it secure!** 🔒

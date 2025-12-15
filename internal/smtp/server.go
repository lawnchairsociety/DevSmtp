package smtp

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"strings"

	"github.com/lawnchairsociety/devsmtp/internal/config"
	"github.com/lawnchairsociety/devsmtp/internal/database"
)

type Server struct {
	config    *config.Config
	db        *database.DB
	logger    *Logger
	tlsConfig *tls.Config
}

func NewServer(cfg *config.Config, db *database.DB, logger *Logger) *Server {
	s := &Server{
		config: cfg,
		db:     db,
		logger: logger,
	}

	if cfg.TLS.Cert != "" && cfg.TLS.Key != "" {
		cert, err := tls.LoadX509KeyPair(cfg.TLS.Cert, cfg.TLS.Key)
		if err != nil {
			s.logger.Warn("Failed to load TLS certificates: %v", err)
		} else {
			s.tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
			s.logger.Info("TLS certificates loaded successfully")
		}
	}

	return s
}

func (s *Server) ListenAndServe() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	defer listener.Close()

	s.logger.Info("SMTP server listening on %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

type session struct {
	server        *Server
	conn          net.Conn
	reader        *bufio.Reader
	writer        *bufio.Writer
	clientIP      string
	helo          string
	mailFrom      string
	rcptTo        []string
	data          []byte
	authenticated bool
	tlsActive     bool
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientIP := conn.RemoteAddr().String()
	if host, _, err := net.SplitHostPort(clientIP); err == nil {
		clientIP = host
	}

	s.logger.Info("New connection from %s", clientIP)

	sess := &session{
		server:   s,
		conn:     conn,
		reader:   bufio.NewReader(conn),
		writer:   bufio.NewWriter(conn),
		clientIP: clientIP,
		rcptTo:   make([]string, 0),
	}

	sess.writeLine("220 DevSmtp ESMTP Service Ready")

	for {
		line, err := sess.reader.ReadString('\n')
		if err != nil {
			s.logger.Info("Connection closed from %s", clientIP)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		cmd := strings.ToUpper(strings.Split(line, " ")[0])
		args := ""
		if idx := strings.Index(line, " "); idx > 0 {
			args = strings.TrimSpace(line[idx+1:])
		}

		s.logger.Debug("[%s] C: %s %s", clientIP, cmd, args)

		quit := sess.handleCommand(cmd, args)
		if quit {
			return
		}
	}
}

func (sess *session) handleCommand(cmd, args string) bool {
	switch cmd {
	case "HELO":
		sess.handleHelo(args)
	case "EHLO":
		sess.handleEhlo(args)
	case "MAIL":
		sess.handleMailFrom(args)
	case "RCPT":
		sess.handleRcptTo(args)
	case "DATA":
		sess.handleData()
	case "RSET":
		sess.handleRset()
	case "NOOP":
		sess.writeLine("250 OK")
	case "QUIT":
		sess.writeLine("221 Bye")
		sess.server.logger.Info("Connection closed by client %s", sess.clientIP)
		return true
	case "VRFY":
		sess.writeLine("252 Cannot VRFY user, but will accept message")
	case "EXPN":
		sess.writeLine("252 Cannot expand mailing list")
	case "STARTTLS":
		sess.handleStartTLS()
	case "AUTH":
		sess.handleAuth(args)
	default:
		sess.server.logger.Warn("[%s] Unknown command: %s", sess.clientIP, cmd)
		sess.writeLine("502 Command not implemented")
	}

	return false
}

func (sess *session) handleHelo(args string) {
	if args == "" {
		sess.writeLine("501 Syntax: HELO hostname")
		return
	}
	sess.helo = args
	sess.server.logger.Info("[%s] HELO %s", sess.clientIP, args)
	sess.writeLine("250 Hello " + args)
}

func (sess *session) handleEhlo(args string) {
	if args == "" {
		sess.writeLine("501 Syntax: EHLO hostname")
		return
	}
	sess.helo = args
	sess.server.logger.Info("[%s] EHLO %s", sess.clientIP, args)

	sess.writeLine("250-Hello " + args)
	sess.writeLine("250-SIZE 10485760")
	sess.writeLine("250-8BITMIME")
	sess.writeLine("250-PIPELINING")

	if sess.server.tlsConfig != nil && !sess.tlsActive {
		sess.writeLine("250-STARTTLS")
	}

	if sess.server.config.Auth.Username != "" {
		sess.writeLine("250-AUTH PLAIN LOGIN")
	}

	sess.writeLine("250 HELP")
}

func (sess *session) handleMailFrom(args string) {
	if sess.server.config.Auth.Required && !sess.authenticated {
		sess.server.logger.Warn("[%s] AUTH required but not authenticated", sess.clientIP)
		sess.writeLine("530 Authentication required")
		return
	}

	upperArgs := strings.ToUpper(args)
	if !strings.HasPrefix(upperArgs, "FROM:") {
		sess.writeLine("501 Syntax: MAIL FROM:<address>")
		return
	}

	addr := args[5:] // Keep original case
	addr = strings.TrimSpace(addr)
	addr = strings.Trim(addr, "<>")

	sess.mailFrom = addr
	sess.server.logger.Info("[%s] MAIL FROM:<%s>", sess.clientIP, addr)
	sess.writeLine("250 OK")
}

func (sess *session) handleRcptTo(args string) {
	if sess.mailFrom == "" {
		sess.writeLine("503 Need MAIL command first")
		return
	}

	upperArgs := strings.ToUpper(args)
	if !strings.HasPrefix(upperArgs, "TO:") {
		sess.writeLine("501 Syntax: RCPT TO:<address>")
		return
	}

	addr := args[3:] // Keep original case
	addr = strings.TrimSpace(addr)
	addr = strings.Trim(addr, "<>")

	sess.rcptTo = append(sess.rcptTo, addr)
	sess.server.logger.Info("[%s] RCPT TO:<%s>", sess.clientIP, addr)
	sess.writeLine("250 OK")
}

func (sess *session) handleData() {
	if len(sess.rcptTo) == 0 {
		sess.writeLine("503 Need RCPT command first")
		return
	}

	sess.server.logger.Info("[%s] DATA started", sess.clientIP)
	sess.writeLine("354 Start mail input; end with <CRLF>.<CRLF>")

	var dataLines []string
	for {
		line, err := sess.reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "." {
			break
		}

		// Handle dot-stuffing
		if strings.HasPrefix(line, "..") {
			line = line[1:]
		}

		dataLines = append(dataLines, line)
	}

	rawData := strings.Join(dataLines, "\r\n")
	sess.data = []byte(rawData)

	// Parse subject from headers
	subject := ""
	body := rawData
	if idx := strings.Index(rawData, "\r\n\r\n"); idx > 0 {
		headers := rawData[:idx]
		body = rawData[idx+4:]
		for _, line := range strings.Split(headers, "\r\n") {
			if strings.HasPrefix(strings.ToLower(line), "subject:") {
				subject = strings.TrimSpace(line[8:])
				break
			}
		}
	} else if idx := strings.Index(rawData, "\n\n"); idx > 0 {
		headers := rawData[:idx]
		body = rawData[idx+2:]
		for _, line := range strings.Split(headers, "\n") {
			if strings.HasPrefix(strings.ToLower(line), "subject:") {
				subject = strings.TrimSpace(line[8:])
				break
			}
		}
	}

	msg := &database.Message{
		Sender:     sess.mailFrom,
		Recipients: strings.Join(sess.rcptTo, ", "),
		Subject:    subject,
		Body:       body,
		RawData:    sess.data,
		Size:       len(sess.data),
		ClientIP:   sess.clientIP,
		IsRead:     false,
	}

	if err := sess.server.db.SaveMessage(msg); err != nil {
		sess.server.logger.Error("[%s] Failed to save message: %v", sess.clientIP, err)
		sess.writeLine("451 Requested action aborted: local error in processing")
		return
	}

	sess.server.logger.Info("[%s] Message received: %s -> %s (%d bytes) Subject: %s",
		sess.clientIP, sess.mailFrom, strings.Join(sess.rcptTo, ", "), len(sess.data), subject)
	sess.writeLine("250 OK: Message queued")

	// Reset session state for next message
	sess.mailFrom = ""
	sess.rcptTo = make([]string, 0)
	sess.data = nil
}

func (sess *session) handleRset() {
	sess.mailFrom = ""
	sess.rcptTo = make([]string, 0)
	sess.data = nil
	sess.server.logger.Debug("[%s] Session reset", sess.clientIP)
	sess.writeLine("250 OK")
}

func (sess *session) handleStartTLS() {
	if sess.server.tlsConfig == nil {
		sess.writeLine("454 TLS not available")
		return
	}

	if sess.tlsActive {
		sess.writeLine("503 TLS already active")
		return
	}

	sess.server.logger.Info("[%s] STARTTLS initiated", sess.clientIP)
	sess.writeLine("220 Ready to start TLS")

	tlsConn := tls.Server(sess.conn, sess.server.tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		sess.server.logger.Error("[%s] TLS handshake failed: %v", sess.clientIP, err)
		return
	}

	sess.conn = tlsConn
	sess.reader = bufio.NewReader(tlsConn)
	sess.writer = bufio.NewWriter(tlsConn)
	sess.tlsActive = true

	sess.server.logger.Info("[%s] TLS handshake successful", sess.clientIP)

	// Reset session state after STARTTLS
	sess.helo = ""
	sess.mailFrom = ""
	sess.rcptTo = make([]string, 0)
	sess.authenticated = false
}

func (sess *session) handleAuth(args string) {
	if sess.server.config.Auth.Username == "" {
		sess.writeLine("503 Authentication not configured")
		return
	}

	parts := strings.SplitN(args, " ", 2)
	mechanism := strings.ToUpper(parts[0])

	sess.server.logger.Info("[%s] AUTH %s attempted", sess.clientIP, mechanism)

	switch mechanism {
	case "PLAIN":
		sess.handleAuthPlain(parts)
	case "LOGIN":
		sess.handleAuthLogin()
	default:
		sess.writeLine("504 Unrecognized authentication mechanism")
	}
}

func (sess *session) handleAuthPlain(parts []string) {
	var credentials string

	if len(parts) > 1 {
		credentials = parts[1]
	} else {
		sess.writeLine("334 ")
		line, err := sess.reader.ReadString('\n')
		if err != nil {
			return
		}
		credentials = strings.TrimSpace(line)
	}

	decoded, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		sess.writeLine("501 Invalid base64")
		return
	}

	// PLAIN format: \0username\0password
	credParts := strings.Split(string(decoded), "\x00")
	if len(credParts) != 3 {
		sess.writeLine("535 Authentication failed")
		return
	}

	username := credParts[1]
	password := credParts[2]

	if username == sess.server.config.Auth.Username && password == sess.server.config.Auth.Password {
		sess.authenticated = true
		sess.server.logger.Info("[%s] AUTH successful for user: %s", sess.clientIP, username)
		sess.writeLine("235 Authentication successful")
	} else {
		sess.server.logger.Warn("[%s] AUTH failed for user: %s", sess.clientIP, username)
		sess.writeLine("535 Authentication failed")
	}
}

func (sess *session) handleAuthLogin() {
	sess.writeLine("334 VXNlcm5hbWU6") // Base64 for "Username:"

	userLine, err := sess.reader.ReadString('\n')
	if err != nil {
		return
	}
	userDecoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(userLine))
	if err != nil {
		sess.writeLine("501 Invalid base64")
		return
	}

	sess.writeLine("334 UGFzc3dvcmQ6") // Base64 for "Password:"

	passLine, err := sess.reader.ReadString('\n')
	if err != nil {
		return
	}
	passDecoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(passLine))
	if err != nil {
		sess.writeLine("501 Invalid base64")
		return
	}

	username := string(userDecoded)
	password := string(passDecoded)

	if username == sess.server.config.Auth.Username && password == sess.server.config.Auth.Password {
		sess.authenticated = true
		sess.server.logger.Info("[%s] AUTH LOGIN successful for user: %s", sess.clientIP, username)
		sess.writeLine("235 Authentication successful")
	} else {
		sess.server.logger.Warn("[%s] AUTH LOGIN failed for user: %s", sess.clientIP, username)
		sess.writeLine("535 Authentication failed")
	}
}

func (sess *session) writeLine(line string) {
	fmt.Fprintf(sess.writer, "%s\r\n", line)
	sess.writer.Flush()
}

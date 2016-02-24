package main

import (
	"flag"
	"log"
	"os"
)

const usage = "usage: MUNKI_REPO_PATH= APE_HTTP_LISTEN_PORT= ape -repo MUNKI_REPO_PATH -port APE_HTTP_LISTEN_PORT"

var (
	flRepo      = flag.String("repo", envString("APE_MUNKI_REPO_PATH", ""), "path to munki repo")
	flPort      = flag.String("port", envString("APE_HTTP_LISTEN_PORT", ""), "port to listen on")
	flBasic     = flag.Bool("basic", envBool("APE_BASIC_AUTH"), "enable basic auth")
	flJWT       = flag.Bool("jwt", envBool("APE_JWT_AUTH"), "enable jwt authentication for api calls")
	flJWTSecret = flag.String("jwt-signing-key", envString("APE_JWT_SIGNING_KEY", ""), "jwt signing key")
	flTLS       = flag.Bool("tls", envBool("APE_USE_TLS"), "use https")
	flTLSCert   = flag.String("tls-cert", envString("APE_TLS_CERT", ""), "path to TLS certificate")
	flTLSKey    = flag.String("tls-key", envString("APE_TLS_KEY", ""), "path to TLS private key")
)

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

func envBool(key string) bool {
	if env := os.Getenv(key); env == "true" {
		return true
	}
	return false
}

func initFlag() {
	flag.Parse()
	if *flRepo == "" {
		flag.Usage()
		log.Fatal(usage)
	}

	if !dirExists(*flRepo) {
		log.Fatal("MUNKI_REPO_PATH must exist")
	}

	// validate port flag
	if *flPort == "" {
		port := defaultPort()
		log.Printf("no port flag specified. Using %v by default", port)
		*flPort = port
	}

	// validate TLS flags
	if *flTLS {
		checkTLSFlags()
	}

	// validate JWT flags
	if *flJWT {
		checkJWTFlags()
	}

	// validate basic auth
	if *flBasic && !*flJWT {
		log.Fatal("Basic Authentication is used to issue JWT Tokens. You must enable JWT as well")
	}
}

func checkJWTFlags() {
	if *flJWTSecret == "" {
		log.Fatal("You must provide a signing key to enable JWT authentication")
	}
}

func checkTLSFlags() {
	if *flTLSKey == "" || *flTLSCert == "" {
		log.Fatal("You must provide a valid path to a TLS cert and key")
	}
}

func defaultPort() string {
	if *flTLS {
		return "443"
	}
	return "80"
}

package armor

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
)

// GetConfigForClient implements the Config.GetClientCertificate callback
func (a *Armor) GetConfigForClient(clientHelloInfo *tls.ClientHelloInfo) (*tls.Config, error) {
	// Get the host from the hello info
	host := a.Hosts[clientHelloInfo.ServerName]
	if len(host.ClientCAs) == 0 {
		return nil, nil
	}

	// Use existing host config if exist
	if host.TLSConfig != nil {
		return host.TLSConfig, nil
	}

	// Build and save the host config
	host.TLSConfig = a.buildTLSConfig(clientHelloInfo, host)

	return host.TLSConfig, nil
}

func (a *Armor) buildTLSConfig(clientHelloInfo *tls.ClientHelloInfo, host *Host) *tls.Config {
	// Copy the configurations from the regular server
	tlsConfig := new(tls.Config)
	*tlsConfig = *a.Echo.TLSServer.TLSConfig

	// Set the client validation and the certification pool
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	tlsConfig.ClientCAs = a.buildClientCertPool(host)

	return tlsConfig
}

func (a *Armor) buildClientCertPool(host *Host) (certPool *x509.CertPool) {
	certPool = x509.NewCertPool()

	// Loop every CA certs given as base64 DER encoding
	for _, clientCAString := range host.ClientCAs {
		// Decode base64
		derBytes, err := base64.StdEncoding.DecodeString(clientCAString)
		if err != nil {
			continue
		}
		if len(derBytes) == 0 {
			continue
		}

		// Parse the DER encoded certificate
		var caCert *x509.Certificate
		caCert, err = x509.ParseCertificate(derBytes)
		if err != nil {
			continue
		}

		// Add the certificate to CertPool
		certPool.AddCert(caCert)
	}

	return certPool
}

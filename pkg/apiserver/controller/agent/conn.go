package agent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/fristonio/xene/pkg/store"
	"github.com/fristonio/xene/pkg/types/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GetAgentConnection returns a GRPC connection corresponding to the provided
// Agent manifest.
func GetAgentConnection(ag *v1alpha1.Agent) (*grpc.ClientConn, error) {
	clientCert, err := store.KVStore.Get(
		context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.SecretKeyPrefix, ag.Spec.ClientCertSecret))
	if err != nil {
		return nil, err
	}

	clientKey, err := store.KVStore.Get(
		context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.SecretKeyPrefix, ag.Spec.ClientKeySecret))
	if err != nil {
		return nil, err
	}

	rootCA, err := store.KVStore.Get(
		context.TODO(),
		fmt.Sprintf("%s/%s", v1alpha1.SecretKeyPrefix, ag.Spec.RootCASecret))
	if err != nil {
		return nil, err
	}

	var certSecret, keySecret, rootCASecret v1alpha1.Secret
	err = json.Unmarshal(clientCert.Data, &certSecret)
	if err != nil {
		return nil, fmt.Errorf("errow while unmarshling client secret: %s", err)
	}

	err = json.Unmarshal(clientKey.Data, &keySecret)
	if err != nil {
		return nil, fmt.Errorf("errow while unmarshling client key: %s", err)
	}

	err = json.Unmarshal(rootCA.Data, &rootCASecret)
	if err != nil {
		return nil, fmt.Errorf("errow while unmarshling root CA certificate: %s", err)
	}

	var cert, key, root []byte
	cert, err = certSecret.GetContent()
	if err != nil {
		return nil, fmt.Errorf("error while decoding cert secret: %s", err)
	}

	key, err = keySecret.GetContent()
	if err != nil {
		return nil, fmt.Errorf("error while decoding key secret: %s", err)
	}

	root, err = rootCASecret.GetContent()
	if err != nil {
		return nil, fmt.Errorf("error while decoding root CA secret: %s", err)
	}

	certificate, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(root)
	if !ok {
		return nil, fmt.Errorf("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   "agent.xene.io",
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	dialOption := grpc.WithTransportCredentials(transportCreds)

	u, err := url.Parse(ag.Spec.Address)
	if err != nil {
		return nil, fmt.Errorf("error while parsing agent address: %s", err)
	}
	return grpc.Dial(u.Host, dialOption)
}

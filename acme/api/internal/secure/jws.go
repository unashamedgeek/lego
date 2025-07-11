package secure

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"fmt"

	"github.com/go-acme/lego/v4/acme/api/internal/nonces"
	jose "github.com/go-jose/go-jose/v4"
)

// JWS Represents a JWS.
type JWS struct {
	privKey crypto.PrivateKey
	kid     string // Key identifier
	nonces  *nonces.Manager
}

// NewJWS Create a new JWS.
func NewJWS(privateKey crypto.PrivateKey, kid string, nonceManager *nonces.Manager) *JWS {
	return &JWS{
		privKey: privateKey,
		nonces:  nonceManager,
		kid:     kid,
	}
}

// SetKid Sets a key identifier.
func (j *JWS) SetKid(kid string) {
	j.kid = kid
}

// SignContent Signs a content with the JWS.
func (j *JWS) SignContent(url string, content []byte) (*jose.JSONWebSignature, error) {
	var alg jose.SignatureAlgorithm
	switch k := j.privKey.(type) {
	case *rsa.PrivateKey:
		alg = jose.RS256
	case *ecdsa.PrivateKey:
		if k.Curve == elliptic.P256() {
			alg = jose.ES256
		} else if k.Curve == elliptic.P384() {
			alg = jose.ES384
		}
	}

	signKey := jose.SigningKey{
		Algorithm: alg,
		Key:       jose.JSONWebKey{Key: j.privKey, KeyID: j.kid},
	}

	options := jose.SignerOptions{
		NonceSource: j.nonces,
		ExtraHeaders: map[jose.HeaderKey]any{
			"url": url,
		},
	}

	if j.kid == "" {
		options.EmbedJWK = true
	}

	signer, err := jose.NewSigner(signKey, &options)
	if err != nil {
		return nil, fmt.Errorf("failed to create jose signer: %w", err)
	}

	signed, err := signer.Sign(content)
	if err != nil {
		return nil, fmt.Errorf("failed to sign content: %w", err)
	}
	return signed, nil
}

// SignEABContent Signs an external account binding content with the JWS.
func (j *JWS) SignEABContent(url, kid string, hmac []byte) (*jose.JSONWebSignature, error) {
	jwk := jose.JSONWebKey{Key: j.privKey}
	jwkJSON, err := jwk.Public().MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("acme: error encoding eab jwk key: %w", err)
	}

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: hmac},
		&jose.SignerOptions{
			EmbedJWK: false,
			ExtraHeaders: map[jose.HeaderKey]any{
				"kid": kid,
				"url": url,
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create External Account Binding jose signer: %w", err)
	}

	signed, err := signer.Sign(jwkJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to External Account Binding sign content: %w", err)
	}

	return signed, nil
}

// GetKeyAuthorization Gets the key authorization for a token.
func (j *JWS) GetKeyAuthorization(token string) (string, error) {
	var publicKey crypto.PublicKey
	switch k := j.privKey.(type) {
	case *ecdsa.PrivateKey:
		publicKey = k.Public()
	case *rsa.PrivateKey:
		publicKey = k.Public()
	}

	// Generate the Key Authorization for the challenge
	jwk := &jose.JSONWebKey{Key: publicKey}

	thumbBytes, err := jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", err
	}

	// unpad the base64URL
	keyThumb := base64.RawURLEncoding.EncodeToString(thumbBytes)

	return token + "." + keyThumb, nil
}

package gpg_key_verifier

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func (gv gpgKeyVerifier) ValidateSignature(data []byte, signature []byte) error {
	plainMessage := crypto.NewPlainMessage(data)
	pgpSignature := crypto.NewPGPSignature(signature)

	if err := gv.keyring.VerifyDetached(plainMessage, pgpSignature, crypto.GetUnixTime()); err != nil {
		return err
	}

	return nil
}

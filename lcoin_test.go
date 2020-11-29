package lcoin

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	u1 User
	u2 User
	u3 User
)

func init() {
	var err error
	u1, err = NewUserFromKey("User1", base64DecodeOrPanic("MHcCAQEEIMbqqXWF57EF7Av0kc9EN5lFXWuZSE/n9eEurq2+F6OXoAoGCCqGSM49AwEHoUQDQgAE1ziOUG1YSzMpRhh2Ez59aVevVsNb2Fv35fb6ieukftq1yZIB95sZhxBD6V31i/8HpZwozBHQs7bAR881b4LKYg=="))
	if err != nil {
		panic(err)
	}

	u2, err = NewUserFromKey("User2", base64DecodeOrPanic("MHcCAQEEIG2RJL2RIJFIhHZ/9wkM5kuWyAih4TEKlSNrooxcU7e9oAoGCCqGSM49AwEHoUQDQgAEY3bd1AyKJFTfxQM3b2LK+KtL9lSXPlgIZsb5vBIr44Z5/RWGVO/JI623NZtRm04wyz1PkSIGII77kfa0ulF9xQ=="))
	if err != nil {
		panic(err)
	}

	u3, err = NewUserFromKey("User3", base64DecodeOrPanic("MHcCAQEEIHetwtaW1YHHFOV3Ztyv6cRr6uAvO0pMT6MOPPRt9d7zoAoGCCqGSM49AwEHoUQDQgAEUfoSxcDlKYAqwWvqpdjX1KIDa6w6J8pGQ8ksvBVLomGkkeG2cLhDsnyxJ3B8EtX0Pp8szbzatLfrQ3pOmAVJEA=="))
	if err != nil {
		panic(err)
	}
}

func TestMessage_Unlock(t *testing.T) {
	tests := []struct {
		name        string
		message     Message
		claimer     User
		claimPubkey ecdsa.PublicKey
		expected    bool
	}{
		{
			name: "User1 sends to User2. User 2 claims sucessfully.",
			message: Message{
				SenderAddress:   u1.Address(),
				ReceiverAddress: u2.Address(),
				Payload:         "This is a massage for user 2",
			},
			claimer:     u2,
			claimPubkey: u2.PrivateKey.PublicKey,
			expected:    true,
		},
		{
			name: "User1 sends to User2. User 3 claims unsucessfully with own public key.",
			message: Message{
				SenderAddress:   u1.Address(),
				ReceiverAddress: u2.Address(),
				Payload:         "This is a massage for user 2",
			},
			claimer:     u3,
			claimPubkey: u3.PrivateKey.PublicKey,
			expected:    false,
		},
		{
			name: "User1 sends to User2. User 3 claims unsucessfully with User 2 public key.",
			message: Message{
				SenderAddress:   u1.Address(),
				ReceiverAddress: u2.Address(),
				Payload:         "This is a massage for user 2",
			},
			claimer:     u3,
			claimPubkey: u2.PrivateKey.PublicKey,
			expected:    false,
		},
	}

	for _, test := range tests {
		messageHash := test.message.Hash()
		claimerSignature, err := ecdsa.SignASN1(rand.Reader, test.claimer.PrivateKey, messageHash[:])
		assert.NoError(t, err)
		assert.Equal(t, test.expected, test.message.Claim(claimerSignature, &test.claimer.PrivateKey.PublicKey))
	}
}

func base64DecodeOrPanic(s string) []byte {
	result, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return result
}

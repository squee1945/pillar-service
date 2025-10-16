package runner

import (
	"context"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
)

func kmsEncrypt(ctx context.Context, keyName string, plaintext []byte) ([]byte, error) {
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating kms client: %v", err)
	}
	defer client.Close()

	req := &kmspb.EncryptRequest{
		Name:      keyName,
		Plaintext: plaintext,
	}

	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("encrypting data: %v", err)
	}

	return resp.Ciphertext, nil
}

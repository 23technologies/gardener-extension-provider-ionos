/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package ionos provides types and functions used for ionos interaction
package ionos

import (
	"context"
	"fmt"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CredentialsStruct struct {
	User     string
	Password string
}

// Credentials contains the necessary ionos credential information.
type Credentials struct {
	ionos    *CredentialsStruct
	ionosMCM *CredentialsStruct
	ionosCCM *CredentialsStruct
	ionosCSI *CredentialsStruct
}

// MCM returns the token used for the Machine Controller Manager.
func (c *Credentials) IonosMCM() CredentialsStruct {
	if c.ionosMCM != nil {
		return *c.ionosMCM
	}
	return *c.ionos
}

// CCM returns the token used for the Cloud Controller Manager.
func (c *Credentials) IonosCCM() CredentialsStruct {
	if c.ionosCCM != nil {
		return *c.ionosCCM
	}
	return *c.ionos
}

// CSI returns the token used for the Container Storage Interface driver.
func (c *Credentials) IonosCSI() CredentialsStruct {
	if c.ionosCSI != nil {
		return *c.ionosCSI
	}
	return *c.ionos
}

// GetCredentials computes for a given context and infrastructure the corresponding credentials object.
//
// PARAMETERS
// ctx       context.Context        Execution context
// c         client.Client          Controller client
// secretRef corev1.SecretReference Secret reference to read credentials from
func GetCredentials(ctx context.Context, c client.Client, secretRef corev1.SecretReference) (*Credentials, error) {
	secret, err := extensionscontroller.GetSecretByReference(ctx, c, &secretRef)
	if err != nil {
		return nil, err
	}
	return ExtractCredentials(secret)
}

func extractUserPass(secret *corev1.Secret, usernameKey, passwordKey string) (*CredentialsStruct, error) {
	user, ok := secret.Data[usernameKey]
	if !ok {
		return nil, fmt.Errorf("missing %q field in secret", usernameKey)
	}

	password, ok := secret.Data[passwordKey]
	if !ok {
		return nil, fmt.Errorf("missing %q field in secret", passwordKey)
	}

	return &CredentialsStruct{User: string(user), Password: string(password)}, nil
}

// ExtractCredentials generates a credentials object for a given provider secret.
//
// PARAMETERS
// secret   *corev1.Secret Secret to extract tokens from
func ExtractCredentials(secret *corev1.Secret) (*Credentials, error) {
	if secret.Data == nil {
		return nil, fmt.Errorf("secret does not contain any data")
	}

	ionos, ionosErr := extractUserPass(secret, IonosCredentialsUserKey, IonosCredentialsPasswordKey)

	mcm, err := extractUserPass(secret, IonosMCMCredentialsUserKey, IonosMCMCredentialsPasswordKey)
	if err != nil && ionosErr != nil {
		return nil, fmt.Errorf("Need either common or machine controller manager specific ionos account credentials: %s, %s", ionosErr, err)
	}
	ccm, err := extractUserPass(secret, IonosCCMCredentialsUserKey, IonosCCMCredentialsPasswordKey)
	if err != nil && ionosErr != nil {
		return nil, fmt.Errorf("Need either common or cloud controller manager specific ionos account credentials: %s, %s", ionosErr, err)
	}
	csi, err := extractUserPass(secret, IonosCSICredentialsUserKey, IonosCSICredentialsPasswordKey)
	if err != nil && ionosErr != nil {
		return nil, fmt.Errorf("Need either common or cloud controller manager specific ionos account credentials: %s, %s", ionosErr, err)
	}

	return &Credentials{
		ionos:    ionos,
		ionosMCM: mcm,
		ionosCCM: ccm,
		ionosCSI: csi,
	}, nil
}

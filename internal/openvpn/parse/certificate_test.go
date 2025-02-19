package parse

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ExtractCert(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		b        []byte
		certData string
		err      error
	}{
		"no input": {
			err: errors.New("cannot extract PEM data: cannot decode PEM encoded block"),
		},
		"bad input": {
			b:   []byte{1, 2, 3},
			err: errors.New("cannot extract PEM data: cannot decode PEM encoded block"),
		},
		"valid key": {
			b:        []byte(validCertPEM),
			certData: validCertData,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			certData, err := ExtractCert(testCase.b)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.certData, certData)
		})
	}
}

const validCertPEM = `
-----BEGIN CERTIFICATE-----
MIIGrDCCBJSgAwIBAgIEAdTnfTANBgkqhkiG9w0BAQsFADB7MQswCQYDVQQGEwJS
TzESMBAGA1UEBxMJQnVjaGFyZXN0MRgwFgYDVQQKEw9DeWJlckdob3N0IFMuQS4x
GzAZBgNVBAMTEkN5YmVyR2hvc3QgUm9vdCBDQTEhMB8GCSqGSIb3DQEJARYSaW5m
b0BjeWJlcmdob3N0LnJvMB4XDTIwMDcwNDE1MjkzNloXDTMwMDcwMjE1MjkzNlow
fTELMAkGA1UEBhMCUk8xEjAQBgNVBAcMCUJ1Y2hhcmVzdDEYMBYGA1UECgwPQ3li
ZXJHaG9zdCBTLkEuMR0wGwYDVQQDDBRjLmoua2xhdmVyQGdtYWlsLmNvbTEhMB8G
CSqGSIb3DQEJARYSaW5mb0BjeWJlcmdob3N0LnJvMIICIjANBgkqhkiG9w0BAQEF
AAOCAg8AMIICCgKCAgEAobp2NlGUHMNBe08YEOnVG3QJjF3ZaXbRhE/II9rmtgJT
NZtDohGChvFlNRsExKzVrKxHCeuJkVffwzQ6fYk4/M1RdYLJUh0UVw3e4WdApw8E
7TJZxDYm4SHQNXUvt1Rt5TjslcXxIpDZgrMSc/kHROYEL9tdgdzPZErUJehXyJPh
EzIrzmAJh501x7WwKPz9ctSVlItyavqEWFF2vyUa6X9DYmD9mQTz5c+VXNO5DkXm
PFBIaEVDnvFtcjGJ56yEvFnWVukL+OUX7ezowrIOFOcp9udjgpeiHq+XvsQ6ER0D
Jt25MiEId3NjkxtZ8BitDftTcLN/kt81hWKT7adMVc3kpIZ80cxrwRCttMd7sHAz
KI9u7pMxv10eUOsIEY87ewBe3l6KvEnjA+9uIjim6gLLebDIaEH50Ee9PzNJ8fqQ
2u54Ab4bt00/H1sUnJ6Ss/+WsQDOK1BsPRKKcnHZntOlHrs2Tu5+txKNU2cOapI8
SjVULUNKrRXASbpfWnLUfri/HO742bJb/TjkOJcOxta3hTPFAhaRWBusVlB41XVH
euH5DAhugYXeSNK6/6Ul8YvKUNH/7QbxuGIGXfth19Xl4QLI1umyEjZopSlt3tOi
O2V1soVNSQCCfxXVoCTMESMLjhkjWdmBDhdy2GTW7S4YoJfqVKiS18rYkN7I4ZMC
AwEAAaOCATQwggEwMAkGA1UdEwQCMAAwCwYDVR0PBAQDAgWgMDQGCWCGSAGG+EIB
DQQnFiVDeWJlckdob3N0IEdlbmVyYXRlZCBVc2VyIENlcnRpZmljYXRlMBEGCWCG
SAGG+EIBAQQEAwIHgDAdBgNVHQ4EFgQULwUtU5s6pL2NN9gPeEnKX0dhwiswga0G
A1UdIwSBpTCBooAU6tdK1g/He5qzjeAoM5eHt4in9iWhf6R9MHsxCzAJBgNVBAYT
AlJPMRIwEAYDVQQHEwlCdWNoYXJlc3QxGDAWBgNVBAoTD0N5YmVyR2hvc3QgUy5B
LjEbMBkGA1UEAxMSQ3liZXJHaG9zdCBSb290IENBMSEwHwYJKoZIhvcNAQkBFhJp
bmZvQGN5YmVyZ2hvc3Qucm+CCQCcVButZsQ0uzANBgkqhkiG9w0BAQsFAAOCAgEA
ystGIMYhQWaEdTqlnLCytrr8657t+PuidZMNNIaPB3wN2Fi2xKf14DTg03mqxjmP
Pb+f+PVNIOV5PdWD4jcQwOP1GEboGV0DFzlRGeAtDcvKwdee4oASJbZq1CETqDao
hQTxKEWC+UBk2F36nOaEI6Sab+Mb4cR9//PAwvzOqrXuGF5NuIOX7eFtCMQSgQq6
lRRqTQjekm0Dxigx4JA92Jo2qZRwCJ0T3IXBJGL831HCFJbDWv8PV3lsfFb/i2+v
r54uywFQVWWp18dYi97gipfuQ4zRg2Ldx5aXSmnhhKpg5ioZvtk043QofF12YORh
obElqavRbvvhZvlCouvcuoq9QKi7IPe5SJZkZ1X7ezMesCwBzwFpt6vRUAcslsNF
bcYS1iSENlY/PTcDqBhbKuc9yAhq+/aUgaY/8VF5RWVzSRZufbf3BPwOkE4K0Uyb
aobO/YX0JOkCacAD+4tdR6YSXNIMMRAOCBQvxbxFXaHzhwhzBAjdsC56FrJKwXvQ
rRLU3tF4P0zFMeNTay8uTtUXugDK7EnklLESuYdpUJ8bUMlAUhJBi6UFI9/icMud
xXvLRvhnBW9EtKib5JnVFUovcEUt+3EJbyst05nkL4YPjQS4TC9DHdo5SyRAy1Tp
iOCYTbretAFZRhh6ycUN5hBeN8GMQxiMreMtDV4PEIQ=
-----END CERTIFICATE-----
`
const validCertData = "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" //nolint:lll

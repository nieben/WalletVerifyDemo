package handler

import (
	"WalletVerifyDemo/tools"
	"WalletVerifyDemo/types"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"net/http/httptest"
	"testing"
)

var (
	e                                        *echo.Echo
	validVerifyRequest, invalidVerifyRequest *types.VerifyRequest
)

func init() {
	e = echo.New()
	e.Use(middleware.Logger())
	e.Validator = tools.NewCustomerValidator()

	buf := make([]byte, 32)
	io.ReadFull(rand.Reader, buf)

	seckey, _ := crypto.LoadECDSA("samplePrvKey")

	address := crypto.PubkeyToAddress(seckey.PublicKey).Hex()

	sig, _ := secp256k1.Sign(buf, crypto.FromECDSA(seckey))

	validVerifyRequest = &types.VerifyRequest{
		Address:       address,
		Message:       base64.StdEncoding.EncodeToString(buf),
		SignedMessage: base64.StdEncoding.EncodeToString(sig),
	}

	invalidVerifyRequest = &types.VerifyRequest{
		Address:       "0x0000000000000000000000000000000000000000", // any other address
		Message:       base64.StdEncoding.EncodeToString(buf),
		SignedMessage: base64.StdEncoding.EncodeToString(sig),
	}
}

func TestVerifyHandler_Message(t *testing.T) {
	req := httptest.NewRequest(echo.GET, "/get_message", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"default",
			args{
				c,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := VerifyHandler{}
			if err := h.Message(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("VerifyHandler.Message error: %v, wantErr: %v", err, tt.wantErr)
			} else {
				b, _ := io.ReadAll(rec.Body)
				t.Logf("VerifyHandler.Message success: %s", string(b))
			}
		})
	}
}

func TestVerifyHandler_Verify(t *testing.T) {
	b, _ := json.Marshal(validVerifyRequest)
	b2, _ := json.Marshal(invalidVerifyRequest)
	req := httptest.NewRequest(echo.POST, "/verify", bytes.NewBuffer(b))
	req2 := httptest.NewRequest(echo.POST, "/verify", bytes.NewBuffer(b2))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c2 := e.NewContext(req2, rec)
	req.Header.Add("Content-Type", "application/json")
	req2.Header.Add("Content-Type", "application/json")

	type args struct {
		c echo.Context
	}
	tests := []struct {
		name             string
		args             args
		expectedVerified bool
	}{
		{
			"validSign",
			args{
				c,
			},
			true,
		},
		{
			"invalidSign",
			args{
				c2,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := VerifyHandler{}
			if err := h.Verify(tt.args.c); err != nil {
				t.Errorf("VerifyHandler.Verify() error: %v, expectedResult: %v", err, tt.expectedVerified)
			} else {
				b, _ := io.ReadAll(rec.Body)
				resp := &types.VerifyResponse{}
				json.Unmarshal(b, resp)
				if resp.Verified != tt.expectedVerified {
					t.Errorf("VerifyHandler.Verify() failed, got: %v, expectedVerified: %v",
						resp.Verified, tt.expectedVerified)
				} else {
					t.Logf("VerifyHandler.Verify() success: %s, expectedVerified: %v", string(b), tt.expectedVerified)
				}
			}
		})
	}
}

package jwt_provider

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/forta-network/forta-core-go/security"
	"github.com/golang-jwt/jwt/v4"
)

func Test_createBotJWT(t *testing.T) {
	dir := t.TempDir()
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := ks.NewAccount("Forta123")
	if err != nil {
		t.Fatal(err)
	}

	key, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		key   *keystore.Key
		data  map[string]interface{}
		botID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "can create valid jwt",
			args: args{
				key:   key,
				botID: "0xbbb",
				data: map[string]interface{}{
					"hash": requestHash("/home", nil).String(),
				},
			},
			want:    "{}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := CreateBotJWT(tt.args.key, tt.args.botID, tt.args.data)
				if (err != nil) != tt.wantErr {
					t.Errorf("createBotJWT() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				token, err := security.VerifyScannerJWT(got)
				if err != nil {
					t.Fatal(err)
				}

				claims, ok := token.Token.Claims.(jwt.MapClaims)
				if !ok {
					t.Fatal("invalid jwt claims")
				}

				if hash := claims["hash"]; hash != tt.args.data["hash"] {
					t.Errorf("createBotJWT() got hash = %v, want hash %v", hash, tt.args.data["hash"])
				}

				if botID := claims["bot-id"]; botID != tt.args.botID {
					t.Errorf("createBotJWT() got bot id = %v, want bot id %v", botID, tt.args.botID)
				}
			},
		)
	}
}

func TestJWTProvider_Start(t *testing.T) {
	t.SkipNow()

	dir := t.TempDir()
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)

	_, err := ks.NewAccount("Forta123")
	if err != nil {
		t.Fatal(err)
	}

	key, err := security.LoadKeyWithPassphrase(dir, "Forta123")
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		cfg JWTProviderConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "spawn jwt provider service for 10 seconds",
			fields: fields{
				cfg: JWTProviderConfig{Key: key},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				j, err := initProvider(&tt.fields.cfg)
				if err != nil {
					t.Fatal(err)
				}

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*155)
				defer cancel()

				if err := j.StartWithContext(ctx); (err != nil) != tt.wantErr {
					t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
				}

				<-ctx.Done()
			},
		)
	}
}

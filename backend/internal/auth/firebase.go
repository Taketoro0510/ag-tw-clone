package auth

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type FirebaseAuth interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type firebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(ctx context.Context, projectID string) (FirebaseAuth, error) {
	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		return nil, err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return &firebaseAuth{client: client}, nil
}

func (a *firebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return a.client.VerifyIDToken(ctx, idToken)
}

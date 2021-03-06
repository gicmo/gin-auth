// Copyright (c) 2016, German Neuroinformatics Node (G-Node),
//                     Adrian Stoewer <adrian.stoewer@rz.ifi.lmu.de>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License. See
// LICENSE file in the root of the Project.

package data

import (
	"database/sql"
	"time"

	"github.com/G-Node/gin-auth/conf"
	"github.com/G-Node/gin-auth/util"
)

// AccessToken represents an OAuth access token
type AccessToken struct {
	Token       string // This is just a random string not the JWT token
	Scope       util.StringSet
	Expires     time.Time
	ClientUUID  string
	AccountUUID sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ListAccessTokens returns all access tokens sorted by creation time.
func ListAccessTokens() []AccessToken {
	const q = `SELECT * FROM AccessTokens WHERE expires > now() ORDER BY createdAt`

	accessTokens := make([]AccessToken, 0)
	err := database.Select(&accessTokens, q)
	if err != nil {
		panic(err)
	}

	return accessTokens
}

// GetAccessToken returns a access token with a given token.
// Returns false if no such access token exists.
func GetAccessToken(token string) (*AccessToken, bool) {
	const q = `SELECT * FROM AccessTokens WHERE token=$1 AND expires > now()`

	accessToken := &AccessToken{}
	err := database.Get(accessToken, q, token)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	return accessToken, err == nil
}

// Create stores a new access token in the database.
// If the token is empty a random token will be generated.
func (tok *AccessToken) Create() error {
	const q = `INSERT INTO AccessTokens (token, scope, expires, clientUUID, accountUUID, createdAt, updatedAt)
	           VALUES ($1, $2, $3, $4, $5, now(), now())
	           RETURNING *`

	tok.Expires = time.Now().Add(conf.GetServerConfig().TokenLifeTime)
	if tok.Token == "" {
		tok.Token = util.RandomToken()
	}

	return database.Get(tok, q, tok.Token, tok.Scope, tok.Expires, tok.ClientUUID, tok.AccountUUID)
}

// UpdateExpirationTime updates the expiration time and stores
// the new time in the database.
func (tok *AccessToken) UpdateExpirationTime() error {
	const q = `UPDATE AccessTokens SET (expires, updatedAt) = ($1, now())
	           WHERE token=$2
	           RETURNING *`

	return database.Get(tok, q, time.Now().Add(conf.GetServerConfig().TokenLifeTime), tok.Token)
}

// Delete removes an access token from the database.
func (tok *AccessToken) Delete() error {
	const q = `DELETE FROM AccessTokens WHERE token=$1`

	_, err := database.Exec(q, tok.Token)
	return err
}

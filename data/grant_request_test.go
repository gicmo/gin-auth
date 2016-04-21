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
	"github.com/G-Node/gin-auth/util"
	"testing"
)

const (
	grantReqTokenAlice = "U7JIKKYI"
	grantReqCodeAlice  = "HGZQP6WE"
	grantReqTokenBob   = "B4LIMIMB"
)

func TestListGrantRequests(t *testing.T) {
	defer util.FailOnPanic(t)
	InitTestDb(t)

	requests := ListGrantRequests()
	if len(requests) != 2 {
		t.Error("Exactly two grant requests expected in list")
	}
}

func TestGetGrantRequest(t *testing.T) {
	defer util.FailOnPanic(t)
	InitTestDb(t)

	req, ok := GetGrantRequest(grantReqTokenAlice)
	if !ok {
		t.Error("Grant request does not exist")
	}
	if req.ScopeRequested[0] != "repo-read" {
		t.Errorf("First requested scope was expected to be 'repo-read'")
	}

	_, ok = GetGrantRequest("doesNotExist")
	if ok {
		t.Error("Grant request should not exist")
	}
}

func TestGetGrantRequestByCode(t *testing.T) {
	defer util.FailOnPanic(t)
	InitTestDb(t)

	req, ok := GetGrantRequestByCode(grantReqCodeAlice)
	if !ok {
		t.Error("Grant request does not exist")
	}
	if req.ScopeRequested[0] != "repo-read" {
		t.Errorf("First requested scope was expected to be 'repo-read'")
	}

	_, ok = GetGrantRequestByCode("doesNotExist")
	if ok {
		t.Error("Grant request should not exist")
	}
}

func TestGrantRequestExchangeCodeForTokens(t *testing.T) {
	defer util.FailOnPanic(t)
	InitTestDb(t)

	req, ok := GetGrantRequest(grantReqTokenAlice)
	if !ok {
		t.Error("Grant request does not exist")
	}

	accessToken, refreshToken, err := req.ExchangeCodeForTokens()
	if err != nil {
		t.Error(err)
	}

	access, ok := GetAccessToken(accessToken)
	if !ok {
		t.Error("Unable to find created access token")
	}
	if access.AccountUUID != uuidAlice {
		t.Error("Access token has a wrong account UUID")
	}

	refresh, ok := GetRefreshToken(refreshToken)
	if !ok {
		t.Error("Unable to find created refresh token")
	}
	if refresh.AccountUUID != uuidAlice {
		t.Error("Refresh token has a wrong account UUID")
	}

	_, ok = GetGrantRequest(grantReqTokenAlice)
	if ok {
		t.Error("Grant request was expected to be deleted")
	}
}

func TestCreateGrantRequest(t *testing.T) {
	InitTestDb(t)

	token := util.RandomToken()
	state := util.RandomToken()
	code := util.RandomToken()
	fresh := GrantRequest{
		Token:          token,
		GrantType:      "code",
		State:          state,
		Code:           sql.NullString{String: code, Valid: true},
		ScopeRequested: util.SqlStringSlice{"foo-read", "foo-write", "foo-admin"},
		ScopeApproved:  util.SqlStringSlice{"foo-read"},
		ClientUUID:     uuidClientGin,
		AccountUUID:    sql.NullString{String: uuidAlice, Valid: true}}

	err := fresh.Create()
	if err != nil {
		t.Error(err)
	}

	check, ok := GetGrantRequest(token)
	if !ok {
		t.Error("Grant request does not exist")
	}
	if check.State != state {
		t.Error("State does not match")
	}
	if check.Code.Valid && check.Code.String != code {
		t.Error("Code does not match")
	}
	if check.ScopeRequested[0] != "foo-read" {
		t.Error("First requested scope was expected to be 'foo-read'")
	}
	if check.ScopeApproved[0] != "foo-read" {
		t.Error("First approved scope was expected to be 'foo-read'")
	}
}

func TestUpdateGrantRequest(t *testing.T) {
	InitTestDb(t)

	newCode := util.RandomToken()
	newState := util.RandomToken()

	req, ok := GetGrantRequest(grantReqTokenAlice)
	if !ok {
		t.Error("Grant request does not exist")
	}

	req.Code = sql.NullString{String: newCode, Valid: true}
	req.State = newState

	err := req.Update()
	if err != nil {
		t.Error(err)
	}

	check, ok := GetGrantRequest(grantReqTokenAlice)
	if !ok {
		t.Error("Grant request does not exist")
	}
	if check.Code.Valid && check.Code.String != newCode {
		t.Error("Code does not match")
	}
	if check.State != newState {
		t.Error("State does not match")
	}
}

func TestDeleteGrantRequest(t *testing.T) {
	InitTestDb(t)

	req, ok := GetGrantRequest(grantReqTokenAlice)
	if !ok {
		t.Error("Grant request does not exist")
	}

	err := req.Delete()
	if err != nil {
		t.Error(err)
	}

	_, ok = GetGrantRequest(uuidClientGin)
	if ok {
		t.Error("Grant request should not exist")
	}
}

func TestGrantRequestGetApproval(t *testing.T) {
	InitTestDb(t)

	req, ok := GetGrantRequest(grantReqTokenAlice)
	if !ok {
		t.Error("Grant request does not exist")
	}
	cli, ok := req.GetClientApproval()
	if !ok {
		t.Error("Approval does not exist")
	}
	if cli.ClientUUID != req.ClientUUID {
		t.Errorf("Client UUID should be '%s' but was '%s'", req.ClientUUID, cli.ClientUUID)
	}
	if cli.AccountUUID != uuidAlice {
		t.Errorf("Account UUID should be '%s' but was '%s'", uuidAlice, cli.AccountUUID)
	}

	req, ok = GetGrantRequest(grantReqTokenBob)
	if !ok {
		t.Error("Grant request does not exist")
	}
	cli, ok = req.GetClientApproval()
	if ok {
		t.Error("Approval should not exist")
	}
}

func TestGrantRequestApproveScopes(t *testing.T) {
	InitTestDb(t)

	token := util.RandomToken()
	request := GrantRequest{
		Token:          token,
		GrantType:      "code",
		State:          util.RandomToken(),
		ScopeRequested: util.SqlStringSlice{"repo-read"},
		ClientUUID:     uuidClientGin,
		AccountUUID:    sql.NullString{String: uuidAlice, Valid: true}}

	err := request.Create()
	if err != nil {
		t.Error(err)
	}

	ok := request.ApproveScopes()
	if !ok {
		t.Error("Scopes not approved")
	}

	check, ok := GetGrantRequest(token)
	if !ok {
		t.Error("Grant request does not exist")
	}
	if len(check.ScopeApproved) != 1 {
		t.Error("Approved scope should have length 1")
	}
	if check.ScopeApproved[0] != "repo-read" {
		t.Errorf("First element in approved scope shoule be 'repo-read' but was '%s'", check.ScopeApproved[0])
	}
}

func TestGrantRequestIsApproved(t *testing.T) {
	InitTestDb(t)

	request := GrantRequest{
		ScopeRequested: util.SqlStringSlice{"repo-read"},
		ScopeApproved:  util.SqlStringSlice{"repo-read", "something-else"},
	}
	ok := request.IsApproved()
	if !ok {
		t.Error("Grant request should be approved")
	}

	request = GrantRequest{
		ScopeRequested: util.SqlStringSlice{"repo-read", "repo-write"},
		ScopeApproved:  util.SqlStringSlice{"repo-read", "repo-write"},
	}
	ok = request.IsApproved()
	if !ok {
		t.Error("Grant request should be approved")
	}

	request = GrantRequest{
		ScopeRequested: util.SqlStringSlice{"repo-read", "repo-write", "something-else"},
		ScopeApproved:  util.SqlStringSlice{"repo-read", "repo-write"},
	}
	ok = request.IsApproved()
	if ok {
		t.Error("Grant request should not be approved")
	}

	request = GrantRequest{
		ScopeRequested: util.SqlStringSlice{},
		ScopeApproved:  util.SqlStringSlice{},
	}
	ok = request.IsApproved()
	if ok {
		t.Error("Grant request should not be approved")
	}
}

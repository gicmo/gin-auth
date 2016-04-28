// Copyright (c) 2016, German Neuroinformatics Node (G-Node),
//                     Adrian Stoewer <adrian.stoewer@rz.ifi.lmu.de>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License. See
// LICENSE file in the root of the Project.

package data

import (
	"github.com/G-Node/gin-auth/util"
	"testing"
	"time"
)

const (
	sessionTokenAlice = "DNM5RS3C"
	sessionTokenBob   = "2MFZZUKI"
)

func TestListSessions(t *testing.T) {
	defer util.FailOnPanic(t)
	InitTestDb(t)

	sessions := ListSessions()
	if len(sessions) != 2 {
		t.Error("Exactly to sessions expected in slice")
	}
}

func TestGetSession(t *testing.T) {
	defer util.FailOnPanic(t)
	InitTestDb(t)

	sess, ok := GetSession(sessionTokenAlice)
	if !ok {
		t.Error("Session does not exist")
	}
	if sess.AccountUUID != uuidAlice {
		t.Errorf("AccountUUID was expectd to be '%s'", uuidAlice)
	}

	_, ok = GetSession("doesNotExist")
	if ok {
		t.Error("Session should not exist")
	}
}

func TestClearOldSessions(t *testing.T) {
	defer util.FailOnPanic(t)
	InitTestDb(t)

	deleted := ClearOldSessions()
	if deleted != 1 {
		t.Error("Exactly one session is supposed to be deleted")
	}

	_, ok := GetSession(sessionTokenBob)
	if ok {
		t.Error("Bobs session should not exist")
	}
}

func TestCreateSession(t *testing.T) {
	InitTestDb(t)

	token := util.RandomToken()
	fresh := Session{
		Token:       token,
		Expires:     time.Now().Add(time.Hour * 12),
		AccountUUID: uuidAlice}

	err := fresh.Create()
	if err != nil {
		t.Error(err)
	}

	check, ok := GetSession(token)
	if !ok {
		t.Error("Token does not exist")
	}
	if check.AccountUUID != uuidAlice {
		t.Errorf("AccountUUID is supposed to be '%s'", uuidAlice)
	}
}

func TestSessionUpdateExpirationTime(t *testing.T) {
	InitTestDb(t)

	sess, ok := GetSession(sessionTokenBob)
	if !ok {
		t.Error("Session does not exist")
	}
	if time.Since(sess.Expires) < 0 {
		t.Error("Session should be expired")
	}

	sess.UpdateExpirationTime()
	if time.Since(sess.Expires) > 0 {
		t.Error("Session should not be expired")
	}

	check, ok := GetSession(sessionTokenBob)
	if !ok {
		t.Error("Session does not exist")
	}
	if time.Since(check.Expires) > 0 {
		t.Error("Session should not be expired")
	}
}

func TestSessionDelete(t *testing.T) {
	InitTestDb(t)

	sess, ok := GetSession(sessionTokenAlice)
	if !ok {
		t.Error("Session does not exist")
	}

	err := sess.Delete()
	if err != nil {
		t.Error(err)
	}

	_, ok = GetSession(sessionTokenAlice)
	if ok {
		t.Error("Session should not exist")
	}
}
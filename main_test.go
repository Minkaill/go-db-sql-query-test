package main

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	clientID := 1

	cl, err := selectClient(db, clientID)

	require.NoError(t, err)
	require.Equal(t, clientID, cl.ID)
	assert.NotEmpty(t, cl.FIO)
	assert.NotEmpty(t, cl.Login)
	assert.NotEmpty(t, cl.Birthday)
	assert.NotEmpty(t, cl.Email)
}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	clientID := -1

	cl, err := selectClient(db, clientID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, cl.FIO)
	assert.Empty(t, cl.Login)
	assert.Empty(t, cl.Birthday)
	assert.Empty(t, cl.Email)
}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	rowID, err := insertClient(db, cl)
	require.NoError(t, err)
	require.NotEmpty(t, rowID)

	cl.ID = rowID

	client, err := selectClient(db, rowID)
	require.NoError(t, err)

	assert.Equal(t, cl.ID, client.ID)
	assert.Equal(t, cl.FIO, client.FIO)
	assert.Equal(t, cl.Login, client.Login)
	assert.Equal(t, cl.Birthday, client.Birthday)
	assert.Equal(t, cl.Email, client.Email)
}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	rowID, err := insertClient(db, cl)
	require.NotEmpty(t, rowID)
	require.NoError(t, err)

	client, err := selectClient(db, rowID)
	require.NoError(t, err)
	require.Equal(t, rowID, client.ID)

	err = deleteClient(db, rowID)
	require.NoError(t, err)

	_, err = selectClient(db, client.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

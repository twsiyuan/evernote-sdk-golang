package utils

import (
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/twsiyuan/evernote-sdk-golang/notestore"
	"github.com/twsiyuan/evernote-sdk-golang/userstore"
)

type EnvironmentType int

const (
	SANDBOX EnvironmentType = iota
	PRODUCTION
)

func GetHost(envType EnvironmentType) string {
	host := "www.evernote.com"
	if envType == SANDBOX {
		host = "sandbox.evernote.com"
	}
	return host
}

func GetUserStore(envType EnvironmentType) (*userstore.UserStoreClient, error) {
	url := fmt.Sprintf("https://%s/edam/user", GetHost(envType))
	c, err := thrift.NewTHttpPostClient(url)
	if err != nil {
		return nil, err
	}
	return userstore.NewUserStoreClientFactory(
		c,
		thrift.NewTBinaryProtocolFactoryDefault(),
	), nil
}

func GetNoteStore(userstore *userstore.UserStoreClient, authenticationToken string) (*notestore.NoteStoreClient, error) {
	urls, err := userstore.GetUserUrls(authenticationToken)
	if err != nil {
		return nil, err
	}

	url := urls.GetNoteStoreUrl()
	c, err := thrift.NewTHttpPostClient(url)
	if err != nil {
		return nil, err
	}

	return notestore.NewNoteStoreClientFactory(
		c,
		thrift.NewTBinaryProtocolFactoryDefault(),
	), nil
}

package main

import (
	"github.com/go-ldap/ldap/v3"
)

var GConn *ldap.Conn

func entryLdap() {
	conn, err := ldapConn()
	if err != nil {
		logme.Fatal(err)
	}
	GConn = conn
}
func ldapConn() (*ldap.Conn, error) {
	for _, v := range Config.GetStringSlice("ldap.servers") {
		ldapConnZ, err := ldap.DialURL(v)
		if err != nil {
			logme.Info(err.Error())
			continue
		}

		err = ldapConnZ.Bind(Config.GetString("ldap.auth.user"), Config.GetString("ldap.auth.password"))
		if err != nil {
			logme.Info(err.Error())
			continue
		}
		return ldapConnZ, nil
	}
	return nil, ErrNotConnectToAnyServersLdap
}

func ldapPrepareSearchRequest() (*ldap.SearchRequest, error) {
	searchRequestDN := Config.GetString("asterisk.extension.dn")
	if searchRequestDN == "" {
		searchRequestDN = Config.GetString("ldap.find.defaultDN")
	}
	if searchRequestDN == "" {
		return &ldap.SearchRequest{}, ErrNotDNForSearchExtension

	}
	searchRequest := ldap.NewSearchRequest(
		searchRequestDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"",
		nil,
		nil,
	)
	return searchRequest, nil

}
func ldapSearch(searchRequest ldap.SearchRequest) (ldap.SearchResult, error) {

	req, err := GConn.Search(&searchRequest)
	if err != nil {
		return ldap.SearchResult{}, err
	}
	if len(req.Entries) == 0 {
		return ldap.SearchResult{}, nil
	}
	//logme.Infof("Finding %d entries in ldap", len(req.Entries))
	return *req, nil
}
func ldapUpdateUserTelephone(cn, number string) error {
	x := ldap.NewModifyRequest(cn, nil)
	x.Replace(Config.GetString("asterisk.extension.owner.UserAttributePhone"), []string{number})
	err := GConn.Modify(x)
	return err
}
func ldapUpdateUserClean(cn string) error {
	x := ldap.NewModifyRequest(cn, nil)
	x.Replace(Config.GetString("asterisk.extension.owner.UserAttributePhone"), nil)
	err := GConn.Modify(x)
	return err
}

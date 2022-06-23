package main

import "errors"

var (
	ErrNotConnectToAnyServersLdap      = errors.New("not Connect to any server ldap")
	ErrNotDNForSearchExtension         = errors.New("not DN for search Extension use (asterisk.extension.dn) or (ldap.find.defaultDN) in config file")
	ErrNotSearchResulFromLdapExtension = errors.New("not result search extension in ldap")
)

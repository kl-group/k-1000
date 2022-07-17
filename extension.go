package main

import (
	"fmt"
	"regexp"
)

var VarMapExtensionLdap MapExtension
var VarMapExtensionAsterisk MapExtension

type MapExtension struct {
	Map map[string]Extension
}

type Extension struct {
	Owner       string
	Name        string
	LdapDN      string
	Number      string
	Description string
	UserId      string
	ExtPassword string
}

func (T *MapExtension) loadFromLdap() error {
	mapExtensionLdap := make(map[string]Extension)
	searchRequest, err := ldapPrepareSearchRequest()
	if err != nil {
		return err
	}
	searchRequest.Filter = Config.GetString("asterisk.extension.ldapfiler")
	searchRequest.Attributes = []string{Config.GetString("asterisk.extension.attribute"), Config.GetString("asterisk.extension.owner.attribute"), "description", Config.GetString("asterisk.extension.owner.passwdattr")}
	req, err := ldapSearch(*searchRequest)
	if err != nil {
		logme.Fatal(err.Error())
	}
	for _, v := range req.Entries {
		objectName := v.GetAttributeValue(Config.GetString("asterisk.extension.attribute"))
		matched, _ := regexp.MatchString(Config.GetString("asterisk.extension.attribute-regexp"), objectName)
		if !matched {
			logme.Infof("Skip entry (%s) not mathed regexp (%s) attribute (%s) in value (%s)", v.DN, Config.GetString("asterisk.extension.attribute-regexp"), Config.GetString("asterisk.extension.attribute"), objectName)
			continue
		}

		var re = regexp.MustCompile(Config.GetString("asterisk.extension.attribute-regexp"))
		s := re.ReplaceAllString(objectName, `$1`)
		if err != nil {
			logme.Warning(err)
			continue
		}
		x := Extension{
			Owner:       v.GetAttributeValue(Config.GetString("asterisk.extension.owner.attribute")),
			LdapDN:      v.DN,
			Number:      s,
			Description: v.GetAttributeValue("description"),
			ExtPassword: v.GetAttributeValue(Config.GetString("asterisk.extension.owner.passwdattr")),
		}
		x.Name = helpersGetNameFromDN(x.Owner)
		mapExtensionLdap[s] = x
	}

	T.Map = mapExtensionLdap
	return nil
}
func (T *MapExtension) loadFromAsterisk() error {
	mapExtensionLdap := make(map[string]Extension)
	sl, err := GraphQLLoadAllExtension()
	if err != nil {
		return err
	}
	for _, v := range sl {
		mapExtensionLdap[v.Number] = v
	}
	T.Map = mapExtensionLdap
	return nil
}
func (T MapExtension) prepareLdap() {
	for _, v := range T.Map {
		if err := v.updateOwnerLdap(); err != nil {
			logme.Warning(err)
			continue
		}
	}

}
func (T MapExtension) prepareAsterisk() {
	for _, v := range T.Map {
		if err := v.checkExtensionAsterisk(); err != nil {
			logme.Warning(err)
			continue
		}
	}

	for _, lv := range VarMapExtensionLdap.Map {
		if _, ok := VarMapExtensionAsterisk.Map[lv.Number]; !ok {
			if err := lv.createExtensionInAsterisk(); err != nil {
				logme.Warning(err)
				continue
			}
		}
	}
}

func (T Extension) checkExtensionAsterisk() error {
	if _, ok := VarMapExtensionLdap.Map[T.Number]; !ok {
		logme.Infof("Remove Extension '%s' from Asterisk", T.Number)
		return nil
	}
	val := VarMapExtensionLdap.Map[T.Number]
	needName := val.Name
	if needName == "" && val.Description != "" {
		needName = val.Description
	}
	if needName == "" {
		logme.Warningf("Name Extension '%s' empty in ldap", T.Number)
		return nil
	}
	if needName == T.Name {
		return nil
	}
	logme.Infof("Rename Extension '%s' old %s new %s", T.Number, T.Name, needName)
	return GraphQLUpdateNameExtension(T, needName)
}

func (T Extension) createExtensionInAsterisk() error {
	return GraphQlAddExtension(T)

}
func (T Extension) updateOwnerLdap() error {
	searchRequest, err := ldapPrepareSearchRequest()
	if err != nil {
		return err
	}
	searchRequest.BaseDN = Config.GetString("ldap.find.userDN")
	searchRequest.Filter = fmt.Sprintf("(&(objectClass=user)(%s=%s))", Config.GetString("asterisk.extension.owner.UserAttributePhone"), T.Number)
	searchRequest.Attributes = []string{Config.GetString("asterisk.extension.owner.UserAttributePhone")}

	req, err := ldapSearch(*searchRequest)
	if err != nil {
		return err
	}
	ownerIs := false
	for _, v := range req.Entries {
		if v.DN == T.Owner {
			ownerIs = true
			continue
		}

		logme.Infof("remove telephone number for user %s", v.DN)
		err = ldapUpdateUserClean(v.DN)
		if err != nil {
			logme.Warning(err)
			continue
		}
	}

	if !ownerIs && T.Owner != "" {
		logme.Infof("set telephone number for user %s, number = %s", T.Owner, T.Number)
		err = ldapUpdateUserTelephone(T.Owner, T.Number)
		if err != nil {
			return err
		}
	}

	return nil
}

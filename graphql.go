package main

import (
	"context"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2/clientcredentials"
	"strconv"
)

var ctx = context.Background()

func GraphQLUpdateNameExtension(e Extension, name string) error {
	id, err := strconv.Atoi(e.Number)
	if err != nil {
		logme.Warningf("string \"%s\" not convert to int", e.Number)
		return err
	}
	client, err := GraphQLClient()
	if err != nil {
		return err
	}

	var vRes struct {
		UpdateExtension struct {
			Status  graphql.Boolean `json:"status"`
			Message graphql.String  `json:"message"`
		} `graphql:"updateExtension(input: $input)" json:"updateExtension"`
	}

	type updateExtensionInput struct {
		ExtensionId graphql.ID     `json:"extensionId" graphql:"extensionId"`
		Name        graphql.String `json:"name" graphql:"name"`
		ExtPassword graphql.String `json:"extPassword" graphql:"extPassword"`
	}
	x := updateExtensionInput{
		ExtensionId: graphql.ID(id),
		Name:        graphql.String(name),
		ExtPassword: graphql.String(e.ExtPassword),
	}
	variables := map[string]interface{}{
		"input": x,
	}
	err = client.Mutate(ctx, &vRes, variables)
	if err != nil {
		return err
	}
	needReloadAsterisk = true
	return nil
}
func GraphQlAddExtension(e Extension) error {
	logme.Infof("Create Extension (%s) in Asterisk", e.Number)
	id, err := strconv.Atoi(e.Number)
	if err != nil {
		logme.Warningf("string \"%s\" not convert to int", e.Number)
		return err
	}

	client, err := GraphQLClient()
	if err != nil {
		return err
	}
	var vRes struct {
		AddExtension struct {
			Status  graphql.Boolean `json:"status"`
			Message graphql.String  `json:"message"`
		} `graphql:"addExtension(input: $input)" json:"addExtension"`
	}
	type addExtensionInput struct {
		ExtensionId graphql.ID     `json:"extensionId" graphql:"extensionId"`
		Name        graphql.String `json:"name" graphql:"name"`
		Tech        graphql.String `json:"tech" graphql:"tech"`
		Email       graphql.String `json:"email" graphql:"email"`
	}
	x := addExtensionInput{
		ExtensionId: graphql.ID(id),
		Tech:        graphql.String("pjsip"),
		Name:        graphql.String(e.Name),
		Email:       graphql.String(""),
	}
	variables := map[string]interface{}{
		"input": x,
	}
	err = client.Mutate(ctx, &vRes, variables)
	if err != nil {
		return err
	}
	needReloadAsterisk = true
	return nil
}
func GraphQLLoadAllExtension() ([]Extension, error) {
	client, err := GraphQLClient()
	if err != nil {
		return nil, err
	}

	var vres struct {
		FetchAllExtensions struct {
			Status    graphql.Boolean
			Extension []struct {
				ExtensionId graphql.String
				User        struct {
					Name        graphql.String
					Id          graphql.String
					ExtPassword graphql.String
				}
				CoreDevice struct {
					DeviceId    graphql.String
					Description graphql.String
				}
			}
		}
	}
	err = client.Query(ctx, &vres, nil)
	if err != nil {
		return nil, err
	}
	var sl []Extension
	for _, v := range vres.FetchAllExtensions.Extension {
		x := Extension{
			Owner:       "",
			Name:        string(v.User.Name),
			LdapDN:      "",
			Number:      string(v.ExtensionId),
			UserId:      string(v.User.Id),
			ExtPassword: string(v.User.ExtPassword),
		}
		sl = append(sl, x)
	}
	return sl, nil
}
func GraphQLReloadAsterisk() error {
	client, err := GraphQLClient()
	if err != nil {
		return err
	}
	var req struct {
		UpdateExtension struct {
			Status        graphql.Boolean `json:"status"`
			Message       graphql.String  `json:"message"`
			TransactionId graphql.ID      `json:"transaction_id" graphql:"transaction_id"`
		} `graphql:"doreload(input: {})" json:"doreload"`
	}
	err = client.Mutate(ctx, &req, nil)
	if err != nil {
		return err
	}
	return nil

}
func GraphQLClient() (*graphql.Client, error) {
	conf := clientcredentials.Config{
		ClientID:     Config.GetString("asterisk.auth.clientid"),
		ClientSecret: Config.GetString("asterisk.auth.secret"),
		TokenURL:     Config.GetString("asterisk.auth.tokenUrl"),
		Scopes:       []string{},
	}
	httpClient := conf.Client(ctx)
	client := graphql.NewClient(Config.GetString("asterisk.auth.qlUrl"), httpClient)
	return client, nil
}

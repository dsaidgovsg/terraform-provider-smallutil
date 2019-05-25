package main

import (
	"crypto"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func hash(objs ...interface{}) []byte {
	digester := crypto.MD5.New()
	for _, ob := range objs {
		fmt.Fprint(digester, reflect.TypeOf(ob))
		fmt.Fprint(digester, ob)
	}
	return digester.Sum(nil)
}

func hashString(objs ...interface{}) string {
	return string(hash(objs))
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}

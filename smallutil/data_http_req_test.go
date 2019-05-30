package smallutil

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestDataHTTPReqRead(t *testing.T) {
	resource.Test(t, resource.TestCase{
		// PreCheck:  func() { testAccPreCheck(t) },
		// Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Config: testAccCheckAwsAmiDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckAwsAmiDataSourceID("data.aws_ami.nat_ami"),
					resource.TestCheckResourceAttr("data.smallutil_http_req2.xxx", "url", "http://localhost:8080/"),
				),
			},
		},
	})
}

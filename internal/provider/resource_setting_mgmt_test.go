package provider

import (
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var settingMgmtLock = sync.Mutex{}

func TestAccSettingMgmt_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			settingMgmtLock.Lock()
			t.Cleanup(func() {
				settingMgmtLock.Unlock()
			})
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingMgmtConfig_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			importStep("unifi_setting_mgmt.test"),
		},
	})
}

func TestAccSettingMgmt_site(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			settingMgmtLock.Lock()
			t.Cleanup(func() {
				settingMgmtLock.Unlock()
			})
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingMgmtConfig_site(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      "unifi_setting_mgmt.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_setting_mgmt.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSettingMgmt_sshKeys(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			settingMgmtLock.Lock()
			t.Cleanup(func() {
				settingMgmtLock.Unlock()
			})
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingMgmtConfig_sshKeys(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      "unifi_setting_mgmt.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_setting_mgmt.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSettingMgmtConfig_basic() string {
	return `
resource "unifi_setting_mgmt" "test" {
	auto_upgrade = true
}
`
}

func testAccSettingMgmtConfig_site() string {
	return `
resource "unifi_site" "test" {
	description = "test"
}

resource "unifi_setting_mgmt" "test" {
	site = unifi_site.test.name
	auto_upgrade = true
}
`
}

func testAccSettingMgmtConfig_sshKeys() string {
	return `
resource "unifi_site" "test" {
	description = "test"
}

resource "unifi_setting_mgmt" "test" {
	site = unifi_site.test.name
	ssh_enabled = true
	ssh_key {
		name = "Test key"
		type = "ssh-rsa"
		key = "AAAAB3NzaC1yc2EAAAADAQABAAACAQDNWqT8zvVtmaks7sLlP+hmWmJVmruyNU9uk8JpLTX0oE+r9hjePsXCThTrft7s+vlaj+bLr8Yf5//TT8KS7LB/YIp2O3jPomOz9A4hIsG5R6FLfSggzQP4a7QSlNLCm/6WjKHP9DhRb7trnFz+KkCNmCVKLZgiyeUm2LydVKJ2QncHopA5yomtSpmb6x66zaKr+DbwzHC13WIEms5Ros0N9pEOcAghsSEVL42bfGBfSH37R+Kaw0nhWei4Y25jO66xsbtyZKoiF1+XXXBuEi77Tv7iQGHHOFRqNKKfGI1QhYvwlcjdzh9wu7Gtzeyh/+jpF8mwCLtFKle+W/zSs+lHCuCihvQEQtCIpZL5FapvxfxPZQJWL5RgsL9jieUaoF8EsWAOM83BCSZa/FB1RyfKdy4f7BQtDCKIm3nD5paCJSfS6DSw1TMvaFPeJLG3PuyHRbNvbVLmHRl9lK03na6/R9JX06nBUuPdP+FLjIZsyZz1DOUSDjCWHFk0+Ne2uEinV7SkOoxC6E2NxqlY/SyMnWZS+p95Zx6yOlNqB9sQ+Q4/YLGY5mUmqJrHPlH6LjXfudybKHMZUuVRF1NX3ESue8NSKc0SlJDQUXtJ9wkjjX1wAWvXCDwI72jtC86r/wzw+mcIfpks3jHQrOhpwCRmQL4vAs5DztA3hKxkgElYaw=="
		comment = "test@example.com"
	}
}
`
}

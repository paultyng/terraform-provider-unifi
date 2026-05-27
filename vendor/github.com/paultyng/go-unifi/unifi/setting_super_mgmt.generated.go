// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// just to fix compile issues with the import
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type SettingSuperMgmt struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Key string `json:"key"`

	AnalyticsDisapprovedFor                  string   `json:"analytics_disapproved_for,omitempty"`
	AutoUpgrade                              bool     `json:"auto_upgrade"`
	AutobackupCronExpr                       string   `json:"autobackup_cron_expr,omitempty"`
	AutobackupDays                           int      `json:"autobackup_days,omitempty"`
	AutobackupEnabled                        bool     `json:"autobackup_enabled"`
	AutobackupGcsBucket                      string   `json:"autobackup_gcs_bucket,omitempty"`
	AutobackupGcsCertificatePath             string   `json:"autobackup_gcs_certificate_path,omitempty"`
	AutobackupLocalPath                      string   `json:"autobackup_local_path,omitempty"`
	AutobackupMaxFiles                       int      `json:"autobackup_max_files,omitempty"`
	AutobackupPostActions                    []string `json:"autobackup_post_actions,omitempty"` // copy_local|copy_s3|copy_gcs|copy_cloud
	AutobackupS3AccessKey                    string   `json:"autobackup_s3_access_key,omitempty"`
	AutobackupS3AccessSecret                 string   `json:"autobackup_s3_access_secret,omitempty"`
	AutobackupS3Bucket                       string   `json:"autobackup_s3_bucket,omitempty"`
	AutobackupTimezone                       string   `json:"autobackup_timezone,omitempty"`
	BackupToCloudEnabled                     bool     `json:"backup_to_cloud_enabled"`
	ContactInfoCity                          string   `json:"contact_info_city,omitempty"`
	ContactInfoCompanyName                   string   `json:"contact_info_company_name,omitempty"`
	ContactInfoCountry                       string   `json:"contact_info_country,omitempty"`
	ContactInfoFullName                      string   `json:"contact_info_full_name,omitempty"`
	ContactInfoPhoneNumber                   string   `json:"contact_info_phone_number,omitempty"`
	ContactInfoShippingAddress1              string   `json:"contact_info_shipping_address_1,omitempty"`
	ContactInfoShippingAddress2              string   `json:"contact_info_shipping_address_2,omitempty"`
	ContactInfoState                         string   `json:"contact_info_state,omitempty"`
	ContactInfoZip                           string   `json:"contact_info_zip,omitempty"`
	DataRetentionSettingPreference           string   `json:"data_retention_setting_preference,omitempty"` // auto|manual
	DataRetentionTimeInHoursFor5MinutesScale int      `json:"data_retention_time_in_hours_for_5minutes_scale,omitempty"`
	DataRetentionTimeInHoursForDailyScale    int      `json:"data_retention_time_in_hours_for_daily_scale,omitempty"`
	DataRetentionTimeInHoursForHourlyScale   int      `json:"data_retention_time_in_hours_for_hourly_scale,omitempty"`
	DataRetentionTimeInHoursForMonthlyScale  int      `json:"data_retention_time_in_hours_for_monthly_scale,omitempty"`
	DataRetentionTimeInHoursForOthers        int      `json:"data_retention_time_in_hours_for_others,omitempty"`
	DefaultSiteDeviceAuthPasswordAlert       string   `json:"default_site_device_auth_password_alert,omitempty"` // false
	Discoverable                             bool     `json:"discoverable"`
	EnableAnalytics                          bool     `json:"enable_analytics"`
	GoogleMapsApiKey                         string   `json:"google_maps_api_key,omitempty"`
	ImageMapsUseGoogleEngine                 bool     `json:"image_maps_use_google_engine"`
	LedEnabled                               bool     `json:"led_enabled"`
	LiveChat                                 string   `json:"live_chat,omitempty"`    // disabled|super-only|everyone
	LiveUpdates                              string   `json:"live_updates,omitempty"` // disabled|live|auto
	MinimumUsableHdSpace                     int      `json:"minimum_usable_hd_space,omitempty"`
	MinimumUsableSdSpace                     int      `json:"minimum_usable_sd_space,omitempty"`
	MultipleSitesEnabled                     bool     `json:"multiple_sites_enabled"`
	OverrideInformHost                       bool     `json:"override_inform_host"`
	OverrideInformHostLocation               string   `json:"override_inform_host_location,omitempty"`
	StoreEnabled                             string   `json:"store_enabled,omitempty"` // disabled|super-only|everyone
	TimeSeriesPerClientStatsEnabled          bool     `json:"time_series_per_client_stats_enabled"`
	XSshPassword                             string   `json:"x_ssh_password,omitempty"`
	XSshUsername                             string   `json:"x_ssh_username,omitempty"`
}

func (dst *SettingSuperMgmt) UnmarshalJSON(b []byte) error {
	type Alias SettingSuperMgmt
	aux := &struct {
		AutobackupDays                           emptyStringInt `json:"autobackup_days"`
		AutobackupMaxFiles                       emptyStringInt `json:"autobackup_max_files"`
		DataRetentionTimeInHoursFor5MinutesScale emptyStringInt `json:"data_retention_time_in_hours_for_5minutes_scale"`
		DataRetentionTimeInHoursForDailyScale    emptyStringInt `json:"data_retention_time_in_hours_for_daily_scale"`
		DataRetentionTimeInHoursForHourlyScale   emptyStringInt `json:"data_retention_time_in_hours_for_hourly_scale"`
		DataRetentionTimeInHoursForMonthlyScale  emptyStringInt `json:"data_retention_time_in_hours_for_monthly_scale"`
		DataRetentionTimeInHoursForOthers        emptyStringInt `json:"data_retention_time_in_hours_for_others"`
		MinimumUsableHdSpace                     emptyStringInt `json:"minimum_usable_hd_space"`
		MinimumUsableSdSpace                     emptyStringInt `json:"minimum_usable_sd_space"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.AutobackupDays = int(aux.AutobackupDays)
	dst.AutobackupMaxFiles = int(aux.AutobackupMaxFiles)
	dst.DataRetentionTimeInHoursFor5MinutesScale = int(aux.DataRetentionTimeInHoursFor5MinutesScale)
	dst.DataRetentionTimeInHoursForDailyScale = int(aux.DataRetentionTimeInHoursForDailyScale)
	dst.DataRetentionTimeInHoursForHourlyScale = int(aux.DataRetentionTimeInHoursForHourlyScale)
	dst.DataRetentionTimeInHoursForMonthlyScale = int(aux.DataRetentionTimeInHoursForMonthlyScale)
	dst.DataRetentionTimeInHoursForOthers = int(aux.DataRetentionTimeInHoursForOthers)
	dst.MinimumUsableHdSpace = int(aux.MinimumUsableHdSpace)
	dst.MinimumUsableSdSpace = int(aux.MinimumUsableSdSpace)

	return nil
}

func (c *Client) getSettingSuperMgmt(ctx context.Context, site string) (*SettingSuperMgmt, error) {
	var respBody struct {
		Meta meta               `json:"meta"`
		Data []SettingSuperMgmt `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/get/setting/super_mgmt", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) updateSettingSuperMgmt(ctx context.Context, site string, d *SettingSuperMgmt) (*SettingSuperMgmt, error) {
	var respBody struct {
		Meta meta               `json:"meta"`
		Data []SettingSuperMgmt `json:"data"`
	}

	d.Key = "super_mgmt"
	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/set/setting/super_mgmt", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

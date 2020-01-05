// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type SettingSuperMgmt struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

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
	DataRetentionTimeEnabled                 bool     `json:"data_retention_time_enabled"`
	DataRetentionTimeInHoursFor5minutesScale int      `json:"data_retention_time_in_hours_for_5minutes_scale,omitempty"`
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
	OverrideInformHost                       bool     `json:"override_inform_host"`
	StoreEnabled                             string   `json:"store_enabled,omitempty"` // disabled|super-only|everyone
	TimeSeriesPerClientStatsEnabled          bool     `json:"time_series_per_client_stats_enabled"`
	UploadDebuggingInfoEnabled               bool     `json:"upload_debugging_info_enabled"`
	XSshPassword                             string   `json:"x_ssh_password,omitempty"`
	XSshUsername                             string   `json:"x_ssh_username,omitempty"`
}

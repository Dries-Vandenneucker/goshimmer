package dashboard

import (
	flag "github.com/spf13/pflag"
)

const (
	// CfgBindAddress defines the config flag of the analysis dashboard binding address.
	CfgBindAddress = "analysis.dashboard.bindAddress"
	// CfgDev defines the config flag of the analysis dashboard dev mode.
	CfgDev = "analysis.dashboard.dev"
	// CfgBasicAuthEnabled defines the config flag of the analysis dashboard basic auth enabler.
	CfgBasicAuthEnabled = "analysis.dashboard.basic_auth.enabled"
	// CfgBasicAuthUsername defines the config flag of the analysis dashboard basic auth username.
	CfgBasicAuthUsername = "analysis.dashboard.basic_auth.username"
	// CfgBasicAuthPassword defines the config flag of the analysis dashboard basic auth password.
	CfgBasicAuthPassword = "analysis.dashboard.basic_auth.password"
	// CfgFPCBufferSize defines the config flag of the analysis dashboard FPC (past conflicts) buffer size.
	CfgFPCBufferSize = "analysis.dashboard.fpc.buffer_size"
)

func init() {
	flag.String(CfgBindAddress, "0.0.0.0:8888", "the bind address of the dashboard")
	flag.Bool(CfgDev, false, "whether the dashboard runs in dev mode")
	flag.Bool(CfgBasicAuthEnabled, false, "whether to enable HTTP basic auth")
	flag.String(CfgBasicAuthUsername, "goshimmer", "HTTP basic auth username")
	flag.String(CfgBasicAuthPassword, "goshimmer", "HTTP basic auth password")
	flag.Uint32(CfgFPCBufferSize, 200, "FPC buffer size")
}

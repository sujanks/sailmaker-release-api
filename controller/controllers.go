package controller

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/functions"
	"github.com/skhatri/api-router-go/router/settings"
)

func Configure(configurer router.ApiConfigurer) {
	var _settings = settings.GetSettings()
	configurer.Get("/api/namespaces", namespaceApiHandler).
		Get("/status", functions.StatusFunc).
		Get("/api/deployments", fetchDeployments).
		Get("/api/statefulsets", fetchStatefulsets).
		Get("/api/jobs", fetchJobs).
		GetIf(_settings.IsToggleOn("can_read_crds")).
			Add("/api/crd-instances", getCrdInstanceList).
			Add("/api/crd-instance", getCrdInstance).
			Add("/api/crds", getCrds).
		Done().
		PostIf(_settings.IsToggleOn("can_write_crds")).Register("/api/release", performRelease).
		GetIf(_settings.IsToggleOn("daemonset_endpoint")).Register("/api/daemonsets", fetchDaemonsets)
}

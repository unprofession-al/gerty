package config

type Configuration struct {
	Port              string `json:"port"`
	Address           string `json:"address"`
	Store             string `json:"store"`
	NodeVarsProviders string `json:"nodevars_providers"`
	JenkinsFileName   string `json:"jenkins_file_name"`
	JenkinsToken      string `json:"-"`
	JenkinsJobName    string `json:"jenkins_job_name"`
	JenkinsBaseUrl    string `json:"jenkins_base_url"`
}

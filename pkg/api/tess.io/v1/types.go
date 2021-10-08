//types.go
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CiConfig is a specification for CI
type CiConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status CiConfigStatus `json:"status"`
	Spec   CiConfigSpec   `json:"spec"`
}

type CiConfigStatus struct {
	DNSName string `json:"dnsName"`
	Phase   string `json:"phase"`
}

// CiConfigSpec is the spec for a CiConfig resource
type CiConfigSpec struct {
	Name      string               `json:"name"`
	Strategy  CiConfigSpecStrategy `json:"strategy"`
	Source    CiConfigSpecSource   `json:"source"`
	Hibernate bool                 `json:"hibernate"`
}

type CiConfigSpecSource struct {
	Git CiConfigSpecSourceGit `json:"git"`
}

type CiConfigSpecSourceGit struct {
	Uri string `json:"uri"`
}

type CiConfigSpecStrategy struct {
	StandardBuild CiConfigSpecStrategyStandardBuild `json:"standardBuild"`
	Master        CiConfigSpecStrategyMaster        `json:"master"`
	Env           []CiConfigSpecStrategyEnv         `json:"env"`
}

type CiConfigSpecStrategyEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CiConfigSpecStrategyStandardBuild struct {
	Spec CiConfigSpecStrategyStandardBuildSpec `json:"spec"`
}

type CiConfigSpecStrategyMaster struct {
	Image      string `json:"image"`
	VolumeSize int    `json:"volumeSize"`
}

type CiConfigSpecStrategyStandardBuildSpec struct {
	Stack       CiConfigSpecStrategyStandardBuildSpecStack      `json:"stack"`
	Identifiers CiConfigSpecStrategyStandardBuildSpecIdentifier `json:"identifiers"`
}

type CiConfigSpecStrategyStandardBuildSpecIdentifier struct {
	AppName          string `json:"appName"`
	AssemblerPomPath string `json:"assemblerPomPath"`
	Owner            string `json:"owner"`
	ServiceName      string `json:"serviceName"`
}

type CiConfigSpecStrategyStandardBuildSpecStack struct {
	Builder string `json:"builder"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CiConfigList is a list of CI resources
type CiConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CiConfig `json:"items"`
}

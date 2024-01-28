/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PodAppSpec defines the desired state of PodApp
type PodAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of PodApp. Edit podapp_types.go to remove/update
	Enable       bool    `json:"Enable,omitempty"`
	PodName      string  `json:"PodName,omitempty"`
	PodNamespace string  `json:"PodNamespace,omitempty"`
	PodSpec      PodSpec `json:"PodSpec"`
}

type PodSpec struct {
	// Add custom pod specifications here
	Containers []ContainerSpec `json:"containers"`
}

type ContainerSpec struct {
	// Name is the name of the container
	Name string `json:"name"`

	// Image is the container image for the container
	Image string `json:"image"`
}

// PodAppStatus defines the observed state of PodApp
type PodAppStatus struct {
	PodStatus PodStatus `json:"podStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type PodStatus string

const (
	PodStatusCreating PodStatus = "Creating"
	PodStatusRunning  PodStatus = "Running"
	PodStatusFailed   PodStatus = "Failed"
	PodStatusUnknown  PodStatus = "Unknown"
)

// +kubebuilder:object:root=true
// PodApp is the Schema for the podapps API
type PodApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodAppSpec   `json:"spec,omitempty"`
	Status PodAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodAppList contains a list of PodApp
type PodAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodApp{}, &PodAppList{})
}

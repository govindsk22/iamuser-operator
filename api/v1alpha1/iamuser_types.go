/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IamUserSpec defines the desired state of IamUser




type IamUserSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Username string `json:"username"`
	// +kubebuilder:validation:Enum=admin;readonly
	Role string `json:"role"`
}

// IamUserStatus defines the observed state of IamUser
type IamUserStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Usercreated bool   `json:"usercreated"`
	UserArn     string `json:"userarn"`
	Username 	string `json:"username"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="UserCreated",type=boolean,JSONPath=`.status.usercreated`
//+kubebuilder:printcolumn:name="UserArn",type=string,JSONPath=`.status.userarn`
//+kubebuilder:printcolumn:name="UserRole",type=string,JSONPath=`.spec.role`
// IamUser is the Schema for the iamusers API
type IamUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IamUserSpec   `json:"spec,omitempty"`
	Status IamUserStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IamUserList contains a list of IamUser
type IamUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IamUser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IamUser{}, &IamUserList{})
}

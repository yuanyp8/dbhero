/*
Copyright 2022 yuanyp8.

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

// DatabaseSpec defines the desired state of Database
type DatabaseSpec struct {

	//+optional
	//kubebuilder:default=utf8mb4
	//+kubebuilder:validation:Enum={utf8mb4,latin1}
	DefaultCharacterSet CharacterSet `json:"defaultCharacterSet,omitempty"`

	//+optional
	//kubebuilder:default=utf8mb4_general_ci
	//+kubebuilder:validation:Enum={utf8mb4_general_ci,latin1_general_cs}
	DefaultCollation Collation `json:"defaultCollation,omitempty"`

	//+optional
	//+kubebuilder:default=mysql
	//+kubebuilder:validation:Enum={mysql,postgresql}
	DBType string `json:"DBType,omitempty"`

	//+kubebuilder:default=5.7
	//+kubebuilder:validation:Enum={5.7,8.0}
	DBVersion string `json:"DBVersion,omitempty"`

	//+optional
	//+kubebuilder:validation:MinLength:=6
	//+kubebuilder:validation:MaxLength:=16
	DBName string `json:"dbName,omitempty"`

	Auth Auth `json:"auth,omitempty"`
}

// DatabaseStatus defines the observed state of Database
type DatabaseStatus struct {
	DBName string        `json:"dbName,omitempty"`
	Auth   Auth          `json:"auth,omitempty"`
	Access AccessAddress `json:"access,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Database is the Schema for the databases API
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSpec   `json:"spec,omitempty"`
	Status DatabaseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DatabaseList contains a list of Database
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}

func (in *Database) SetAnnotation() *Database {
	if in.GetAnnotations() == nil {
		in.Annotations = map[string]string{}
	}
	return in
}

func (in *Database) SetLabels(labels map[string]string) {
	if in.ObjectMeta.Labels == nil {
		in.ObjectMeta.Labels = map[string]string{
			"control-plane": "dbhero",
		}
	}

	if len(labels) != 0 {
		for k, v := range labels {
			in.ObjectMeta.Labels[k] = v
		}
	}
	return
}

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}

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
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DataSourceConnection struct {
	//+optional

	MySQL *MysqlConnection `json:"mysql,omitempty"`

	//+optional
	// TODO(yuanyp8): not implement
	Postgre *PostgreConnection `json:"postgre,omitempty"`
}

// DataSourceSpec defines the desired state of DataSource
type DataSourceSpec struct {
	Connection DataSourceConnection `json:"connection,omitempty"`
}

// DataSourceStatus defines the observed state of DataSource
type DataSourceStatus struct {
	//+kubebuilder:default=false
	//+optional

	IsConnected bool   `json:"isConnected"`
	LastPing    string `json:"lastPing"`

	//+kubebuilder:default="UnKnown"
	//+optional

	Type    DataSourceType `json:"type"`
	Auth    Auth           `json:"auth,omitempty"`
	Version string         `json:"version,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Namespace",type=string,JSONPath=`.metadata.namespace`,priority=1
//+kubebuilder:printcolumn:name="TYPE",type="string",JSONPath=".status.type"
//+kubebuilder:printcolumn:name="Version",type="string",JSONPath=".status.version"
//+kubebuilder:printcolumn:name="CONNECTED",type="boolean",JSONPath=".status.isConnected"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// DataSource is the Schema for the datasource API
type DataSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataSourceSpec   `json:"spec,omitempty"`
	Status DataSourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DataSourceList contains a list of DataSource
type DataSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataSource{}, &DataSourceList{})
}

func (in *DataSource) GetType() (DataSourceType, error) {
	if in.Spec.Connection.MySQL != nil {
		return MYSQL, nil
	}
	if in.Spec.Connection.Postgre != nil {
		return POSTGRE, nil
	}
	return UNKNOWN, fmt.Errorf("no database connection configured for database: %s", in.Name)

}

func (in *DataSource) GetMysqlAuth() Auth {
	return in.Spec.Connection.MySQL.Auth
}

func (in *DataSource) GetPostgreAuth() Auth {
	// TODO: not implement
	return Auth{}
}

func (in *DataSource) SetAnnotation() *DataSource {
	if in.GetAnnotations() == nil {
		in.Annotations = map[string]string{}
	}
	return in
}

func (in *DataSource) SetLabels(labels map[string]string) {
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

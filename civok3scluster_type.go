package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CivoK3sClusterConditionType defines possible states of a CivoK3sCluster
type CivoK3sClusterConditionType string

// CivoK3sClusterSpec defines the desired state of CivoK3sCluster
type CivoK3sClusterSpec struct {
	Masters      int64   `json:"masters"`
	MasterSize   *string `json:"master_size,omitempty"`
	Version      string  `json:"version"`
	StorageClass string  `json:"storageclass,omitempty"`
	CNIPlugin    string  `json:"cni_plugin,omitempty"`
}

// CivoK3sClusterStatus defines the observed state of CivoK3sCluster
type CivoK3sClusterStatus struct {
	State          string `json:"state,omitempty"`
	ClusterVersion string `json:"clusterversion,omitempty"`
}

// CivoK3sCluster is the Schema for the civok3sclusters API
type CivoK3sCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CivoK3sClusterSpec   `json:"spec,omitempty"`
	Status CivoK3sClusterStatus `json:"status,omitempty"`
}

// CivoK3sClusterList contains a list of CivoK3sCluster
type CivoK3sClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CivoK3sCluster `json:"items"`
}

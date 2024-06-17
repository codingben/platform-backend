package types

import (
	cappv1 "github.com/dana-team/container-app-operator/api/v1alpha1"
)

type CappRevision struct {
	Metadata    Metadata                  `json:"metadata" binding:"required"`
	Annotations []KeyValue                `json:"annotations" binding:"required"`
	Labels      []KeyValue                `json:"labels" binding:"required"`
	Spec        cappv1.CappRevisionSpec   `json:"spec" binding:"required"`
	Status      cappv1.CappRevisionStatus `json:"status" binding:"required"`
}

type CappRevisionList struct {
	CappRevisions []cappv1.CappRevision `json:"capprevisions"`
	Count         int                   `json:"count"`
}

type CappRevisionNamespaceUri struct {
	NamespaceName string `uri:"namespaceName" binding:"required"`
}

type CappRevisionUri struct {
	NamespaceName    string `uri:"namespaceName" binding:"required"`
	CappRevisionName string `uri:"cappRevisionName" binding:"required"`
}

type CappRevisionQuery struct {
	Labels []KeyValue `form:"labels"`
}

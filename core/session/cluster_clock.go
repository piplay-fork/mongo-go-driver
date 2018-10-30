// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package session

import (
	"sync"

	"github.com/mongodb/mongo-go-driver/bson"
)

// ClusterClock represents a logical clock for keeping track of cluster time.
type ClusterClock struct {
	clusterTime bson.Doc
	lock        sync.Mutex
}

// GetClusterTime returns the cluster's current time.
func (cc *ClusterClock) GetClusterTime() bson.Doc {
	var ct bson.Doc
	cc.lock.Lock()
	ct = cc.clusterTime
	cc.lock.Unlock()

	return ct
}

// AdvanceClusterTime updates the cluster's current time.
func (cc *ClusterClock) AdvanceClusterTime(clusterTime bson.Doc) {
	cc.lock.Lock()
	cc.clusterTime = MaxClusterTime(cc.clusterTime, clusterTime)
	cc.lock.Unlock()
}

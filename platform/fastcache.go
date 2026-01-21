// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
package main

import "context"

var fastCache *FastCache

type FastCache struct {
	// Whether delivery HLS in high performance mode.
	HLSHighPerformance bool
	// Whether deliver HLS in low latency mode.
	HLSLowLatency bool
}

func NewFastCache() *FastCache {
	return &FastCache{}
}

func (v *FastCache) Refresh(ctx context.Context) error {
	pipe := rdb.Pipeline()
	hlsLowLatencyCmd := pipe.HGet(ctx, SRS_LL_HLS, "hlsLowLatency")
	noHlsCtxCmd := pipe.HGet(ctx, SRS_HP_HLS, "noHlsCtx")

	// We ignore the Exec error because we check individual command results.
	// We want to tolerate redis.Nil (and potentially other errors by defaulting to false,
	// though logging them would be better, but we stick to existing logic of defaulting to false).
	_, _ = pipe.Exec(ctx)

	if val, err := hlsLowLatencyCmd.Result(); err == nil && val == "true" {
		v.HLSLowLatency = true
	} else {
		v.HLSLowLatency = false
	}

	if val, err := noHlsCtxCmd.Result(); err == nil && val == "true" {
		v.HLSHighPerformance = true
	} else {
		v.HLSHighPerformance = false
	}

	return nil
}

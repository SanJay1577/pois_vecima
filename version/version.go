/*
 * Copyright, 2020 - 2021, Vecima Networks Inc.,
 *   as an unpublished work. This document contains confidential and
 *   proprietary information, including trade secrets, of Vecima
 *   Networks Inc.  Any use, reproduction or transmission
 *   of any part or the whole of this document is expressly prohibited
 *   without the prior written permission of Vecima Networks
 *   Inc.
 */
package version

import (
	"fmt"
)

// VersionInfo holds information about an application version.
type VersionInfo struct {
	Name     string
	FullName string
	Release  string
	Extra    string
}

func (v VersionInfo) Version() string {
	return fmt.Sprintf("%s%s", v.Release, v.Extra)
}

func (v VersionInfo) Application() string {
	return v.Name
}

func (v VersionInfo) ApplicationName() string {
	return v.FullName
}
func (v VersionInfo) UserAgent() string {
	return fmt.Sprintf("Vecima-Networks-Inc-%s/%s", v.Application(), v.Version())
}

func (v VersionInfo) String() string {
	return fmt.Sprintf("Vecima Networks Inc. © 2023 %s Version %s", v.ApplicationName(), v.Version())
}

var Version = VersionInfo{
	Name:     "pois",
	FullName: "Placement Opportunity Information Service",
	Release:  "1.0.0",
	Extra:    "",
}

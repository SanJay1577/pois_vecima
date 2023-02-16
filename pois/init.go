/*
 * Copyright, 2020, Vecima Technology Inc.,
 *   as an unpublished work. This document contains confidential and
 *   proprietary information, including trade secrets, of Vecima
 *   Technology Inc.  Any use, reproduction or transmission
 *   of any part or the whole of this document is expressly prohibited
 *   without the prior written permission of Vecima Technology
 *   Inc.
 */

package pois

import (
	alias "pois/pois/alias"
	ccms "pois/pois/ccms"
)

func Initialize() error {

	if err := alias.Initialize(); err != nil {
		return err
	}

	if err := ccms.Initialize(); err != nil {
		return err
	}

	return nil
}

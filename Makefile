# Copyright (c) 2020 Contaim, LLC
# 
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

bootstrap: # bootstrap the build by downloading additional tools that may be used by devs
	@go generate -tags tools tools/tools.go

test:
	@go test ./...
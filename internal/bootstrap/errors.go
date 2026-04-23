package bootstrap

import "errors"

var InitError = errors.New("env var BAILUP_EMAIL, BAILUP_PASS, BAILUP_REGULATION need to be set")

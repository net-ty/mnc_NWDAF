module nwdaf.com

go 1.14

replace nwdaf.com/logger => ../logger

replace nwdaf.com/service => ../service

replace nwdaf.com/factory => ../factory

replace nwdaf.com/util => ../util

replace nwdaf.com/consumer => ../consumer

replace nwdaf.com/context => ../context

replace nwdaf.com/mtlf => ../mtlf

replace nwdaf.com/anlf => ../AnLF

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.0
	github.com/free5gc/http2_util v1.0.0
	github.com/free5gc/http_wrapper v1.0.0
	github.com/free5gc/logger_conf v1.0.0
	github.com/free5gc/logger_util v1.0.0
	github.com/free5gc/openapi v1.0.0
	github.com/free5gc/path_util v1.0.0
	github.com/free5gc/version v1.0.0
	github.com/gin-gonic/gin v1.6.3
	github.com/google/go-cmp v0.5.1 // indirect
	github.com/google/uuid v1.3.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/urfave/cli v1.22.4
	golang.org/x/sys v0.0.0-20201214210602-f9fddec55a1e // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

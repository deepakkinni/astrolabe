module github.com/vmware-tanzu/astrolabe

go 1.13

replace github.com/vmware/gvddk => ./vendor/vmware.com/gvddk

require (
	github.com/aws/aws-sdk-go v1.29.19
	github.com/go-openapi/errors v0.19.4
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/runtime v0.19.12
	github.com/go-openapi/spec v0.19.7
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.8
	github.com/go-openapi/validate v0.19.7
	github.com/go-swagger/go-swagger v0.23.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/vmware/govmomi v0.22.2
	github.com/vmware/gvddk v0.8.1
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	k8s.io/klog v1.0.0
	k8s.io/utils v0.0.0-20191218082557-f07c713de883
)

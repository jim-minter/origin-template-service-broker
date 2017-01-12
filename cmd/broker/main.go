package main

import (
	"net"
	"net/http"

	"github.com/golang/glog"
	"github.com/jim-minter/origin-template-service-broker/pkg/broker/template"
	"github.com/jim-minter/origin-template-service-broker/pkg/handler"
	"github.com/openshift/origin/pkg/cmd/flagtypes"
	"github.com/openshift/origin/pkg/cmd/util/clientcmd"
	"github.com/spf13/pflag"
	"k8s.io/kubernetes/pkg/util/logs"

	_ "github.com/openshift/origin/pkg/api/install"
)

/*
	_ "k8s.io/kubernetes/pkg/api/install"
	_ "k8s.io/kubernetes/pkg/apis/autoscaling/install"
	_ "k8s.io/kubernetes/pkg/apis/batch/install"
	_ "k8s.io/kubernetes/pkg/apis/extensions/install"
*/

var addr = pflag.String("addr", ":8000", "listen address")
var createProjects = pflag.Bool("create-projects", false, "set true to create a project per service instance")

func main() {
	logs.InitLogs()
	flagtypes.GLog(pflag.CommandLine)

	factory := clientcmd.New(pflag.CommandLine)
	pflag.Parse()

	broker, err := template.NewBroker(factory, *createProjects)
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	glog.Info("running")

	err = http.Serve(l, handler.NewHandler(broker))
	if err != nil {
		panic(err)
	}
}

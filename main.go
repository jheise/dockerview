package main

import (
	// std
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	// external
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
)

var (
	client  *docker.Client
	port    string
	address string
)

type ContainerNet struct {
	IPAddr  string
	NetType string
	Port    int64
}

type Container struct {
	Name  string
	Ports []ContainerNet
}

type DockerList struct {
	Containers []*Container
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	// get list of all containers
	opts := docker.ListContainersOptions{true, false, 0, "", "", nil, nil}
	containers, err := client.ListContainers(opts)
	if err != nil {
		panic(err)
	}

	// build list of container structs
	dockerlist := new(DockerList)

	for _, container := range containers {
		// fmt.Fprintln(w, container.Names[0])
		newcontainer := new(Container)
		newcontainer.Name = container.Names[0]
		if len(container.Ports) > 0 {
			for _, port := range container.Ports {
				// fmt.Fprintf(w, "%s %s:%d\n", port.Type, port.IP, port.PublicPort)
				newnet := ContainerNet{port.IP, port.Type, port.PublicPort}
				newcontainer.Ports = append(newcontainer.Ports, newnet)

			}
		}
		dockerlist.Containers = append(dockerlist.Containers, newcontainer)
	}

	for _, x := range dockerlist.Containers {
		fmt.Printf("%s\n", x)
	}

	// load template and return html
	template_bytes, err := ioutil.ReadFile("templates/dockerview.template")
	if err != nil {
		panic(err)
	}
	template_data := string(template_bytes)

	tmpl := template.New("dockerview-template")
	tmpl, err = tmpl.Parse(template_data)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, *dockerlist)
}

func init() {
	flag.StringVar(&address, "address", "0.0.0.0", "address to bind to")
	flag.StringVar(&port, "port", "9999", "port to bind to")
	flag.Parse()

	endpoint := "unix:///var/run/docker.sock"
	localclient, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	client = localclient

}

func main() {
	socket_str := address + ":" + port

	router := mux.NewRouter()
	router.HandleFunc("/", ListHandler).Methods("GET")
	fmt.Printf("Listenting on %s\n", socket_str)
	http.ListenAndServe(socket_str, router)
}

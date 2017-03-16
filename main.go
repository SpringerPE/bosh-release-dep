package main

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/awalterschulze/gographviz/parser"
	"gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
)

var (
	releaseFolder = kingpin.Arg("release-folder", "Folder containing the bosh release").Required().String()
)

var (
	version = "0.0.1"
)

type PackageSpec struct {
	Packages []string `yaml:"dependencies"`
}

type JobSpec struct {
	Packages []string `yaml:"packages"`
}

func quote(name string) string {
	return "\"" + name + "\""
}

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	graphAst, err := parser.ParseString(`digraph Deps {}`)
	if err != nil {
		log.Fatal(err)
	}

	graph := gographviz.NewGraph()
	gographviz.Analyse(graphAst, graph)

	packagesFolder := path.Join(*releaseFolder, "packages")
	packages, err := ioutil.ReadDir(packagesFolder)
	if err != nil {
		log.Fatalf("error: %v", err)

	}

	jobsFolder := path.Join(*releaseFolder, "jobs")
	jobs, err := ioutil.ReadDir(jobsFolder)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, pkg := range packages {
		if pkg.IsDir() {
			spec := PackageSpec{}
			data, err := ioutil.ReadFile(path.Join(packagesFolder, pkg.Name(), "spec"))
			if err != nil {
				log.Fatalf("error: %v", err)

			}
			err = yaml.Unmarshal([]byte(data), &spec)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			graph.AddNode("Deps", quote(pkg.Name()), map[string]string{"type": "package"})
			for _, dep := range spec.Packages {
				graph.AddNode("Deps", quote(dep), map[string]string{"type": "package"})
				graph.AddEdge(quote(pkg.Name()), quote(dep), true, nil)
			}
		}
	}

	for _, job := range jobs {
		if job.IsDir() {
			spec := JobSpec{}
			data, err := ioutil.ReadFile(path.Join(jobsFolder, job.Name(), "spec"))
			if err != nil {
				log.Fatalf("error: %v", err)

			}
			err = yaml.Unmarshal([]byte(data), &spec)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			jobName := quote(job.Name() + "_job")
			graph.AddNode("Deps", jobName, map[string]string{"type": "job"})
			for _, pkg := range spec.Packages {
				graph.AddNode("Deps", quote(pkg), map[string]string{"type": "package"})
				graph.AddEdge(jobName, quote(pkg), true, nil)
			}
		}
	}

	output := graph.String()
	fmt.Println(output)

}

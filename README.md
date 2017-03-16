# Bosh Release Dependencies

## Description
Simple tool written in go to print a graphviz compatible dot file describing the dependencies between jobs and packages in a boshrelease

## Installation
go get github.com/SpringerPE/bosh-release-dep

## Usage
```
bosh-release-dep path-to-the-boshrelease-folder > dependencies.dot
```
In order to visualise the graph install graphviz and render the image:

```
bosh-release-dep path-to-the-boshrelease-folder | dot -Tpng -o dep.png
```

Enjoy!!!

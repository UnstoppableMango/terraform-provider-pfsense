package main

import (
	"charm.land/log/v2"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func Patch(root string) error {
	pkgs, err := loader.LoadRoots(root)
	if err != nil {
		return err
	}

	reg := &markers.Registry{}
	if err = Register(reg); err != nil {
		return err
	}

	for _, pkg := range pkgs {
		log.Info("Patching package", "name", pkg.Name)
		if err := patchPackage(reg, pkg); err != nil {
			return err
		}
	}

	return nil
}

func patchPackage(reg *markers.Registry, pkg *loader.Package) error {
	rt, err := genall.FromOptions(reg, []string{})
	if err != nil {
		return err
	}

	rt.Run()

	return nil
}

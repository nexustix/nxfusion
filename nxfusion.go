package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	bp "github.com/nexustix/boilerplate"
	nrc "github.com/nexustix/nxReplicatorCommon"
)

// nxfusion add amazingMolecule amazingAtom

func main() {
	//fmt.Println("Hello World")
	args := os.Args

	usr, err := user.Current()
	bp.FailError(err)
	workingDir := usr.HomeDir
	atomDir := nrc.InitWorkFolder(workingDir, ".nxreplicator", "atoms")
	moleculeDir := nrc.InitWorkFolder(workingDir, ".nxreplicator", "molecules")

	atomManager := nrc.AtomManager{WorkingDir: atomDir}
	//molecule := nrc.Molecule{}

	action := bp.StringAtIndex(1, args)
	packName := bp.StringAtIndex(2, args)
	providerID := bp.StringAtIndex(3, args)
	atomID := bp.StringAtIndex(4, args)

	fmt.Printf("%s | %s | %s | %s\n", action, packName, providerID, atomID)

	switch action {
	case "add":
		addAtom(&atomManager, moleculeDir, packName, providerID, atomID)
	case "remove":
		removeAtom(moleculeDir, packName, providerID, atomID)
	case "list":
		listAtoms(moleculeDir, packName)
	}
}

//addAtom adds atom to molecule
func addAtom(atomManager *nrc.AtomManager, moleculeDir, packName, providerID, atomID string) {
	filePath := path.Join(moleculeDir, packName+".nxrm")
	if atomManager.HasEntry(providerID, atomID) {
		tmpAtom := atomManager.GetEntry(providerID, atomID)

		tmpMolecule := nrc.Molecule{}
		tmpMolecule.LoadFromFile(filePath)

		tmpMoleculeItem := nrc.MoleculeItem{Explicit: true, ProviderID: tmpAtom.Provider, AtomID: tmpAtom.ID, Dir: tmpAtom.RelativePath}
		tmpMolecule.AddItem(tmpMoleculeItem)

		for _, v := range tmpAtom.Dependencies {
			var tmpDepItem nrc.MoleculeItem
			if atomManager.HasEntry(tmpAtom.Provider, v) {
				tmpDependency := atomManager.GetEntry(tmpAtom.Provider, v)
				tmpDepItem = nrc.MoleculeItem{Explicit: false, ProviderID: tmpDependency.Provider, AtomID: tmpDependency.ID, Dir: tmpDependency.RelativePath}
			} else {
				fmt.Printf("<!> WARNING no source for '%s' found (adding wildcard entry)\n", v)
				tmpDepItem = nrc.MoleculeItem{Explicit: false, ProviderID: "", AtomID: v, Dir: "mods"}
			}
			tmpMolecule.AddItem(tmpDepItem)
		}

		tmpMolecule.SaveToFile(filePath)
	} else {
		fmt.Printf("<!> ERROR Atom >%s< not found\n", atomID)
	}
}

//removeAtom removes atom from molecule
func removeAtom(moleculeDir, packName, providerID, atomID string) {
	filePath := path.Join(moleculeDir, packName+".nxrm")

	tmpMolecule := nrc.Molecule{}
	tmpMolecule.LoadFromFile(filePath)

	tmpMolecule.RemoveItem(providerID, atomID)

	tmpMolecule.SaveToFile(filePath)
}

//listAtoms list atoms of molecule
func listAtoms(moleculeDir, packName string) {
	filePath := path.Join(moleculeDir, packName+".nxrm")

	tmpMolecule := nrc.Molecule{}
	tmpMolecule.LoadFromFile(filePath)
	fmt.Printf("-----\n")
	for k, v := range tmpMolecule.MoleculeItems {
		fmt.Printf("(%v) <%v> <%s> %s\n", k, v.Explicit, v.ProviderID, v.AtomID)
	}
	fmt.Printf("-----\n")
}

//func getAtom(){}

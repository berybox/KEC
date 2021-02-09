package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/urfave/cli"

	KEC "../core"
)

func exclude(c *cli.Context) error {
	if target == "" || nonTarget == "" || output == "" {
		cli.ShowCommandHelp(c, "exclude")
		err := fmt.Errorf("\n   Mandatory parameter(s) (target, nontarget, output) not set")
		showMsg("")
		showError(err, true)
		return err
	}

	if fileInfo(target) == 0 || fileInfo(nonTarget) == 0 {
		err := fmt.Errorf("Target or/and nontarget sequences not found")
		showError(err, true)
		return err
	}

	//showMsg("Exclude parameters:\n\nNontarget: %s\nTarget: %s\nOutput: %s\nK: %d\nMin: %d\nMax: %d", nonTarget, target, output, k, min, max)

	kecTime := time.Now()
	showMsg("\nStarting KEC - searching unique sequences\n\n")
	kec := KEC.New(k, min, max, false, reverse)
	var actionTime time.Time
	var actionMem runtime.MemStats

	//Add nontarget
	for _, fn := range fileList(nonTarget) {
		actionTime = time.Now()
		showMsg("Reading NONTARGET file: \"%s\" ", filepath.Base(fn))
		kec.AddNontargetFasta(fn)
		showMsg("took %s\n", time.Since(actionTime))

		//Show memory stats
		/*
			runtime.GC()
			runtime.ReadMemStats(&actionMem)
			showMsg("-- HeapAlloc: %d MB,  Mallocs: %d\n", actionMem.HeapAlloc/1024/1024, actionMem.Mallocs)
			//*/
	}

	//Add target
	for _, fn := range fileList(target) {
		actionTime = time.Now()
		showMsg("Reading TARGET file: \"%s\" ", filepath.Base(fn))
		kec.AddTargetFastaCR(fn)
		showMsg("took %s\n", time.Since(actionTime))
	}

	//*
	runtime.GC()
	runtime.ReadMemStats(&actionMem)
	showMsg("-- HeapAlloc: %d MB,  Mallocs: %d\n", actionMem.HeapAlloc/1024/1024, actionMem.Mallocs)
	//*/

	//Merge
	actionTime = time.Now()
	showMsg("Merging K-mers (%d nt) to sequences longer than %d ", k, min)
	kec.MergeCrossRef()
	showMsg("took %s\n", time.Since(actionTime))

	fmt.Println("")

	numSeq := kec.NumCrossRefSeq()

	//Save results
	if numSeq > 0 {
		actionTime = time.Now()
		showMsg("Saving results to: %s ", filepath.Base(output))
		kec.SaveCrossRefFasta(output)
		showMsg("took %s\n", time.Since(actionTime))
	}

	showMsg("Whole search took %s\n", time.Since(kecTime))
	showMsg("Found %d unique sequences\n", numSeq)

	return nil
}

func include(c *cli.Context) error {
	if master == "" || pool == "" || output == "" {
		cli.ShowCommandHelp(c, "include")
		err := fmt.Errorf("\n   Mandatory parameter(s) (master, pool, output) not set")
		showMsg("")
		showError(err, true)
		return err
	}

	if fileInfo(master) == 0 || fileInfo(pool) == 0 {
		err := fmt.Errorf("Master or/and pool sequences not found")
		showError(err, true)
		return err
	}

	//showMsg("Include parameters:\n\nMaster: %s\nPool: %s\nOutput: %s\nK: %d\nMin: %d\nMax: %d", master, pool, output, k, min, max)
	kecTime := time.Now()
	showMsg("\nStarting KEC - searching consensus sequences\n\n")
	kec := KEC.New(k, min, max, true, false)
	var actionTime time.Time

	//Add pool
	for _, fn := range fileList(pool) {
		actionTime = time.Now()
		showMsg("Reading POOL file: \"%s\" ", filepath.Base(fn))
		kec.AddNontargetFasta(fn)
		showMsg("took %s\n", time.Since(actionTime))
	}

	//Add master
	for _, fn := range fileList(master) {
		actionTime = time.Now()
		showMsg("Reading MASTER file: \"%s\" ", filepath.Base(fn))
		kec.AddTargetFastaCR(fn)
		showMsg("took %s\n", time.Since(actionTime))
	}

	//Merge
	actionTime = time.Now()
	showMsg("Merging K-mers (%d nt) to sequences longer than %d ", k, min)
	kec.MergeCrossRef()
	showMsg("took %s\n", time.Since(actionTime))

	fmt.Println("")

	numSeq := kec.NumCrossRefSeq()

	//Save results
	if numSeq > 0 {
		actionTime = time.Now()
		showMsg("Saving results to: %s ", filepath.Base(output))
		kec.SaveCrossRefFasta(output)
		showMsg("took %s\n", time.Since(actionTime))
	}

	showMsg("Whole search took %s\n", time.Since(kecTime))
	showMsg("Found %d consensus sequences\n", numSeq)

	return nil
}

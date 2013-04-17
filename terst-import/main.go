/*
Command terst-import will smuggol "github.com/robertkrimen/terst" into a local Go package/repository

Install

    go get github.com/robertkrimen/terst-import

Usage

    Usage: terst-import [target]                                               
     -quiet=false: Be absolutely quiet                                         
     -update=false: Update (go get -u) package first                           
     -verbose=false: Be more verbose                                           
                                                                               
        # Import "github.com/robertkrimen/terst" into the current directory    
        $ terst-import                                                         
                                                                               
        # Import "github.com/robertkrimen/terst" into another directory        
        $ terst-import ./xyzzy                                                 
*/
package main

import (
	"github.com/robertkrimen/smuggol"
	_ "github.com/robertkrimen/terst"
)

func main() {
	smuggol.Main("terst-import", "github.com/robertkrimen/terst", nil)
}

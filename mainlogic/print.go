// Copyright 2016 Politecnico di Torino
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mainlogic

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func PrintL2Switch(name string) {
	if sw, ok := switches[name]; ok {
		//TODO Improve. First implementation
		spew.Dump(sw)
	}
}

func PrintL2Switches() {
	for swname, _ := range switches {
		PrintL2Switch(swname)
	}
}

func PrintRouter(name string) {
	if router, ok := routers[name]; ok {
		//TODO Improve. First implementation
		spew.Dump(router)
	}
}

func PrintRouters() {
	for routername, _ := range routers {
		PrintRouter(routername)
	}
}

func PrintMainLogic() {
	fmt.Printf("\nSwitches\n\n")
	PrintL2Switches()
	fmt.Printf("\nRouters\n\n")
	PrintRouters()
}
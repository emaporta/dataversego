# dataversego
A Go library for interacting with the [Dataverse API](https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/overview).

## Installation
To install dataversego, use go get:

``` golang
go get github.com/emaporta/dataversego
```

## Usage
Import the package into your Go code:

``` golang
import "github.com/emaporta/dataversego"
```
To use the library, you will need a Dataverse application user. You can create one by following the instructions [here](https://learn.microsoft.com/en-us/power-platform/admin/manage-application-users).

Here is an example of how to use the dataversego library to get a Dataverse row:

``` golang
package main

import (
	"fmt"

	"github.com/emaporta/dataversego"
)

func main() {
    // Use your own CLIENTID, SECRET, TOKEN and ORGURL
	auth := dataversego.Authenticate("CLIENTID", "SECRET", "TOKEN", "ORGURL")

	fmt.Println(auth)

	columns := []string{"fullname"}

	retrieveParameters := dataversego.RetrieveSignature{
		Auth:      auth,
		TableName: "contacts",
		Id:        "CONTACT_GUID",
		Columns:   columns,
	}

	ent, err := dataversego.Retrieve(retrieveParameters)

	if err != nil {

		fmt.Println(err)
	}
	fmt.Println(ent)
}
```
## Documentation
For complete documentation of the library's functions and types, see the [GoDoc](https://godoc.org/github.com/emaporta/dataversego) page.

## Contributing
If you find a bug or have a feature please raise an issue
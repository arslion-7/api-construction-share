package main

import (
	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/routes"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
	initializers.SyncDB()
	// initializers.InitDataFill()
	// initializers.CreateDirs()
}

func main() {
	r := routes.SetupRoutes()

	r.Run(":8081")
}

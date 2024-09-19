//go:build windows

package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"log"
	"os"
	"vanbandmate/controllers"
	"vanbandmate/initializers"
	"vanbandmate/routes"
)

// Default name value for the service
var (
	svcName = "vbmSrv"
	svcDsc  = "Web service VanBandMate"
)

//func startWeb() {
//	//changeCurrentDirectory()
//	initializers.LoadEnvVariables()
//	initializers.ConnectToDatabase()
//
//	if os.Getenv("MIGRATE_DB") == "true" {
//		initializers.MigrateDB()
//	}
//	if os.Getenv("GEN_DATA_DB") == "true" {
//		initializers.GenData()
//	}
//	engine := html.New("./views", ".html")
//	engine.AddFuncMap(fiber.Map{
//		//"GetCodeUser":           controllers.GetCodeUser,
//		"IsChecked":         controllers.IsChecked,
//		"IsVerify":          controllers.IsVerify,
//		"GetUserID":         controllers.GetUserID,
//		"FormatDate":        controllers.FormatDate,
//		"FormatTime":        controllers.FormatTime,
//		"FormatTimeComment": controllers.FormatTimeComment,
//		"FormatPrice":       controllers.FormatPrice,
//		"IsTimeAfterNow":    controllers.IsTimeAfterNow,
//		"IsSelected":        controllers.IsSelected,
//		"CheckPermission":   controllers.CheckPermission,
//	})
//	app := fiber.New(fiber.Config{
//		Views:             engine,
//		PassLocalsToViews: true,
//		BodyLimit:         1024 * 1024 * 100,
//	})
//
//	app.Static("/", "./public")
//	routes.RouteInit(app)
//
//	//go func() {
//	if os.Getenv("HOST_WEB") == "https" {
//		app.ListenTLS(":"+os.Getenv("PORT"), os.Getenv("PEM"), os.Getenv("KEY"))
//
//	} else {
//		app.Listen(":" + os.Getenv("PORT"))
//
//	}
//	//}()
//
//}

func startWeb() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()

	if os.Getenv("MIGRATE_DB") == "true" {
		initializers.MigrateDB()
	}
	if os.Getenv("GEN_DATA_DB") == "true" {
		initializers.GenData()
	}

	engine := html.New("./views", ".html")
	engine.AddFuncMap(fiber.Map{
		"IsChecked":         controllers.IsChecked,
		"IsVerify":          controllers.IsVerify,
		"GetUserID":         controllers.GetUserID,
		"FormatDate":        controllers.FormatDate,
		"FormatTime":        controllers.FormatTime,
		"FormatTimeComment": controllers.FormatTimeComment,
		"FormatPrice":       controllers.FormatPrice,
		"IsTimeAfterNow":    controllers.IsTimeAfterNow,
		"IsSelected":        controllers.IsSelected,
		"CheckPermission":   controllers.CheckPermission,
	})

	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
		BodyLimit:         1024 * 1024 * 100,
	})

	app.Static("/", "./public")
	routes.RouteInit(app)

	// Lắng nghe trên cổng được cung cấp bởi Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Cổng mặc định khi chạy cục bộ
	}

	app.Listen(":" + port)
}

func main() {
	startWeb()
}

func usage(errmsg string) {
	_, err := fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	if err != nil {
		log.Println("Error printing usage message to the 'Stderr': ", err)
	}

	os.Exit(2)
}

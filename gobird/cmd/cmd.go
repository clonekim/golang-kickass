package cmd

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"net/http"
)

var Version string = "0.0.1"
var globalContext context.Context = context.Background()

type GoBird struct {
	echo.Context
	DriveServ *drive.Service
	SheetServ *sheets.Service
	ServiceAccountEmail *string
	Debug bool
	Port int
}


func RunCmd(bird *GoBird) *cli.App {

	return &cli.App{
		Name:  "gobird " + Version,
		Usage: "google drive app implemented by Go.",
		Flags: [] cli.Flag {

			&cli.BoolFlag{
				Name: "test",
				Value: false,
				Usage: "Test mode",
			},

			&cli.BoolFlag{
				Name: "debug",
				Value: false,
				Usage: "debug option",
			},

			&cli.IntFlag{
				Name: "port",
				Value: 8000,
				Usage: "port number",
			},

			&cli.StringFlag{
				Name: "credential",
				Value: "",
				Usage: "credential file to google auth",
			},

		},

		Action: func(c *cli.Context) error {

			port := c.Int("port")
			debug := c.Bool("debug")
			test := c.Bool("test")

			if !test {

				if c.String("credential")  == "" {
					return errors.New("credential file is required")
				}

				b , err := ioutil.ReadFile(c.String("credential")); if err != nil {
					return err
				}

				conf, err := createJWTConfig(b); if err != nil {
					return err
				}

				client := conf.Client(globalContext)
				drv, err := createDriveService(client); if err != nil {
					return err
				}

				sheet, err := createSheetService(client); if err != nil {
					return err
				}

				bird.ServiceAccountEmail = &conf.Email
				bird.DriveServ = drv
				bird.SheetServ = sheet

			}

			bird.Port = port
			bird.Debug =  debug

			return nil

		},
	}

}


func createJWTConfig(data []byte) (*jwt.Config, error) {
	log.Println("creating HttpClient")
	//*jwt.Config
	conf, err := google.JWTConfigFromJSON(data,
		drive.DriveScope,
		drive.DriveFileScope,
		drive.DriveAppdataScope,
		drive.DriveScriptsScope,
		drive.DriveMetadataScope,
		drive.DriveMetadataReadonlyScope,
		drive.DrivePhotosReadonlyScope,
		sheets.SpreadsheetsScope,
		sheets.SpreadsheetsReadonlyScope,
	)

	if err != nil {
		return nil, err
	}

	return conf, err
}



func createDriveService(client *http.Client) (*drive.Service, error) {
	log.Println("creating drive service")
	return drive.NewService(globalContext, option.WithHTTPClient(client))
}

func createSheetService(client *http.Client) (*sheets.Service, error) {
	log.Println("creating sheet service")
	return sheets.NewService(globalContext, option.WithHTTPClient(client))
}


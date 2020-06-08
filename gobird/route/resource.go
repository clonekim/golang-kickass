package route

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"gobird/cmd"
	"google.golang.org/api/drive/v3"
)

/*
  모든 파일 가져오기
  폴더, 파일이든 모드 파일로 취급
 */
func GetResources(c echo.Context) error {
	bird := c.(*cmd.GoBird)

	results, err := bird.DriveServ.Files.List().SupportsAllDrives(true).
		Fields("files(id, kind, mimeType, name, description, starred, trashed, parents, owners, webContentLink, webViewLink, iconLink, spaces, thumbnailLink, createdTime, modifiedTime, sharingUser, size, capabilities, permissions)").
		Do()

	if err != nil {
		return c.JSON(500, err)
	} else {
		return c.JSON(200, results)
	}

}

func GetResource(c echo.Context) error {
	id := c.Param("id")
	bird := c.(*cmd.GoBird)
	result, err := bird.DriveServ.Files.Get(id).Do()


	if err != nil {
		return c.JSON(500, err)
	} else {
		return c.JSON(200, result)
	}
}

type Resource struct {
	Name     string `json:"name" validate:"required"`
	MimeType string `json:"mimeType" validate:"required"`
	Parents []string `json:"parents,omitempty"`
	Email string `json:"email,omitempty" validate:"email"`
}



func CreateNewResource( c echo.Context) (err error) {

	r := new(Resource)

	err = c.Bind(r); if err != nil {
		return
	}

	err = c.Validate(r); if err != nil {
		return c.JSON(400, err)
	}

	bird := c.(*cmd.GoBird)
	file := &drive.File{Name: r.Name, MimeType: r.MimeType }

	if len(r.Parents) > 0 {
		file.Parents = r.Parents
	}

	
	result, err := bird.DriveServ.Files.Create(file).Do()

	if err != nil {
		return c.JSON(500, err)
	} else {
		if r.Email != "" {
			bird.DriveServ.Permissions.Create(result.Id, &drive.Permission{EmailAddress: r.Email, Role: "user"})
		}

		return c.JSON(200, result)
	}

}

func DeleteResource( c echo.Context) error {
	id := c.Param("id")
	bird := c.(*cmd.GoBird)
	err := bird.DriveServ.Files.Delete(id).Do()

	if err != nil {
		return c.JSON(500, err)
	} else {
		return c.NoContent(200)
	}
}


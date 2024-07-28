package cmd

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/sevensolutions/tiny-repo/core"
	"github.com/spf13/cobra"
)

var address string

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull an artifact",
	Long:  `Pull an artifact`,
	Run: func(cmd *cobra.Command, args []string) {
		println("Address " + address)
		if address == "" {
			panic("missing address")
		}

		rawSpec := args[0]

		spec, err := core.ParseVersionSpec(rawSpec)
		if err != nil {
			panic(err)
		}

		println(spec.Name)

		fullUrl := address + "/" + spec.Namespace + "/" + spec.Name + "/" + spec.Version.String()

		println(fullUrl)

		err = downloadFile("test.jpg", fullUrl)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	pullCmd.PersistentFlags().StringVar(&address, "address", "", "The TinyServer address")

	rootCmd.AddCommand(pullCmd)
}

func downloadFile(filepath string, url string) (err error) {

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJwcmVmaXgiOiJicnMifQ.dWhwiYMDZ33-yKUI-VhVu6yWj99UU-2r0GgQtHg3U8U")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition == "" {
		return errors.New("missing Content-Disposition header")
	}

	_, params, err := mime.ParseMediaType(contentDisposition)
	filename := params["filename"]

	if filename == "" {
		return errors.New("missing filename in Content-Disposition header")
	}

	println(filename)

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"bytes"
	"fmt"
	"gotDots/utils"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"strings"
)

func pushPackage(packName string) {
	foundArchives := findPackageArchive(packName)

	// TODO: Push newer version
	var packArchive string
	if len(foundArchives) == 0 {
		fmt.Printf("Could not find package with name: %s\n", packName)
		os.Exit(1)
	} else if len(foundArchives) == 1 {
		packArchive = foundArchives[0]
	} else {
		fmt.Printf("Following packages found with name: %s\n", packName)
		utils.ListNames("   ", foundArchives)
		fmt.Print("Choose by entering number: ")
		var choice int
		_, scanErr := fmt.Scanln(&choice)
		if scanErr != nil {
			fmt.Println("Could not parse input")
			os.Exit(1)
		}

		packArchive = foundArchives[choice-1]
	}

	if packArchive == "" {
		fmt.Printf("Could not find package with name %s\n", packName)
		return
	}

	token := readToken()
	backendUrl := os.Getenv("CREATE_PACKAGE_URL")

	manifest := utils.ReadManifestFromTar(packArchive)

	var client *http.Client
	var remoteURL string
	{
		//setup a mocked http client.
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := httputil.DumpRequest(r, true)
			if err != nil {
				handleError(err, true)
			}
			fmt.Printf("%s", b)
		}))
		defer ts.Close()
		client = ts.Client()
		remoteURL = backendUrl
	}

	// appsList := fmt.Sprintf("[%s]", strings.Join(returnList(manifest.IncludedApps), ","))
	appsList := encodeIncludedApps(manifest.IncludedApps)

	values := map[string]io.Reader{
		"Id":               strings.NewReader(manifest.Id),
		"Name":             strings.NewReader(manifest.Name),
		"Version":          strings.NewReader(manifest.Version),
		"IncludedAppsJson": strings.NewReader(appsList),
		"PackageArchive":   mustOpen(packArchive),
		"Visibility":       strings.NewReader(manifest.Visibility),
	}

	uploadErr := Upload(client, remoteURL, values, token)
	if uploadErr != nil {
		fmt.Printf("Could not push to registry. Error: %s\n", uploadErr.Error())
		os.Exit(1)
	}

	fmt.Printf("Package %s successfully pushed to registry\n", packName)
}

func Upload(client *http.Client, url string, values map[string]io.Reader, token string) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	//goland:noinspection GoUnhandledErrorResult
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", token)

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == 460 {
			//goland:noinspection GoErrorStringFormat
			err = fmt.Errorf("Package name exists. " +
				"If you want to push new version, use 'dots update <pack name>' instead")
		} else if res.StatusCode == http.StatusUnauthorized {
			//goland:noinspection GoErrorStringFormat
			err = fmt.Errorf("Please login and try again")
		} else {
			err = fmt.Errorf("Could not push to registry Error: %s\n", res.Status)
		}
	}
	return err
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

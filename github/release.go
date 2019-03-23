package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
)

type Release struct {
	URL             string `json:"url"`
	AssetsURL       string `json:"assets_url"`
	UploadURL       string `json:"upload_url"`
	HTMLURL         string `json:"html_url"`
	ID              int    `json:"id"`
	NodeID          string `json:"node_id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Author          struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		URL      string      `json:"url"`
		ID       int         `json:"id"`
		NodeID   string      `json:"node_id"`
		Name     string      `json:"name"`
		Label    interface{} `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}

type Releases []Release

func GetReleases(owner string, repo string) (Releases, error) {

	github := viper.GetString("github")
	resp, err := http.Get(fmt.Sprintf("%s/repos/%s/%s/releases", github, owner, repo))
	if err != nil {
		// handle error
		fmt.Printf("http.Get failed: %v", err)
	}

	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Printf("ioutil.ReadAll failed: %v", err)
	}

	//fmt.Println(string(body))

	res := Releases{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		// handle error
		fmt.Printf("json.Unmarshal failed: %v", err)
	}

	fmt.Printf("Return StatusCode: %v\n", resp.StatusCode)

	//fmt.Printf("Releases: %v", res)

	for _, r := range res {
		fmt.Printf("%4d    %s\n", r.ID, r.TagName)
	}

	return res, nil
}

func GetRelease(owner string, repo string, releaseId string) {

	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/%s", github, owner, repo, releaseId)

	resp, err := http.Get(url)
	if err != nil {
		// handle error
		fmt.Printf("http.Get failed: %v", err)
	}

	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Printf("ioutil.ReadAll failed: %v", err)
	}

	//fmt.Println(string(body))

	res := Release{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		// handle error
		fmt.Printf("json.Unmarshal failed: %v", err)
	}

	fmt.Printf("Release: %v", res)
}

func GetReleaseByTag(owner string, repo string, tag string) {

	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s", github, owner, repo, tag)

	resp, err := http.Get(url)
	if err != nil {
		// handle error
		fmt.Printf("http.Get failed: %v", err)
	}

	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Printf("ioutil.ReadAll failed: %v", err)
	}

	//fmt.Println(string(body))

	res := Release{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		// handle error
		fmt.Printf("json.Unmarshal failed: %v", err)
	}

	fmt.Printf("Release: %v", res)
}

type RequestCreateRelease struct {
	TagName         string `json:"tag_name"`         // Required. The name of the tag.
	TargetCommitish string `json:"target_commitish"` // Specifies the commitish value that determines where the Git tag is created from. Can be any branch or commit SHA. Unused if the Git tag already exists. Default: the repository's default branch (usually master).
	Name            string `json:"name"`             // The name of the release.
	Body            string `json:"body"`             // Text describing the contents of the tag.
	Draft           bool   `json:"draft"`            // true to create a draft (unpublished) release, false to create a published one. Default: false
	Prerelease      bool   `json:"prerelease"`       // true to identify the release as a prerelease. false to identify the release as a full release. Default: false
}

func CreateRelease(owner string, repo string) {

	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases", github, owner, repo)
	token := viper.GetString("token")

	request := RequestCreateRelease{}
	request.TagName = viper.GetString("tag_name")
	request.TargetCommitish = viper.GetString("target_commitish")
	request.Name = viper.GetString("name")
	request.Body = viper.GetString("body")
	request.Draft = viper.GetBool("draft")
	request.Prerelease = viper.GetBool("prerelease")

	requestByte, _ := json.Marshal(request)

	method := http.MethodPost
	if viper.GetBool("edit") {
		method = http.MethodPatch
		url += fmt.Sprintf("/%s", viper.GetString("id"))
	}

	fmt.Printf("create release: %v\n", string(requestByte))
	fmt.Printf("create release: %v\n", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestByte))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(token, "x-oauth-basic")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Printf("http.Client.Do failed: %v", err)
		return
	}
	//noinspection GoUnhandledErrorResult
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	statusCode := resp.StatusCode
	header := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	fmt.Println(statusCode)
	fmt.Println(header)

	if statusCode == http.StatusNotFound {
		var result map[string]interface{}
		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			fmt.Println("Unmarshal failed, ", err)
			return
		}

		fmt.Println(result["message"])
		fmt.Println(result["documentation_url"])
	}
}

func DeleteRelease(owner string, repo string) {

	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/%s", github, owner, repo, viper.GetString("id"))
	token := viper.GetString("token")

	method := http.MethodDelete

	fmt.Printf("delete release: %v\n", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(token, "x-oauth-basic")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Printf("http.Client.Do failed: %v", err)
		return
	}
	//noinspection GoUnhandledErrorResult
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	statusCode := resp.StatusCode
	header := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	fmt.Println(statusCode)
	fmt.Println(header)

	if statusCode == http.StatusNotFound {
		var result map[string]interface{}
		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			fmt.Println("Unmarshal failed, ", err)
			return
		}

		fmt.Println(result["message"])
		fmt.Println(result["documentation_url"])
	}
}

func ListAssets(owner string, repo string) {

	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/%s/assets", github, owner, repo, viper.GetString("id"))
	token := viper.GetString("token")

	method := http.MethodGet

	fmt.Printf("List assets for a release: %v\n", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(token, "x-oauth-basic")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Printf("http.Client.Do failed: %v", err)
		return
	}
	//noinspection GoUnhandledErrorResult
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	statusCode := resp.StatusCode
	header := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	fmt.Println(statusCode)
	fmt.Println(header)

	if statusCode == http.StatusNotFound {
		var result map[string]interface{}
		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			fmt.Println("Unmarshal failed, ", err)
			return
		}

		fmt.Println(result["message"])
		fmt.Println(result["documentation_url"])
	}
}

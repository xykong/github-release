package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
)

type Release struct {
	Url             string `json:"url"`
	AssetsUrl       string `json:"assets_url"`
	UploadUrl       string `json:"upload_url"`
	HtmlUrl         string `json:"html_url"`
	Id              int    `json:"id"`
	NodeId          string `json:"node_id"`
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Author          struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []struct {
		Url      string      `json:"url"`
		Id       int         `json:"id"`
		NodeId   string      `json:"node_id"`
		Name     string      `json:"name"`
		Label    interface{} `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadUrl string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballUrl string `json:"tarball_url"`
	ZipballUrl string `json:"zipball_url"`
	Body       string `json:"body"`
}

type Releases []Release

func SendRequest(url string, method string, body []byte, token string, mime string, v interface{}) (*http.Response, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("http.NewRequest failed: %v", err)
		return nil, err
	}

	if mime == "" {
		mime = "application/json"
	}
	req.Header.Set("Content-Type", mime)

	if token != "" {
		req.SetBasicAuth(token, "x-oauth-basic")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("http.Client.Do failed: %v", err)
		return nil, err
	}

	//noinspection GoUnhandledErrorResult
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("http.Client.Do failed: %v", err)
		return resp, err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Printf("json.Unmarshal failed: %v", err)
		return resp, err
	}

	return resp, nil
}

func validate(input map[string]string) error {

	for k, v := range input {

		if v == "" {
			return fmt.Errorf("%s is required", k)
		}
	}

	return nil
}

// ListReleases
func ListReleases(owner string, repo string) (Releases, error) {

	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases", github, owner, repo)

	err := validate(map[string]string{
		"user": owner,
		"repo": repo,
	})
	if err != nil {
		return nil, fmt.Errorf("list releases for a repository: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Infof("List releases for a repository")

	var releases = Releases{}
	_, err = SendRequest(url, http.MethodGet, nil, "", "", &releases)
	if err != nil {
		return nil, err
	}

	for _, r := range releases {
		fmt.Printf("%4d    %s\n", r.Id, r.TagName)
	}

	return releases, nil
}

func GetRelease(owner string, repo string, releaseId string) (*Release, error) {

	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/%s", github, owner, repo, releaseId)

	err := validate(map[string]string{
		"user":       owner,
		"repo":       repo,
		"release_id": releaseId,
	})
	if err != nil {
		return nil, fmt.Errorf("get a single release: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Infof("Get a single release")

	var release = Release{}
	_, err = SendRequest(url, http.MethodGet, nil, "", "", &release)
	if err != nil {
		return nil, err
	}

	result, err := json.MarshalIndent(release, "", "\t")
	fmt.Printf("get a single release:\n%s\n", string(result))

	return &release, nil
}

func GetReleaseByTag(owner string, repo string, tag string) (*Release, error) {

	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s", github, owner, repo, tag)

	err := validate(map[string]string{
		"user": owner,
		"repo": repo,
		"tag":  tag,
	})
	if err != nil {
		return nil, fmt.Errorf("get a single release: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Infof("Get a single release")

	var release = Release{}
	_, err = SendRequest(url, http.MethodGet, nil, "", "", &release)
	if err != nil {
		return nil, err
	}

	result, err := json.MarshalIndent(release, "", "\t")
	fmt.Printf("get a single release:\n%s\n", string(result))

	return &release, nil
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

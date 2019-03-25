package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
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

	logrus.WithFields(logrus.Fields{
		"StatusCode": resp.StatusCode,
		"Header":     resp.Header,
		"Body":       string(data),
	}).Debug("send request")

	if err != nil {
		fmt.Printf("ioutil.ReadAll failed: %v", err)
		return resp, err
	}

	if len(data) > 0 && v != nil {
		err = json.Unmarshal(data, v)
		if err != nil {
			fmt.Printf("json.Unmarshal failed: %v", err)
			return resp, err
		}
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
	token := viper.GetString("token")

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
	_, err = SendRequest(url, http.MethodGet, nil, token, "", &releases)
	if err != nil {
		return nil, err
	}

	color.Green("%20s    %10s    %5s    %s\n", "created", "id", "draft", "tag")
	for _, r := range releases {
		fmt.Printf("%20v    %10d    %5v    %s\n",
			r.CreatedAt.Format("2006-01-02 15:04:05"), r.Id, r.Draft, r.TagName)
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

func CreateRelease(owner string, repo string) error {

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

	err := validate(map[string]string{
		"user":  owner,
		"repo":  repo,
		"token": token,
	})
	if err != nil {
		return fmt.Errorf("create a release: %v", err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Infof("Create a release")

	var result map[string]interface{}
	resp, err := SendRequest(url, method, requestByte, token, "", &result)
	if err != nil {
		return err
	}

	content, err := json.MarshalIndent(result, "", "\t")
	fmt.Printf("create a %v:\n%s\n", color.GreenString("release"), string(content))

	logrus.WithFields(logrus.Fields{
		"StatusCode": resp.StatusCode,
		"Header":     resp.Header,
		"Body":       string(content),
	}).Debug("create a release")

	if resp.StatusCode == http.StatusCreated {

		logrus.WithFields(logrus.Fields{
			"id":       int64(result["id"].(float64)),
			"tag_name": result["tag_name"],
			"url":      result["url"],
		}).Info("create a release success")

		return nil
	}

	if val, ok := result["errors"]; ok {

		if errors, ok := val.([]interface{}); ok {

			logrus.Errorf("create a release failed with %d errors:", len(errors))

			for i, v := range errors {

				if item, ok := v.(map[string]interface{}); ok {
					if _, ok := item["message"]; ok {
						logrus.Infof("%v: %v", color.GreenString("%v", i), item["message"])
					} else {
						logrus.Infof("%v: %v %v", color.GreenString("%v", i), item["code"], item["field"])
					}
				}
			}
		}
	}

	return nil
}

func DeleteRelease(owner string, repo string) error {

	desc := "delete a release"
	github := viper.GetString("github")
	id := viper.GetString("id")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/%s", github, owner, repo, id)
	token := viper.GetString("token")

	method := http.MethodDelete

	err := validate(map[string]string{
		"user":       owner,
		"repo":       repo,
		"release_id": id,
		"token":      token,
	})
	if err != nil {
		return fmt.Errorf("%s: %v", desc, err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Infof(desc)

	var result map[string]interface{}
	resp, err := SendRequest(url, method, nil, token, "", &result)
	if err != nil {
		return err
	}

	//content, err := json.MarshalIndent(result, "", "\t")
	content, err := json.Marshal(result)
	logrus.WithFields(logrus.Fields{
		"StatusCode": resp.StatusCode,
		"Header":     resp.Header,
		"Body":       string(content),
	}).Debug(desc)

	if resp.StatusCode == http.StatusNoContent {

		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Infof("%s success", desc)

		return nil
	}

	if resp.StatusCode == http.StatusNotFound {

		logrus.WithFields(logrus.Fields{
			"documentation_url": result["documentation_url"],
			"message":           result["message"],
		}).Errorf("%s failed", desc)
	}

	return nil
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

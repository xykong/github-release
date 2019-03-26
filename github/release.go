package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gabriel-vasile/mimetype"
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

	//logrus.WithFields(logrus.Fields{
	//	"StatusCode": resp.StatusCode,
	//	"Header":     resp.Header,
	//	"Body":       string(data),
	//}).Debug("send request")

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		color.Green("StatusCode: %v\n", resp.StatusCode)
		color.Green("Header: %v\n", resp.Header)
		color.Green("Body: %s\n", string(data))
	}

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

	desc := "list releases for a repository"
	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases", github, owner, repo)
	token := viper.GetString("token")

	err := validate(map[string]string{
		"user": owner,
		"repo": repo,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %v", desc, err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Info(desc)

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

func ListAssets(owner string, repo string) error {

	desc := "list assets for a release"
	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/%s/assets", github, owner, repo, viper.GetString("id"))
	token := viper.GetString("token")
	method := http.MethodGet

	err := validate(map[string]string{
		"user": owner,
		"repo": repo,
	})
	if err != nil {
		return fmt.Errorf("%s: %v", desc, err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Info(desc)

	var result []interface{}
	_, err = SendRequest(url, method, nil, token, "", &result)
	if err != nil {
		return err
	}

	color.Green("%20s    %10s    %10s    %s\n", "created", "id", "size", "name")
	for _, v := range result {
		item, _ := v.(map[string]interface{})
		fmt.Printf("%20v    %10d    %10d    %v\n",
			item["created_at"], int64(item["id"].(float64)), int64(item["size"].(float64)), item["name"])
	}

	return nil
}

func GetRelease(owner string, repo string, releaseId string) (*Release, error) {

	desc := "get a single release"
	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/%s", github, owner, repo, releaseId)

	err := validate(map[string]string{
		"user":       owner,
		"repo":       repo,
		"release_id": releaseId,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %v", desc, err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Info(desc)

	var release = Release{}
	_, err = SendRequest(url, http.MethodGet, nil, "", "", &release)
	if err != nil {
		return nil, err
	}

	result, err := json.MarshalIndent(release, "", "\t")
	fmt.Printf("%s:\n%s\n", desc, string(result))

	return &release, nil
}

func GetReleaseByTag(owner string, repo string, tag string) (*Release, error) {

	desc := "get a single release"
	github := viper.GetString("github")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/tags/%s", github, owner, repo, tag)

	err := validate(map[string]string{
		"user": owner,
		"repo": repo,
		"tag":  tag,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %v", desc, err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Info(desc)

	var release = Release{}
	_, err = SendRequest(url, http.MethodGet, nil, "", "", &release)
	if err != nil {
		return nil, err
	}

	result, err := json.MarshalIndent(release, "", "\t")
	fmt.Printf("%s:\n%s\n", desc, string(result))

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

	desc := "create a release"
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
		return fmt.Errorf("%s: %v", desc, err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Info(desc)

	var result map[string]interface{}
	resp, err := SendRequest(url, method, requestByte, token, "", &result)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusCreated {

		logrus.WithFields(logrus.Fields{
			"id":       int64(result["id"].(float64)),
			"tag_name": result["tag_name"],
			"url":      result["url"],
		}).Infof("%s success", desc)

		return nil
	}

	printErrors(desc, result)

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

func UploadAsset(owner string, repo string, filename string, label string) error {

	desc := "upload a release asset"
	//github := viper.GetString("github")
	github := "https://uploads.github.com"
	id := viper.GetString("id")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/%s/assets?name=%s", github, owner, repo, id, filename)
	token := viper.GetString("token")
	method := http.MethodPost

	if label != "" {
		url += fmt.Sprintf(",label=%s", label)
	}

	err := validate(map[string]string{
		"user":     owner,
		"repo":     repo,
		"token":    token,
		"id":       id,
		"filename": filename,
	})
	if err != nil {
		return fmt.Errorf("%s: %v", desc, err)
	}

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Info(desc)

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%s, file read failed: %v", desc, err)
	}

	mime, _ := mimetype.Detect(buf)

	var result map[string]interface{}
	resp, err := SendRequest(url, method, buf, token, mime, &result)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusCreated {

		logrus.WithFields(logrus.Fields{
			"id":   int64(result["id"].(float64)),
			"name": result["name"],
			"url":  result["url"],
		}).Infof("%s success", desc)

		return nil
	}

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest {

		logrus.WithFields(logrus.Fields{
			"documentation_url": result["documentation_url"],
			"message":           result["message"],
		}).Errorf("%s failed", desc)
	}

	printErrors(desc, result)

	return nil
}

func printErrors(desc string, result map[string]interface{}) {

	if val, ok := result["errors"]; ok {

		if errors, ok := val.([]interface{}); ok {

			logrus.Errorf("%s failed with %d errors:", desc, len(errors))

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
}

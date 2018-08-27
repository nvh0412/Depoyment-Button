package main

import (
	"os"
	"fmt"
	"github.com/nvh0412/go-octokit/octokit"
  "strconv"
)

var owner = os.Getenv("GITHUB_OWNER")
var repo = os.Getenv("GITHUB_REPO")
var accessToken = os.Getenv("GITHUB_TOKEN")

func getTagName(client octokit.Client) string {
	urlLatest, err := octokit.ReleasesLatestURL.Expand(octokit.M{"owner": owner, "repo": repo})
	if err != nil {
		fmt.Println(err)
		return ""
	}

	release, result := client.Releases(urlLatest).Latest()

	if result.HasError() {
		fmt.Println(result)
		return ""
	}

	tag, err := strconv.ParseFloat(release.TagName, 64)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strconv.FormatFloat(tag + 0.1, 'f', 1, 64)
}

func createNewRelease(client octokit.Client, tagName string) *octokit.Release {
	params := octokit.ReleaseParams{
		TagName: tagName,
		TargetCommitish: "master",
		Name: "Release Production",
	}

	url, _ := octokit.ReleasesURL.Expand(octokit.M{"owner": owner, "repo": repo})
	release, res := client.Releases(url).Create(params)

	if res.HasError() {
		fmt.Println(res)
	}

	return release
}

func createDeployment(client octokit.Client, tagName string) *octokit.Deployment {
	params := octokit.DeploymentParams{
		Ref: tagName,
		Description: "Release production",
	}

	deployment, res := client.Deployments().Create(nil, octokit.M{"owner": owner, "repo": repo}, params)

	if res.HasError() {
		fmt.Println(res)
	}

	return deployment
}

func HandleRequest() (string, error) {
	token := octokit.TokenAuth{AccessToken: accessToken}
	client := octokit.NewClient(token)

	tagName := getTagName(*client)
	release := createNewRelease(*client, tagName)
	deployment := createDeployment(*client, tagName)

	return fmt.Sprintf("Release version %s has been released!\n Deployment succeed! URL: %s", release.TagName, deployment.URL), nil
}

func main() {
	lambda.Start(HandleRequest)
}

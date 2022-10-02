package content

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-shiori/go-readability"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// Save save the content to the destination github repository
func Save(url string, name string) {
	content, domain, title, err := getWebContent(url)
	cobra.CheckErr(err)

	ghClient := buildGithubClient()
	username := viper.GetString("username")
	repo := viper.GetString("repository")
	branch := viper.GetString("branch")
	commit := viper.GetString("commit")

	content = addMetadata(content, title)

	opt := &github.RepositoryContentFileOptions{
		Message: &commit,
		Content: content,
		Branch:  &branch,
	}

	_, _, err = ghClient.Repositories.CreateFile(
		context.Background(),
		username,
		repo,
		buildContentPath(domain, name),
		opt,
	)
	cobra.CheckErr(err)

	fmt.Fprintln(os.Stdout, "Successfully save the content!")
}

func getWebContent(webURL string) (markdown []byte, domain string, title string, err error) {
	u, err := url.Parse(webURL)
	if err != nil {
		return
	}

	domain = u.Hostname()

	resp, err := http.Get(webURL)
	if err != nil {
		return
	}

	article, err := readability.FromReader(resp.Body, u)
	if err != nil {
		return
	}

	title = article.Title

	converter := md.NewConverter("", true, nil)
	markdownStr, err := converter.ConvertString(article.Content)
	if err != nil {
		return
	}

	markdown = []byte(markdownStr)
	return
}

func addMetadata(content []byte, name string) []byte {
	strContent := string(content)

	updatedContent := fmt.Sprintf("%s\n%s", buildMetadata(name), strContent)

	return []byte(updatedContent)
}

func buildMetadata(title string) string {
	var metaData string = fmt.Sprintf("---\ntitle: \"%s\"\ntags:", title)

	tags := strings.Split(viper.GetString("tags"), ",")
	for _, tag := range tags {
		metaData = fmt.Sprintf("%s\n- %s", metaData, tag)
	}

	metaData = fmt.Sprintf("%s\n---\n", metaData)

	return metaData
}

func buildGithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: viper.GetString("token")})
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func buildContentPath(domain, name string) string {
	return fmt.Sprintf("%s/%s/%s.md", viper.GetString("folder"), domain, name)
}

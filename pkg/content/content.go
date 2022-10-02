package content

import (
	"context"
	"fmt"
	"net/url"
	"os"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// Save save the content to the destination github repository
func Save(url string, name string) {
	content, domain, err := getWebContent(url)
	cobra.CheckErr(err)

	ghClient := buildGithubClient()
	username := viper.GetString("username")
	repo := viper.GetString("repository")
	branch := viper.GetString("branch")
	commit := viper.GetString("commit")

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

func getWebContent(webURL string) (markdown []byte, domain string, err error) {
	u, err := url.Parse(webURL)
	if err != nil {
		return
	}

	domain = u.Hostname()
	converter := md.NewConverter(domain, true, nil)

	markdownStr, err := converter.ConvertURL(webURL)
	if err != nil {
		return
	}

	markdown = []byte(markdownStr)
	return
}

func buildGithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: viper.GetString("token")})
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func buildContentPath(domain, name string) string {
	return fmt.Sprintf("%s/%s/%s", viper.GetString("folder"), domain, name)
}

package cmd

import (
	"context"
	"fmt"
	"github.com/google/go-github/v25/github"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"os"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "",
	Long: ``,
	Run: applyEntryPoint,
}

func init() {
	rootCmd.AddCommand(applyCmd)
}


func listRepos() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	type Repository struct {
		Name        githubv4.String
		Description githubv4.String
	}

	var query struct {
		Viewer struct {
			Repositories struct {
				Nodes []Repository
				PageInfo  struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"repositories(first: $repositoryPageSize, after: $repositoryCursor)"`
			Login     githubv4.String
			CreatedAt githubv4.DateTime
		}
	}

	variables := map[string]interface{}{
		"repositoryPageSize": githubv4.Int(5),
		"repositoryCursor":  (*githubv4.String)(nil),
	}


	for {
		err := client.Query(context.Background(), &query, variables)
		if err != nil {
			panic(err)
		}

		if !query.Viewer.Repositories.PageInfo.HasNextPage {
			break
		}

		for _, x := range query.Viewer.Repositories.Nodes {
			fmt.Printf("%+v\n", x)
		}

		variables["repositoryCursor"] = githubv4.NewString(query.Viewer.Repositories.PageInfo.EndCursor)
	}

}

func createRepo(name string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repo := &github.Repository{
		Name:    github.String(name),
		Private: github.Bool(false),
	}
	_, _, err := client.Repositories.Create(ctx, "", repo)

	if err != nil {
		panic(err)
	}

}


func applyEntryPoint(cmd *cobra.Command, args []string) {
	createRepo("the-best")
}

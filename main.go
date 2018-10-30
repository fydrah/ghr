// Copyright (C) 2018 fydrah
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
// fydrah <flav.hardy@gmail.com>

// ghr creates GitHub releases based on a given annotated tag
// Usage of ghr:
//  -owner string
//      GitHub repository owner name. GHR_OWNER
//  -repository string
//      GitHub repository name. GHR_REPOSITORY
//  -tag string
//      Git annotated tag to release. GHR_TAG
//  -token string
//      GitHub Token. GHR_TOKEN
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"regexp"
	"strings"
)

type ghr struct {
	token        string
	repository   string
	tag          string
	owner        string
	message      string
	githubClient *github.Client
	ctx          context.Context
}

var g ghr

func genClient() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), ctx
	// list all repositories for the authenticated user
	//repos, _, err := client.Repositories.List(ctx, "", nil)
}

func getTag() (*github.Tag, error) {
	var (
		err error
		ref *github.Reference
		tag *github.Tag
	)
	if ref, _, err = g.githubClient.Git.GetRef(g.ctx, g.owner, g.repository, fmt.Sprintf("tags/%v", g.tag)); err != nil {
		return nil, fmt.Errorf("Failed to retrieve tag reference: %v", err)
	}
	if tag, _, err = g.githubClient.Git.GetTag(g.ctx, g.owner, g.repository, *ref.Object.SHA); err != nil {
		return nil, fmt.Errorf("Failed to retrieve tag: %v", err)
	}
	return tag, nil
}

// https://gist.github.com/r0l1/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func createRelease() error {
	r := github.RepositoryRelease{
		TagName: &g.tag,
		Name:    &g.tag,
		Body:    &g.message,
	}
	if _, _, err := g.githubClient.Repositories.CreateRelease(g.ctx, g.owner, g.repository, &r); err != nil {
		return fmt.Errorf("Failed to create release: %v", err)
	}
	return nil
}

func init() {
	flag.StringVar(&g.repository, "repository", os.Getenv("GHR_REPOSITORY"), "GitHub repository name. GHR_REPOSITORY")
	flag.StringVar(&g.owner, "owner", os.Getenv("GHR_OWNER"), "GitHub repository owner name. GHR_OWNER")
	flag.StringVar(&g.token, "token", os.Getenv("GHR_TOKEN"), "GitHub Token. GHR_TOKEN")
	flag.StringVar(&g.tag, "tag", os.Getenv("GHR_TAG"), "Git annotated tag to release. GHR_TAG")
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			fmt.Printf("%v required\n\n", f.Name)
			flag.PrintDefaults()
			os.Exit(1)
		}
	})
	g.githubClient, g.ctx = genClient()
}

func main() {
	fmt.Print(`
#####################
GHR - GitHub Releaser
#####################

`)
	fmt.Printf("[%v] Creating release %v...\n\n", g.repository, g.tag)
	if tag, err := getTag(); err != nil {
		log.Fatal(err)
	} else {
		re := regexp.MustCompile("-----BEGIN PGP SIGNATURE-----(.|\n)*-----END PGP SIGNATURE-----")
		g.message = re.ReplaceAllString(*tag.Message, "")
		fmt.Printf("### BEGIN message ###\n\n%v\n### END message ###\n\n", g.message)
	}

	c := askForConfirmation(fmt.Sprintf("Create %v release with the following message ?", g.tag))

	if !c {
		fmt.Print("Abort...\n")
		return
	}

	if err := createRelease(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print("Release created with success\n")
	}
}

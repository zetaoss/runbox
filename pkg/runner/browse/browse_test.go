package browse

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"github.com/zetaoss/runbox/pkg/testutil"
)

var browse1 *Browse

func init() {
	d := testutil.NewDocker()
	browse1 = New(box.New(d))
}

func TestRun(t *testing.T) {
	testCases := []struct {
		urlString string
		want      string
	}{
		{
			"https://api.github.com",
			"<html><head></head><body><pre style=\"word-wrap: break-word; white-space: pre-wrap;\">{\n  \"current_user_url\": \"https://api.github.com/user\",\n  \"current_user_authorizations_html_url\": \"https://github.com/settings/connections/applications{/client_id}\",\n  \"authorizations_url\": \"https://api.github.com/authorizations\",\n  \"code_search_url\": \"https://api.github.com/search/code?q={query}{&amp;page,per_page,sort,order}\",\n  \"commit_search_url\": \"https://api.github.com/search/commits?q={query}{&amp;page,per_page,sort,order}\",\n  \"emails_url\": \"https://api.github.com/user/emails\",\n  \"emojis_url\": \"https://api.github.com/emojis\",\n  \"events_url\": \"https://api.github.com/events\",\n  \"feeds_url\": \"https://api.github.com/feeds\",\n  \"followers_url\": \"https://api.github.com/user/followers\",\n  \"following_url\": \"https://api.github.com/user/following{/target}\",\n  \"gists_url\": \"https://api.github.com/gists{/gist_id}\",\n  \"hub_url\": \"https://api.github.com/hub\",\n  \"issue_search_url\": \"https://api.github.com/search/issues?q={query}{&amp;page,per_page,sort,order}\",\n  \"issues_url\": \"https://api.github.com/issues\",\n  \"keys_url\": \"https://api.github.com/user/keys\",\n  \"label_search_url\": \"https://api.github.com/search/labels?q={query}&amp;repository_id={repository_id}{&amp;page,per_page}\",\n  \"notifications_url\": \"https://api.github.com/notifications\",\n  \"organization_url\": \"https://api.github.com/orgs/{org}\",\n  \"organization_repositories_url\": \"https://api.github.com/orgs/{org}/repos{?type,page,per_page,sort}\",\n  \"organization_teams_url\": \"https://api.github.com/orgs/{org}/teams\",\n  \"public_gists_url\": \"https://api.github.com/gists/public\",\n  \"rate_limit_url\": \"https://api.github.com/rate_limit\",\n  \"repository_url\": \"https://api.github.com/repos/{owner}/{repo}\",\n  \"repository_search_url\": \"https://api.github.com/search/repositories?q={query}{&amp;page,per_page,sort,order}\",\n  \"current_user_repositories_url\": \"https://api.github.com/user/repos{?type,page,per_page,sort}\",\n  \"starred_url\": \"https://api.github.com/user/starred{/owner}{/repo}\",\n  \"starred_gists_url\": \"https://api.github.com/gists/starred\",\n  \"topic_search_url\": \"https://api.github.com/search/topics?q={query}{&amp;page,per_page}\",\n  \"user_url\": \"https://api.github.com/users/{user}\",\n  \"user_organizations_url\": \"https://api.github.com/user/orgs\",\n  \"user_repositories_url\": \"https://api.github.com/users/{user}/repos{?type,page,per_page,sort}\",\n  \"user_search_url\": \"https://api.github.com/search/users?q={query}{&amp;page,per_page,sort,order}\"\n}\n</pre></body></html>\n",
		},
		{
			"https://jsonplaceholder.typicode.com/todos/100",
			"<html><head></head><body><pre style=\"word-wrap: break-word; white-space: pre-wrap;\">{\n  \"userId\": 5,\n  \"id\": 100,\n  \"title\": \"excepturi a et neque qui expedita vel voluptate\",\n  \"completed\": false\n}</pre></body></html>\n",
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.urlString), func(t *testing.T) {
			got, err := browse1.Run(tc.urlString)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

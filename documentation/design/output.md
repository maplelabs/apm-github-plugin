# Github Audit Plugin Output Json Response

Github audit plugin will be providing output in following json format:

## Commits related
### Type: commit 
```json
{
    "type": "commit",
    "repo_type":"github",
    "repo_name":"test_repo",
    "repo_url":"https://github.com/testurl",
    "created_at": "2022-08-30T16:25:04Z",
    "message": "initial PR",
    "url": "https://api.github.com/repos/pramurthy/sf-apm-agent/commits/9a5a338a2b6f9d435faa9adbda1f952276c1aea8",
    "committer": {
        "id": "1233",
        "user": "name1"
    },
    "sha": "9a5a338a2b6f9d435faa9adbda1f952276c1aea8"
}
```

## Pull requests related
### Type: pull request
```json
{
    "type": "pull_request",
    "repo_type":"github",
    "repo_name":"test_repo",
    "repo_url":"https://github.com/testurl",
    "created_at": "2022-08-30T16:25:04Z",
    "updated_at": "2022-08-30T16:25:04Z",
    "closed_at": "",
    "state": "open",
    "pull_request_no": "1",
    "title": "initial PR",
    "url": "https://api.github.com/repos/maplelabs/github-audit/pulls/1",
    "merged_at": "",
    "closed_at": "",
    "merge_commit_sha": "87157431fa8922d17f4dadas3437c05fc72f12ae",
    "assignees": [
        {
            "id": "1233",
            "user": "name1"
        },
        {
            "id": "1234",
            "user": "name2"
        }
    ],
    "reviewers": [
        {
            "id": "1233",
            "user": "name1"
        },
        {
            "id": "1234",
            "user": "name2"
        }
    ],
    "request_from_repo": {
        "name": "repo1",
        "url": "",
        "private": "false",
        "sha": "87157431fa8922d17f43ce7c697c05fc72f12ae",
        "branch": "master",
        "by_user": {
            "id": "1234",
            "name": "name2"
        }
    },
    "merge_to_repo": {
        "name": "repo2",
        "url": "",
        "private": "false",
        "sha": "87157431fa8922d17f43ce7c697c05fc72f12ae",
        "branch": "master"
    },
    "merged_by": {
        "id": "1233",
        "user": "name1"
    }
}
```
### Type: pull request commits
```json
{
    "type": "pull_request_commits",
    "repo_type":"github",
    "repo_name":"test_repo",
    "repo_url":"https://github.com/testurl",
    "pull_request_no": "1",
    "title": "initial PR",
    "url": "https://api.github.com/repos/maplelabs/github-audit/pulls/1",
    "commits": [
        {
            "created_at": "2022-08-30T16:25:04Z",
            "message": "initial PR",
            "url": "https://api.github.com/repos/pramurthy/sf-apm-agent/commits/9a5a338a2b6f9d435faa9adbda1f952276c1aea8",
            "committer": {
                "id": "1233",
                "user": "name1"
            },
            "sha": "9a5a338a2b6f9d435faa9adbda1f952276c1aea8"
        }
    ]
}
```
### Type: pull request comments
```json
{
    "type": "pull_request_comments",
    "repo_type": "github",
    "repo_name":"test_repo",
    "repo_url":"https://github.com/testurl",
    "pull_request_no": "1",
    "title": "initial PR",
    "url": "https://api.github.com/repos/maplelabs/github-audit/pulls/1",
    "comments": [
        {
            "comment_id": "963040674",
            "message": "some comment",
            "created_at": "2022-08-30T16:25:04Z",
            "updated_at": "2022-08-30T16:25:04Z",
            "url": "https://api.github.com/repos/maplelabs/github-audit/pulls/comments/963040674",
            "created_by": {
                "id": "1233",
                "user": "name1"
            }
        }
    ]
}
```

### Type: pull request issues 
```json
{
    "type": "pull_request_issues",
    "repo_type": "github",
    "repo_name":"test_repo",
    "repo_url":"https://github.com/testurl",
    "pull_request_no": "1",
    "title": "initial PR",
    "url": "https://api.github.com/repos/maplelabs/github-audit/pulls/1",
    "issues": [
        {
            "issue_no": "3",
            "title": "some issue",
            "state":"open",
            "created_at": "2022-08-30T16:25:04Z",
            "updated_at": "2022-08-30T16:25:04Z",
            "closed_at": "",
            "url": "https://api.github.com/repos/maplelabs/github-audit/issues/3",
            "created_by": {
                "id": "1233",
                "user": "name1"
            },
            "assignees": [
                {
                    "id": "1233",
                    "user": "name1"
                },
                {
                    "id": "1234",
                    "user": "name2"
                }
            ]
        }
    ]
}
```
## Issues related
### Type: issue
```json
{
    "type": "issue",
    "repo_type":"github",
    "repo_name":"test_repo",
    "repo_url":"https://github.com/testurl",
    "issue_no": "3",
    "state":"open",
    "title": "some issue",
    "created_at": "2022-08-30T16:25:04Z",
    "updated_at": "2022-08-30T16:25:04Z",
    "closed_at": "",
    "url": "https://api.github.com/repos/maplelabs/github-audit/issues/3",
    "created_by": {
        "id": "1233",
        "user": "name1"
    },
    "assignees": [
        {
            "id": "1233",
            "user": "name1"
        },
        {
            "id": "1234",
            "user": "name2"
        }
    ]
}
```
# git-branch-delete

Interactively delete local git branches.
## How to run

```bash
make build && go install

git-branch-delete
```

## Under the hood

This tool makes use of [`git branch --merged/--no-merged`](https://git-scm.com/docs/git-branch)
to display an enriched list of local branches and their status (merged or not) against `origin/master` (to be 
configurable in the future).

You can navigate the list interactively and pick which branches to delete locally.


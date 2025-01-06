# GitHub Readme Stats 📊

This project builds upon the previous work by Jainam Desai[^1]. Please consider
starring this repository if you find it useful! ⭐️

## Example Output

```
🕰️ I get my jam on during the evening!

🌞 Morning  	180    commits	███░░░░░░░░░░░░░░░░░░░░░░░░░░░	4.37%
🌆 Daytime  	881    commits	████████████████░░░░░░░░░░░░░░	21.38%
🌃 Evening  	1619   commits	██████████████████████████████	39.30%
🌙 Night    	1440   commits	██████████████████████████░░░░	34.95%
```

```
📅 I'm most productive on Saturdays!

Monday      	666    commits	██████████████████████████░░░░	16.17%
Tuesday     	635    commits	█████████████████████████░░░░░	15.41%
Wednesday   	500    commits	████████████████████░░░░░░░░░░	12.14%
Thursday    	468    commits	██████████████████░░░░░░░░░░░░	11.36%
Friday      	485    commits	███████████████████░░░░░░░░░░░	11.77%
Saturday    	750    commits	██████████████████████████████	18.20%
Sunday      	616    commits	████████████████████████░░░░░░	14.95%
```

```
🧪 Python for the win!

Python      	21     repos	██████████████████████████████	15.22%
HTML        	11     repos	███████████████░░░░░░░░░░░░░░░	7.97%
C           	10     repos	██████████████░░░░░░░░░░░░░░░░	7.25%
Makefile    	10     repos	██████████████░░░░░░░░░░░░░░░░	7.25%
Shell       	10     repos	██████████████░░░░░░░░░░░░░░░░	7.25%
```

## Installation

To integrate this action into your repository, start by creating a GitHub Action
workflow file. In your repository, navigate to/create
`.github/workflows/update-readme.yaml`.

```yaml
name: Update Readme

on:
  schedule:
    - cron: "1 3 * * *" # best to run at random time of day
  workflow_dispatch:

jobs:
  Update:
    name: Update Stats
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: mahyarmirrashed/github-readme-stats@main # v1.0 and other tags exist, as well
        with:
          GITHUB_TOKEN: ${{ secrets.PAT }}
          TIMEZONE: "America/Winnipeg"
          FEATURES: "DAY_STATS,WEEK_STATS,LANGUAGE_STATS"
```

Next, generate a Personal Access Token (PAT) in your
[GitHub Settings Page](https://github.com/settings/tokens). There, select the
`Generate new token: fine-grained, repo-scoped` option.

Give your token a name (e.g. "Profile Stats"), an expiration (e.g. 30 days), and
a description.

> [!IMPORTANT]  
> For the permissions, you need to give access to `All repositories`, and
> selecting `read-only` `Contents` access for the `Repository Permissions`.
> Then, generate the token.

Once the token is generated, copy it, and create a new repository secret for
your repository. (It will be located at an address like
`https://github.com/<user>/<repo>/settings/secrets/actions`.) There, create a
new repository secret, name it `PAT`, and paste the generated token in the body!

> [!WARNING]  
> Generally, it is not safe to give such encompassing permissions for a Personal
> Access Token, like `read` permissions for all repositories. However, there is
> <ins>no other way</ins> to get the statistics you desire without this
> permission.
>
> I created this GitHub action to keep the processing of this information local
> to GitHub's infrastructure (unlike WakaTime[^2]) to reduce potential attack
> surfaces.

## Usage

This action operates by updating a specific section in your `README.md` file.

### Adding the Markers

Insert the following markers where you want the statistics to appear:

```markdown
<!-- README-STATS:START -->

<!-- README-STATS:END -->
```

When the GitHub Action runs, it will populate the section between these markers
with your chosen statistics.

### Example Configuration

To include daily, weekly, and language statistics, set the `FEATURES` parameter
as follows:

```yaml
FEATURES: "DAY_STATS,WEEK_STATS,LANGUAGE_STATS"
```

## Configuration

You can customize the statistics included and their order by modifying the
`FEATURES` parameter in your workflow file.

### Available Includes

- `DAY_STATS`: Commit statistics based on the time of day.
- `WEEK_STATS`: Commit statistics based on the day of the week.
- `LANGUAGE_STATS`: Programming language usage statistics.

## Troubleshooting

- **`README.md` Not Found:** Ensure that `README.md` exists in the root of your
  repository and that you’ve added the `<!-- README-STATS:START -->` and
  `<!-- README-STATS:END -->` markers.
- **Permissions Error:** Verify that your PAT has the necessary scopes (`repo`
  and `read:user`).
- **Incorrect Timezone:** Make sure the timezone string is valid (e.g.,
  "America/Winnipeg").

## Contributions

Contributions are welcome! Please open an issue or submit a pull request for any
enhancements or bug fixes.

[^1]: https://github.com/th3c0d3br34ker/github-readme-info

[^2]: https://wakatime.com/plugins/status?onboarding=true

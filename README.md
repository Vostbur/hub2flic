App for migration from GitHub to GitFlic
----------------------------------------

**Edit configuration file in YAML format. For example:**

    github_name: <username>
    github_token: <GitHub token>
    per_page: 100
    gitflic_name: <username>
    gitflic_pass: <userpass>
    gitflic_token: <GitFlic token>
    clone_path: ./tmp

**Run CLI:**

    Usage of ./hub2flic:
    -config string
            Path to YAML config file
    -gist string
            Clone gists: 'no' (or without key), 'yes' or 'single' (default "no")
    -repo string
            Repository name

**Example:**

transfer all repos without gists

    $ ./hub2flic -config conf.yaml

transfer repo by name

    $ ./hub2flic -config conf.yaml -repo github_repo_name

transfer all repos and all gists each in its own repo

    $ ./hub2flic -config conf.yaml -gist yes

transfer all repos and all gists in one repo

    $ ./hub2flic -config conf.yaml -gist single

_TODO: Подробнее допишу позже_
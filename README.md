# GHR - GitHub Releaser

Create a GitHub released based on:

* Tag name, ex: "2.3.0"
* Tag Annotation


## Install

```shell
go get github.com/fydrah/ghr
```

## Usage

```shell
Usage of ghr:
  -owner string
    	GitHub repository owner name. GHR_OWNER
  -repository string
    	GitHub repository name. GHR_REPOSITORY
  -tag string
    	Git annotated tag to release. GHR_TAG
  -token string
    	GitHub Token. GHR_TOKEN
```

## Requirements

* Tag should be an **annotated** tag
* Your token must have correct [scopes](https://developer.github.com/apps/building-oauth-apps/understanding-scopes-for-oauth-apps/)
* You must own the repository or have correct access to organization repository

## Example

```shell
$ export GHR_TOKEN=mygithubtoken_cancreatereleases
$ ghr -owner fydrah \
    -repository loginapp \
    -tag 2.4.0
    #####################
    GHR - GitHub Releaser
    #####################

    [loginapp] Creating release 2.4.0...

    ### BEGIN message ###

    release 2.4.0

    Image for this release available at:

        quay.io/fydrah/loginapp:2.4.0

    Features:

        * Extract HTML templates from binary (#6)

    Thanks for using loginapp !

    ### END message ###

    Create 2.4.0 release with the following message ? [y/n]: y
    Release created with success
```

You can see and example of producted releases [here](https://github.com/fydrah/loginapp/releases).

Enjoy !

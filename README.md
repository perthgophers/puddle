[![Go Report Card](https://goreportcard.com/badge/github.com/perthgophers/puddle)](https://goreportcard.com/report/github.com/perthgophers/puddle)

# Puddle
The Perth Gophers Slackbot Mud Project

# Introduction

This is Puddle, the Perth Gopher's collaborative project. Created to help us learn Go and to teach others Go. 

Puddle is a mud. Puddle is a slackbot. Puddle is a slackbot mud. How do we design a Mud? Messily.

* [Mud Design concepts](https://www.gammon.com.au/forum/?id=10147)
* [How do Muds work?](http://www.livinginternet.com/d/dw.htm)
* [Similar project for inspiration](https://github.com/Streamweaver/pogomud)

# Spinup

[Go 1.6](https://golang.org/) and above is required for this project.

```
$ go get github.com/perthgophers/puddle
$ go get github.com/Masterminds/glide
$ cd  $GOPATH/src/github.com/perthgophers/puddle
$ glide install
$ go run main.go --slack_token <SLACK_API_TOKEN>
```

# Collaborative etiquette

This project is an **Open Source** project encouraging everyone to be fearless in their contributions. This document lists some ground rules to foster a happy collaborative atmosphere.

* Individuals making valuable contributions are added as collaborators (commit access).

* This project is more like an open wiki than a standard guarded project.


## For issue reporters

* __Don't be afraid!__ You don't need to submit code to contribute to the success of a project! Don't hold back on submitting issues, commenting on discussions and helping people out.

* __Tell us everything.__ When filing bug reports, be generous in the details: environment, OS, browser, steps to replicate, and so on.

## For contributors

* __Check development notes.__ Any details on setting up a development environment, running tests, releasing versions and such are kept in a file like `NOTES.md`.

* __Be a good citizen.__ Try your best to adhere to the established styles of the project. This doesn't mean that you shouldn't break them, but be prepared to have a reason if you do.

* __Be informative.__ Format your pull requests nicely. Include screenshots if applicable.

* __Don't Panic.__ This project is going to be messy. Don't be afraid of submitting any code, even if you yourself think it's absolutely terrible. This is, above all else, a learning project.

## For collaborators

Individuals making valuable contributions are encouraged to be added as collaborators (commit access). This status should be given *liberally.* If things break, we don't care. Someone will fix it, and then give a valuable critique of your code!

* __Work in branches then send a PR.__ Just because you have full commit access doesn't mean you shouldn't use pull requests. PRs are a great way to solicit feedback from co-collaborators and to give them a nice overview of what's going on.

* __Review outstanding PRs.__ Feel free to merge any you see fit, and leave comments on anything that needs revisions. If you don't feel comfortable merging them, at least comment with a `:+1:` to signal your co-collaborators that it's passed your review.

* __Push directly for micro-fixes only.__ Only push to `master` for trivial updates that would be too noisy to notify your teammates of, such as typo fixes.

* __You can self-merge your PRs.__ Sometimes the rest of the team may be inactive, in which case, use your best judgement to self-merge PRs. If you mess it up, so what?

* __Keep history clean.__ No `--force` pushing on branches that aren't yours. No --force pushing at all, please.

* __Communicate.__ You can use GitHub issues to communicate with your co-collaborators. Feel free to `@mention` them in issues. You can also chat in the #puddle channel on Slack.

* __Be nice.__ The collaborator team is nice, so we are all nice. Don't be an asshole to anyone or you'll be fired, from a cannon, into the sun.

## For founders

* __Set up continuous-integration.__ A CI service like [Travis] will inspect PR's so you don't have to.

* __Be decisive.__ Don't be afraid to put your foot down on issues like design decisions.

* __Thank people.__ People put in their hard work. Gratitude goes a long way.

* __Be nice.__ The open source community is nice, so we are nice.

## Acknowledgements

> Thank you to the Perth Gophers community.

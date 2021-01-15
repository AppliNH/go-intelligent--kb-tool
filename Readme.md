# Go-Intelligent Knowledge base tool

*This has been made in the context of an assignment*


*By : Thomas MARTIN, Séraphin HENRY, Victor CAVERO*
## What is it

Go-KB is CLI which aims to store the data you give in order to make the best suggestions the next time you enter data.

Here, it works with a combination of name and surname.

Let's say you give `Name=Toto` and `Surname=Paolo`, then `Name=Titi`and `Surname=Paolo`.

The next time you enter `Surname=Paolo` it will propose you to chose between `Toto` and `Titi` as a name. Or you can also enter a new one.

Everything is saved made to lower key and saved in a Key-Value database.

## Run it

You will need to have the golang runtime to be installed [(see here)](https://golang.org/).

Then you can :

- run `go get -u github.com/applinh/go-intelligent-kb-tool`
- Then you can run the tool with the command `go-intelligent-kb-tool`

OR

- clone the repo, `cd` inside it, and run `go run main.go`

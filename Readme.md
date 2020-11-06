# Go-Intelligent Knowledge base tool

*This has been made in the context of an assignment*
*Thomas MARTIN, Séraphin HENRY, Victor CAVERO*
## What is it

Go-KB is CLI which aims to store the data you give in order to make the best suggestions the next time you enter data.

Here, it works with a combination of name and surname.

Let's say you give `Name=Toto` and `Surname=Paolo`, then `Name=Titi`and `Surname=Paolo`.

The next time you enter `Surname=Paolo` it will propose you to chose between `Toto` and `Titi` as a name. Or you can also enter a new one.

## Run it

You will need to have the golang runtime to be installed [(see here)](https://golang.org/).

Then you can run :

`go run main.go`
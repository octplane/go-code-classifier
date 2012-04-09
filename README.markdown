GO code classifier
==================

This is an experiment to work around the fact that most public paste website are not able to guess the pasted syntax from the pasted snipper (cf http://stackoverflow.com/questions/9465215/pastie-with-api-and-language-detection ). I've decided to try and learn Go and using a [Bayesian Classifier](https://github.com/jbrukh/bayesian), I'm writing a language classifier.

Once I've something that can generate a working scanner, I will implement the website to paste your snippets.

Building the code
=================

Clone the repository and initialize the submodule:

```shell
git clone git://github.com/octplane/go-code-classifier.git
cd go-code-classifier
git submodule init
git submodule update
export GOLANG=$PWD
go build src/main.go
```

That's all for now.

Licence
=======

This software is released under the WTFPL2.

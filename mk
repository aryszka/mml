#! /usr/bin/env mml --require mml/mk

envDefault("image", "foo-image")
flagDefault("test", "false")

println(mk.test)

// this currently not supported by the language, because it is not yet defined:
// maybe just need to be moved below?
mk.build = "main.mml" -> when(mk.generate)

mk.format = "**/*.mml" -> format // this could also be a default

let parse mk.wrap(require mml/parse)
mk.generate = "syntax.g"
	-> parse.read
	-> parse.generate("syntax.mml")
	-> always

let docker sh("docker build -t" mk.image ".")
mk.docker       = docker
mk.dockerUpdate = docker -> when(mk.build)

// similarly, there can be an mml/script, and mml/mk can be just an extension to it

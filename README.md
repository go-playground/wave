Package wave
=============
![Project status](https://img.shields.io/badge/beta-0.1.0-orange.svg)
[![GoDoc](https://godoc.org/github.com/go-playground/wave?status.svg)](https://godoc.org/github.com/go-playground/wave)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Package wave is a thin helper layer on top of Go's net/rpc

### Why?
Ths intention of this library is to provide a thin wrapper around the std net/rpc package allowing
the user to add functionality via hooks instead of creating a whole new framework.

**NOTES:** currently there are no hooks, first one will be hooking into the `Register` and `RegisterName`
functions to allow any service discovery to be handled. This project is very much in it's early stages
any contributions are welcome.


Package Versioning
----------
I'm jumping on the vendoring bandwagon, you should vendor this package as I will not
be creating different version with gopkg.in like allot of my other libraries.

Why? because my time is spread pretty thin maintaining all of the libraries I have + LIFE,
it is so freeing not to worry about it and will help me keep pouring out bigger and better
things for you the community.

License
------
Distributed under MIT License, please see license file in code for more details.

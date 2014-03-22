stein - the best way to consume TAP
===================================

`stein` parses TAP-Y, TAP-J and vanilla TAP and displays it in a web page.

TAP-Y and TAP-J follow the [TAPOUT](https://github.com/rubyworks/tapout) [specification](https://github.com/rubyworks/tapout/wiki/TAP-Y-J-Specification), and vanilla TAP follows the original TAP [specification](http://testanything.org/tap-specification.html) and [TAP13](http://testanything.org/tap-version-13-specification.html).

status
------

`stein` is under heavy development and should not be considered production ready. The API may change at any time and there are most definitely bugs.

installation
------------

    go get github.com/beatgammit/stein

license
-------

`stein` itself is released under the MIT license (see LICENSE.MIT for details), but its dependencies have other licenses:

* [goyaml](https://github.com/goyaml/yaml): licensed under the LGPL with an exception clause (see project for details)
* [go](https://golang.org): licensed under 3-clause BSD license

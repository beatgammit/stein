stein - the best way to consume TAP
===================================

`stein` parses TAP-Y, TAP-J and vanilla TAP and displays it in a web page.

TAP-Y and TAP-J follow the [TAPOUT](https://github.com/rubyworks/tapout) [specification](https://github.com/rubyworks/tapout/wiki/TAP-Y-J-Specification), and vanilla TAP follows the original TAP [specification](http://testanything.org/tap-specification.html) and [TAP13](http://testanything.org/tap-version-13-specification.html) (currently a TODO item).

status
------

`stein` is under heavy development and should not be considered production ready. The API may change at any time and there are most definitely bugs.

installation
------------

    go get github.com/beatgammit/stein

documentation
-------------

API documentation is on [apiary](http://docs.stein.apiary.io), and source code documentation is viewable with `godoc`.

license
-------

`stein` itself is released under the 3-clause BSD license (see LICENSE.BSD for details), but its dependencies have other licenses:

* [yaml](https://github.com/go-yaml/yaml): licensed under the LGPL with an exception clause (see project for details)

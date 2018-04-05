# Go Carnival

Similar to Go Playground, Carnival can execute simple go code and return the result.

The main difference is that Carnival runs every job posted to it in it's own docker
container.

This way if code is unsafe for some reason or another can still be executed.
There are also no clock restrictions in place (currently).

---

Future goals:

* Use vgo when compiling binaries to allow for more than the stl to be used.

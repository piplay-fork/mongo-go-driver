# Sessions Design - [Design Ticket](https://jira.mongodb.org/browse/GODRIVER-463)
This document detains a design for the Driver's Session API for the Go driver. This design takes the
design laid out in the [Sessions
Specification](https://github.com/mongodb/specifications/blob/master/source/sessions/driver-sessions.rst)
and augments with components that help ensure users don't accidentally misuse a session.

## Definitions
<dl>
</dl>

## Motivation
To be considered a complete driver, the Go driver needs to implement the [Drivers Session
Specification](https://github.com/mongodb/specifications/blob/master/source/sessions/driver-sessions.rst).
While the currently design works well for many drivers there are a number of sharp edges that make a
straight forward implementation in the Go driver not feasible. For instance, the specifications
requires that a sessions be a parameter to the collection methods. Since Go does not have method
overloading this would usually require either adding the parameter to all methods or doubling the
number of methods. Default `nil` parameters are an anti-pattern in Go and doubling the number of
methods can only be done once before it becomes unwieldy. The Go driver has already accepted a
proposed design that deviates from another specifications that required doubling the number of
methods, so this design follows that design and does not double the number of methods.

Additionally, user feedback has indicated that the API provided by the driver sessions specification
has usability concerns. While it is not possible to completely fix these usability issues one driver
at a time, an intermediary step is suggested in this design that allows users to avoid the
potentially usability issue if they desire.

## High Level Design
This design follows the driver sessions specification for the majority of the design. For instance,
a session is started from a `mongo.Client` object, and a `mongo.Session` is a valid parameter for
every collection method. The `mongo.Session` type has methods that mirror those of the
`mongo.Collection` type, except that the methods on `mongo.Session` take a `mongo.Collection` as a
parameter. The `mongo.Client` attached to the `mongo.Session` object will be used for all operations
and the `mongo.Client` attached to the `mongo.Collection` type will not be used.

### `mongo.Session` As An Option
The `mongo.Session` type is a valid option for each of the `mongo.Collection` CRUD methods. This
should be implemented by first adding an `Empty` type to each of the option packages. This type
would be an empty struct type that would implement each of the interfaces for that options package.
The `mongo.Session` type embeds each of these `Empty` types, which will allow a `mongo.Session` to
be used as an option for any of the CRUD methods.

### `mongo.Session` methods
The `mongo.Session` type will have a method for each of the CRUD Specification methods. The
signatures of these methods will mirror that of the signatures of the `mongo.Collection` methods.
Additionally, the `RunCommand` method from the `mongo.Database` type will also be available on the
`mongo.Session` type and it will take a `mongo.Database` instead of a `mongo.Collection`. The method
signature for the `InsertMany` method of the `mongo.Collection` type is:
```go
func (coll *Collection) InsertMany(context.Context, []interface{}, ...insertopt.Many) (*InsertManyResult, error)
```
The method signature for the `InsertMany` method of the `mongo.Session` type is:
```go
func (s *Session) InsertMany(context.Context, *Collection, []interface{}, ...insertopt.Many) (*InsertManyResult, error)
```
If another session is passed as an option to the CRUD methods of the `mongo.Session` type, that
method will return an error.

This mirroring of methods allows a user to use a both a `mongo.Session` in a way in which they
cannot forget to pass the session into a method while also allowing them to reuse the
`mongo.Collection` and `mongo.Database` types they've used before.
## Code

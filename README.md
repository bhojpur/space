# Bhojpur Space - Data Processing Engine

The `Bhojpur Space` is a high-performance, distributed spatial data processing engine applied within
[Bhojpur.NET Platform](https://github.com/bhojpur/platform) ecosystem for delivery of applications or
services. It features an in-memory geolocation data store, spatial index, and realtime geofence. It
supports a variety of object types including lat/lon points, bounding boxes, XYZ tiles, Geohashes,
and GeoJSON.

## Key Features

- Spatial index with [search](#searching) methods, such as: Nearby, Within, and Intersects.
- Realtime [geofencing](#geofencing) through `webhooks` or [pub/sub channels](#pubsub-channels).
- Object types of [lat/lon](#latlon-point), [bbox](#bounding-box), [Geohash](#geohash), [GeoJSON](#geojson),
[QuadKey](#quadkey), and [XYZ tile](#xyz-tile).
- Support for lots of `client-side` libraries written in different programming languages.
- Variety of protocols, including [http](#http) (curl), [websockets](#websockets), [telnet](#telnet),
and the [Redis RESP](https://redis.io/topics/protocol).
- Server responses are [RESP](https://redis.io/topics/protocol) or [JSON](https://www.json.org).
- Full [command line interface](#client).
- Leader / follower [replication](#replication).
- In-memory database that persists on disk.

## Core Components

- `spacesvr` - The server engine
- `spacectl` - The command line interface tool

## Getting Started

### Getting Space

Perhaps the easiest way to get the latest `Bhojpur Space` is to use one of the pre-built release
binaries which are available for OSX, Linux, FreeBSD, and Windows. Instructions for using these
binaries are on the GitHub [releases page](https://github.com/bhojpur/space/releases).

### Docker

To run the latest stable version of the `Bhojpur Space`:

```bash
$ docker pull bhojpur/space
$ docker run -p 9851:9851 bhojpur/space
```

Visit the [Docker Hub page](https://hub.docker.com/r/bhojpur/space/) for more information.

### Homebrew (macOS)

Install the `Bhojpur Space` using [Homebrew](https://brew.sh/)

```bash
$ brew install bhojpur-space
$ spacesvr
```

### Building Bhojpur Space

The `Bhojpur Space` can be compiled and used on Linux, OSX, Windows, FreeBSD, and probably others
since the codebase is 100% Go. We support both 32 bit and 64 bit systems. [Go](https://golang.org/dl/)
must be installed on the build machine.

To build everything simply:

```bash
$ make
```

To test:

```bash
$ make test
```

### Running

For the command line options invoke:

```bash
$ ./spacesvr -h
```

To run a single server:

```bash
$ ./spacesvr

# The `spacectl` shell connects to localhost:9851

$ ./spacectl
> help
```

#### Prometheus Metrics

The `Bhojpur Space` can natively export the Prometheus metrics by setting the `--metrics-addr`
command line flag (disabled, by default). This example exposes HTTP metrics server on port 4321:

```bash
# start server and enable Prometheus metrics, listen on local interface only
$ ./spacesvr --metrics-addr=127.0.0.1:4321

# access metrics
$ curl http://127.0.0.1:4321/metrics
```

If you need to access the `/metrics` endpoints from a different host you'll have to set the flag
accordingly, e.g. set it to `0.0.0.0:<<port>>` to listen on all interfaces.

Use the [redis_exporter](https://github.com/oliver006/redis_exporter) for more advanced use cases
like extracting key values or running a lua script.

## Playing with Bhojpur Space

Basic operations:

```bash
$ ./spacectl

# add a couple of points named 'truck1' and 'truck2' to a collection named 'fleet'.
> set fleet truck1 point 33.5123 -112.2693   # on the Loop 101 in Phoenix
> set fleet truck2 point 33.4626 -112.1695   # on the I-10 in Phoenix

# search the 'fleet' collection.
> scan fleet                                 # returns both trucks in 'fleet'
> nearby fleet point 33.462 -112.268 6000    # search 6 kilometers around a point. returns one truck.

# key value operations
> get fleet truck1                           # returns 'truck1'
> del fleet truck2                           # deletes 'truck2'
> drop fleet                                 # removes all 
```

The `Bhojpur Space` has a ton of [great commands](https://docs.bhojpur.net/commands).

## Fields

The `Fields` are extra data that belongs to an object. A field is always a double precision floating
point. There is no limit to the number of fields that an object can have.

To set a field when setting an object:

```
> set fleet truck1 field speed 90 point 33.5123 -112.2693             
> set fleet truck1 field speed 90 field age 21 point 33.5123 -112.2693
```

To set a field when an object already exists:

```
> fset fleet truck1 speed 90
```

## Searching

The `Bhojpur Space` has support to search for objects and points that are within or intersects other
objects. All object types can be searched including Polygons, MultiPolygons, GeometryCollections, etc.

### Within

WITHIN searches a collection for objects that are fully contained inside a specified bounding area.

### Intersects

INTERSECTS searches a collection for objects that intersect a specified bounding area.

### Nearby

NEARBY searches a collection for objects that intersect a specified radius.

## Search options

**WHERE** - This option allows for filtering out results based on [field](#fields) values. For
example<br>```nearby fleet where speed 70 +inf point 33.462 -112.268 6000``` will return only
the objects in the 'fleet' collection that are within the 6 km radius **and** have a field named
`speed` that is greater than `70`. <br><br>Multiple WHEREs are concatenated as **and** clauses.
```WHERE speed 70 +inf WHERE age -inf 24``` would be interpreted as *speed is over 70 <b>and</b>
age is less than 24.*<br><br>The default value for a field is always `0`. Thus if you do a WHERE
on the field `speed` and an object does not have that field set, the server will pretend that the
object does and that the value is Zero.

**MATCH** - MATCH is similar to WHERE except that it works on the object id instead of fields.<br>
```nearby fleet match truck* point 33.462 -112.268 6000``` will return only the objects in the
'fleet' collection that are within the 6 km radius **and** have an object id that starts with
`truck`. There can be multiple MATCH options in a single search. The MATCH value is a simple
[glob pattern](https://en.wikipedia.org/wiki/Glob_(programming)).

**CURSOR** - CURSOR is used to iterate though many objects from the search results. An iteration
begins when the CURSOR is set to Zero or not included with the request, and completes when the
cursor returned by the server is Zero.

**NOFIELDS** - NOFIELDS tells the server that you do not want field values returned with the
search results.

**LIMIT** - LIMIT can be used to limit the number of objects returned for a single search request.

## Geofencing

A <a href="https://en.wikipedia.org/wiki/Geo-fence">geofence</a> is a virtual boundary that can
detect when an object enters or exits the area. This boundary can be a radius, bounding box, or
a polygon. The Bhojpur Space can turn any standard search into a geofence monitor by adding the
FENCE keyword to the search.

The *Bhojpur Space* also allows for [Webhooks](https://docs.bhojpur.net/commands/sethook) to be
assigned to Geofences.

A simple example:

```
> nearby fleet fence point 33.462 -112.268 6000
```

This command opens a geofence that monitors the 'fleet' collection. The server will respond with:

```
{"ok":true,"live":true}
```

And, the connection will be kept open. If any object enters or exits the 6 km radius around
`33.462,-112.268` the `Bhojpur Space` server will respond in realtime with a message such as:

```
{"command":"set","detect":"enter","id":"truck02","object":{"type":"Point","coordinates":[-112.2695,33.4626]}}
```

The `Bhojpur Space` server will notify the client, if the `command` is `del | set | drop`.

- `del` notifies the client that an object has been deleted from the collection that is being fenced.
- `drop` notifies the client that the entire collection is dropped.
- `set` notifies the client that an object has been added or updated, and when it's position is
detected by the fence.

The `detect` may be one of the following values.

- `inside` is when an object is inside the specified area.
- `outside` is when an object is outside the specified area.
- `enter` is when an object that **was not** previously in the fence has entered the area.
- `exit` is when an object that **was** previously in the fence has exited the area.
- `cross` is when an object that **was not** previously in the fence has entered **and** exited the area.

These can be used when establishing a geofence, to pre-filter responses. For instance, to limit
responses to `enter` and `exit` detections:

```
> nearby fleet fence detect enter,exit point 33.462 -112.268 6000
```

## Publish/Subscribe channels

The `Bhojpur Space` supports delivering geofence notications over pub/sub channels.

To create a static geofence that sends notifications, when a bus is within 200 meters of a point
and sends to the `busstop` channel:

```
> setchan busstop nearby buses fence point 33.5123 -112.2693 200
```

Subscribe on the `busstop` channel:

```
> subscribe busstop
```

## Object types

All `object types` except for `XYZ Tiles` and `QuadKeys` can be stored in a collection. The XYZ Tiles
and QuadKeys are reserved for the SEARCH keyword only.

### Lat/lon point

The most basic object type is a point that is composed of a `latitude` and a `longitude`. There is
an optional `z` member that may be used for auxiliary data such as elevation or a timestamp.

```
set fleet truck1 point 33.5123 -112.2693     # plain lat/lon
set fleet truck1 point 33.5123 -112.2693 225 # lat/lon with z member
```

### Bounding Box

A `bounding box` consists of two points. The first being the southwestern most point and the second
is the north-eastern most point.

```
set fleet truck1 bounds 30 -110 40 -100
```

### GeoHash

A [geohash](https://en.wikipedia.org/wiki/Geohash) is a string representation of a point. With
the length of the string indicating the precision of the point.

```
set fleet truck1 hash 9tbnthxzr # this would be equivalent to 'point 33.5123 -112.2693'
```

#### GeoJSON

A [GeoJSON](https://tools.ietf.org/html/rfc7946) is an industry standard format for representing
a variety of object types including a point, multi-point, linestring, multi-linestring, polygon,
multi-polygon, geometry-collection, feature, and feature-collection.

<i>* All ignored members will not persist.</i>

**Important to note that all coordinates are in Longitude, Latitude order.**

```
set city tempe object {"type":"Polygon","coordinates":[[[0,0],[10,10],[10,0],[0,0]]]}
```

### XYZ Tile

An `XYZ Tile` is rectangle bounding area on earth that is represented by an X, Y coordinate
and a Z (zoom) level. Check out [maptiler.org](http://www.maptiler.org/google-maps-coordinates-tile-bounds-projection/) for an interactive example.

### QuadKey

A `QuadKey` used the same coordinate system as an `XYZ Tile` except that the string
representation is a string characters composed of 0, 1, 2, or 3. For a detailed explanation
checkout [The Bing Maps Tile System](https://msdn.microsoft.com/en-us/library/bb259689.aspx).

## Network protocols

It's recommended to use a `client-side` library or the [Bhojpur Space CLI](#running), but there
are times when only HTTP is available or when you need to test from a remote terminal. In those
cases we provide an HTTP and telnet options.

### HTTP

One of the simplest ways to call a `Bhojpur Space` command is to use HTTP. From the command line
you can use [curl](https://curl.haxx.se/). For example:

```bash
# call with request in the body
$ curl --data "set fleet truck3 point 33.4762 -112.10923" localhost:9851

# call with request in the url path
$ curl localhost:9851/set+fleet+truck3+point+33.4762+-112.10923
```

### Web Sockets

Websockets can be used when you need to Geofence and keep the connection alive. It works just
like the HTTP example above, with the exception that the connection stays alive and the data
is sent from the server as text websocket messages.

### Telnet

There is the option to use a plain telnet connection. The default output through telnet is
[RESP](https://redis.io/topics/protocol).

```bash
$ telnet localhost 9851
set fleet truck3 point 33.4762 -112.10923
+OK
```

The `Bhopur Space` server will respond in [JSON](https://json.org) or [RESP](https://redis.io/topics/protocol)
depending on which protocol is used when initiating the first command.

- HTTP and WebSockets use JSON.
- Telnet and RESP clients use RESP.

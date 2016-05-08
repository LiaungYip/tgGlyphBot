package main

var databaseFilename = "glyphbot.boltdb"

// Used inside the cache DB.
// Incrementing this version number invalidates the cache, causing subsequent
// images to be created from scratch.
var programVersion = "0.0.3"

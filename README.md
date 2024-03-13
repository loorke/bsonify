# bsonify

This repository contains converters from Go collection types (i.e. maps and structs)
to BSON documents suitable for usage with MongoDB driver. The main usage scenario
is smart converting of an object for $set update without resetting existing fields 
with zero values. 

## E.g.
```Go
type b struct {
}

type c struct {
    Ca *int  `bson:"ca,omitempty"`
    Cb *uint `bson:"cb"`
}

type a struct {
    Aa string `bson:"aa,omitempty"`
    Ab bool   `bson:"custom_name,omitempty"`
    Ac string

    Ad b `bson:"ad,omitempty"`
    Ae c `bson:"ae"`
}

mgo.Database("test").Collection("test").InsertOne(ctx, a{Aa: "do_not_reset"})

mgo.Database("test").Collection("test").UpdateOne(ctx, bson.M{
    "_id": "test_doc_id",
}, bson.M{
    "$set": bsonify.SetUpdateD(a{Ab: true, Ae: c{Ca: &n}}),
})

mgo.Database("test").Collection("test").FindOne(ctx, bson.M{"_id": "test_doc_id"})

/*
    {
        "aa": "do_not_reset",
        "custom_name": true,
        "Ac": "",
        "ae": {
            "ca": 1,
            "cb": null
        }
    }
*/
```

## Converter functions

* `func SetUpdateM(v any) bson.M` 
Accepts a struct or map and returns a bson.M for mongo update $set operation. The struct or map can contain pointers, interfaces, maps, and structs. Map keys should be strings. BSON tags with omitempty are supported.
* `func SetUpdateD(v any) bson.D`
Accepts a struct or map and returns a bson.D for mongo update $set operation.
The struct or map can contain pointers, interfaces, maps, and structs.
Map keys should be strings. BSON tags with omitempty are supported.
* `func Dump(v any) bson.D`
Accepts a struct or map and returns a bson.D that resembles the original
object structure. Map keys should be strings. The struct or map can contain
pointers, interfaces, maps, and structs. BSON tags with omitempty are
supported.
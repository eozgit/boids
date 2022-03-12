package main

import "github.com/go-redis/redis/v8"

type Boid struct {
	id       int
	strId    string
	velocity *Vector
}

func (boid *Boid) getPosition() *Vector {
	pos := rdb.GeoPos(ctx, boid.strId, boid.strId).Val()[0]
	return &Vector{pos.Latitude * 10, pos.Longitude * 10}
}

func (boid *Boid) setPosition(position *Vector) {
	cmd := rdb.GeoAdd(ctx, boid.strId, &redis.GeoLocation{Name: boid.strId, Latitude: position.x / 10, Longitude: position.y / 10})
	er := cmd.Err()
	if er != nil {
		panic(er)
	}
}

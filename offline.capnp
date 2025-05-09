using Go = import "/go.capnp";
@0xda3a0d9284ca402f;
$Go.package("main");
$Go.import("offline");

struct Way {
  id @0 :Int64;
  name @1 :Text;
  ref @2 :Text;
  maxSpeed @3 :Float64;
  minLat @4 :Float64;
  minLon @5 :Float64;
  maxLat @6 :Float64;
  maxLon @7 :Float64;
  nodes @8 :List(Coordinates);
  lanes @9 :UInt8;
  advisorySpeed @10 :Float64;
  hazard @11 :Text;
  oneWay @12 :Bool;
  maxSpeedForward @13 :Float64;
  maxSpeedBackward @14 :Float64;
  maxSpeedPractical @15 :Float64;
  maxSpeedPracticalForward @16 :Float64;
  maxSpeedPracticalBackward @17 :Float64;
}

struct Coordinates {
  latitude @0 :Float64;
  longitude @1 :Float64;
}

struct Offline {
  minLat @0 :Float64;
  minLon @1 :Float64;
  maxLat @2 :Float64;
  maxLon @3 :Float64;
  ways @4 :List(Way);
  overlap @5 :Float64;
}

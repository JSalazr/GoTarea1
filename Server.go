package main

import (
	"net/http"
	"encoding/json"
	"log"
	"googlemaps.github.io/maps"
	"github.com/kr/pretty"
	"golang.org/x/net/context"
	"fmt"
	"bytes"
	"strconv"
	"regexp"
	//"encoding/binary"
)

type Address struct {
	Origen string
	Destino string
}

type NearbyRestaurants struct{
	Location string
}

type Dir struct{
	Lat float64 "json: \"lat\""
	Lon float64 "json: \"lon\""
}


func ejercicio1(w http.ResponseWriter, r *http.Request) {
	var address Address
  if r.Body == nil {
      http.Error(w, "Please send a request body", 400)
      return
  }
  er := json.NewDecoder(r.Body).Decode(&address)
  if er != nil {
      http.Error(w, er.Error(), 400)
      return
  }
	fmt.Println(address)
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyAzzrnc71pLvEvOdY322DQwwbUsFQZT7Vg"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	pat := regexp.MustCompile("[a-zA-Z\\d ]*, [a-zA-Z\\d ]*")
	sap := pat.FindString(address.Origen)
	sap2 := pat.FindString(address.Destino)
	fmt.Println(sap);
	a := &maps.DirectionsRequest{
		Origin:      sap,
		Destination: sap2,
	}
	resp, _, b := c.Directions(context.Background(), a)
	if b != nil {
		log.Fatalf("fatal error: %s", b)
	}
	pretty.Println(resp)
	buffer := new(bytes.Buffer)
	buffer.WriteString("{\"ruta\":[")
	json.NewDecoder(r.Body).Decode(&resp)
	//ret := make(map[string]Dir, len(resp[0].Legs[0].Steps)+1)
	//pretty.Println(resp[0].Legs[0].Steps[0].StartLocation)
	for x :=0; x< len(resp[0].Legs[0].Steps); x++{
		buffer.WriteString("{\"lat\":")
		buffer.WriteString(strconv.FormatFloat(resp[0].Legs[0].Steps[x].StartLocation.Lat,'f',5, 64))
		buffer.WriteString(", ")
		buffer.WriteString("\"lon\":")
		buffer.WriteString(strconv.FormatFloat(resp[0].Legs[0].Steps[x].StartLocation.Lng,'f',5, 64))
		buffer.WriteString("}, ")
		if x == (len(resp[0].Legs[0].Steps) - 1){
			buffer.WriteString("{\"lat\":")
			buffer.WriteString(strconv.FormatFloat(resp[0].Legs[0].Steps[x].EndLocation.Lat,'f',5, 64))
			buffer.WriteString(", ")
			buffer.WriteString("\"lon\":")
			buffer.WriteString(strconv.FormatFloat(resp[0].Legs[0].Steps[x].EndLocation.Lng,'f',5, 64))
			buffer.WriteString("} ")
		}
	}

	buffer.WriteString("]}")
	pretty.Println(buffer.String());
	fmt.Fprintf(w, buffer.String());
		/*data := map[string]interface{}{"ruta": ret }
		j, y := json.Marshal(data)
		if y != nil {
			log.Fatalf("fatal error: %s", y)
		}
		pretty.Println(ret)
		pretty.Println(string(j))*/
}

func ejercicio2(w http.ResponseWriter, r *http.Request) {
	var nearby NearbyRestaurants
  if r.Body == nil {
      http.Error(w, "Please send a request body", 400)
      return
  }
	json.NewDecoder(r.Body).Decode(&nearby)
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyAzzrnc71pLvEvOdY322DQwwbUsFQZT7Vg"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	a := &maps.DirectionsRequest{
		Origin:      nearby.Location,
		Destination: nearby.Location,
	}
	resp, _, b := c.Directions(context.Background(), a)
	if b != nil {
		log.Fatalf("fatal error: %s", b)
	}
	json.NewDecoder(r.Body).Decode(&resp)
	pretty.Println(resp)

	c, _ = maps.NewClient(maps.WithAPIKey("AIzaSyDYY9Dny4Iu9eFN0wU6EuVeAq2KsPQXKcY"))
	t := &maps.NearbySearchRequest{
		Location:   &maps.LatLng{resp[0].Legs[0].Steps[0].StartLocation.Lat, resp[0].Legs[0].Steps[0].StartLocation.Lng},
		Radius: 1000,
		Type: "restaurant",
	}

	resp1, _ := c.NearbySearch(context.Background(), t)
	json.NewDecoder(r.Body).Decode(&resp1)

	buffer := new(bytes.Buffer)
	buffer.WriteString("{\"restaurantes\":[")
	for x :=0; x< len(resp1.Results); x++{
		buffer.WriteString("{\"nombre\":\"")
		buffer.WriteString(resp1.Results[x].Name)
		buffer.WriteString("\", ")
		buffer.WriteString("\"lat\":")
		buffer.WriteString(strconv.FormatFloat(resp1.Results[x].Geometry.Location.Lat,'f',5, 64))
		buffer.WriteString(", ")
		buffer.WriteString("\"lon\":")
		buffer.WriteString(strconv.FormatFloat(resp1.Results[x].Geometry.Location.Lng,'f',5, 64))
		if x == (len(resp1.Results) - 1){
			buffer.WriteString("}")
		}else{
			buffer.WriteString("}, ")
		}
	}
	buffer.WriteString("]}")
	fmt.Fprintf(w, buffer.String())
}

func main() {
	http.HandleFunc("/ejercicio1", ejercicio1);
	http.HandleFunc("/ejercicio2", ejercicio2);
	http.ListenAndServe(":8080", nil)
}

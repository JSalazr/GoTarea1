package main

import (
	"net/http"
	"encoding/json"
	"googlemaps.github.io/maps"
	"github.com/kr/pretty"
	"golang.org/x/net/context"
	"fmt"
	"bytes"
	"strconv"
	"regexp"
	"encoding/base64"
	"image"
	"image/color"
	"strings"
	"golang.org/x/image/bmp"
	//"io"
	//"encoding/binary"
)

type Address struct {
	Origen string
	Destino string
}

type NearbyRestaurants struct{
	Origen string
}

type ImgStr struct{
	Nombre string
	Data string
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
  er1 := json.NewDecoder(r.Body).Decode(&address)
	if er1 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen\"}", 400);
		return
	}
	fmt.Println(address)
	c, er1 := maps.NewClient(maps.WithAPIKey("AIzaSyAzzrnc71pLvEvOdY322DQwwbUsFQZT7Vg"))
	if er1 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen\"}", 500);
		return
	}

	pat := regexp.MustCompile("a-zA-Z\\d\\, ]*")
	sap := pat.FindString(address.Origen)
	sap2 := pat.FindString(address.Destino)
	fmt.Println(sap);
	a := &maps.DirectionsRequest{
		Origin:      sap,
		Destination: sap2,
	}
	resp, _, er1 := c.Directions(context.Background(), a)
	if er1 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen\"}", 500);
		return
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
	er2 :=json.NewDecoder(r.Body).Decode(&nearby)
	if er2 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen\"}", 400);
		return
	}
	c, er2 := maps.NewClient(maps.WithAPIKey("AIzaSyAzzrnc71pLvEvOdY322DQwwbUsFQZT7Vg"))
	if er2 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen\"}", 500);
		return
	}

	a := &maps.DirectionsRequest{
		Origin:      nearby.Origen,
		Destination: nearby.Origen,
	}
	resp, _, er2 := c.Directions(context.Background(), a)
	if er2 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen\"}", 500);
		return
	}
	json.NewDecoder(r.Body).Decode(&resp)
	pretty.Println(resp)

	c, er2 = maps.NewClient(maps.WithAPIKey("AIzaSyDYY9Dny4Iu9eFN0wU6EuVeAq2KsPQXKcY"))
	if er2 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen\"}", 500);
		return
	}
	t := &maps.NearbySearchRequest{
		Location:   &maps.LatLng{resp[0].Legs[0].Steps[0].StartLocation.Lat, resp[0].Legs[0].Steps[0].StartLocation.Lng},
		Radius: 1000,
		Type: "restaurant",
	}

	resp1, er2 := c.NearbySearch(context.Background(), t)
	if er2 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen\"}", 500);
		return
	}
	er2 = json.NewDecoder(r.Body).Decode(&resp1)
	if er2 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen\"}", 500);
		return
	}

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

func ejercicio3(w http.ResponseWriter, req *http.Request) {
	var imgInc ImgStr
	if req.Body == nil {
      http.Error(w, "Please send a request body1", 400)
      return
  }
	er3 := json.NewDecoder(req.Body).Decode(&imgInc)
	if er3 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen2\"}", 400);
		return
	}

	u, er3 := base64.StdEncoding.DecodeString(imgInc.Data)
	if er3 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen3\"}", 500);
		return
	}
	r := bytes.NewReader(u)
	im,er3:= bmp.Decode(r)
	if er3 !=nil {
		http.Error(w, "{\"error\":\"No se especifico origen4\"}", 500);
		return
	}

	pretty.Println(im.Bounds().Min.Y)

	bounds := im.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	imgbw := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			old := im.At(x, y)
			r, g, b, _ := old.RGBA()

			avg := (r + g + b) / 3
			new := color.Gray{uint8(avg / 256)}

			imgbw.Set(x, y, new)
		}
	}

	//json.NewEncoder(w).Encode(imgSet)

	imgBuf := new(bytes.Buffer)
	bmp.Encode(imgBuf, imgbw)
	bwImg := base64.StdEncoding.EncodeToString(imgBuf.Bytes())
	fmt.Println(bwImg)
	sep := strings.Split(imgInc.Nombre, ".");
	buffer := new(bytes.Buffer)
	buffer.WriteString("{\"nombre\": \"")
	buffer.WriteString(sep[0])
	buffer.WriteString("(blanco y negro).")
	buffer.WriteString(sep[1])
	buffer.WriteString("\", \"data\": \"")
	buffer.WriteString(bwImg)
	buffer.WriteString("\"}")
	pretty.Println(buffer.String())
	fmt.Fprintf(w, buffer.String())
	//pretty.Println(width, height)
}

func main() {
	http.HandleFunc("/ejercicio1", ejercicio1);
	http.HandleFunc("/ejercicio2", ejercicio2);
	http.HandleFunc("/ejercicio3", ejercicio3);
	http.ListenAndServe(":8080", nil)
}

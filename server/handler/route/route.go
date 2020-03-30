package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

type Route struct {
	CoordinatesAccuracy string `json:"coordinates_accuracy"`
	Altitude string `json:"altitude"`
	AltitudeAccuracy string `json:"altitude_accuracy"`
	Bearing string `json:"bearing"`
	BearingAccuracy string `json:"bearing_accuracy"`
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
	AMapCoordinates string `json:"amap_coordinates"`
	Speed string `json:"speed"`
	SpeedAccuracy string `json:"speed_accuracy"`
	TimeSeconds string `json:"time_seconds"`
}

func Record(w http.ResponseWriter, r *http.Request) {
	uid := r.FormValue("uid")
	key := r.FormValue("key")
	if key != os.Args[6] {
		return
	}
	Msg := Route{
		CoordinatesAccuracy: filter(r.FormValue("coordinates_accuracy")),
		Altitude:            filter(r.FormValue("altitude")),
		AltitudeAccuracy:    filter(r.FormValue("altitude_accuracy")),
		Bearing:             filter(r.FormValue("bearing")),
		BearingAccuracy:     filter(r.FormValue("bearing_accuracy")),
		Latitude:            filter(r.FormValue("latitude")),
		Longitude:           filter(r.FormValue("longitude")),
		Speed:               filter(r.FormValue("speed")),
		SpeedAccuracy:       filter(r.FormValue("speed_accuracy")),
		TimeSeconds:         filter(r.FormValue("time_seconds")),
	}
	url := fmt.Sprintf("https://restapi.amap.com/v3/assistant/coordinate/convert?key=%s&locations=%s,%s&coordsys=gps", os.Args[7], Msg.Longitude, Msg.Latitude)
	//fmt.Println("url1", url)
	resp, err := http.Get(url)
	type aMapS struct {
		Status string `json:"status"`
		Info string `json:"info"`
		Locations string `json:"locations"`
	}

	aMap := &aMapS{}

	if err != nil {
		Msg.AMapCoordinates = ""
	}else {
		s, err := ioutil.ReadAll(resp.Body)
		_=resp.Body.Close()
		if err == nil {
			_=json.Unmarshal(s, aMap)
			Msg.AMapCoordinates = aMap.Locations
		}
	}

	url = fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?key=%s&location=%s&poitype=&radius=0&extensions=all&batch=false&roadlevel=0", os.Args[7], Msg.AMapCoordinates)
	//fmt.Println("url2", url)
	resp, err = http.Get(url)
	aMapLoc := ""
	if err == nil {
		s, err := ioutil.ReadAll(resp.Body)
		_=resp.Body.Close()
		if err == nil {
			//fmt.Printf(string(s))
			aMapLoc = string(s)
		}
	}

	data, err := json.Marshal(Msg)
	if err != nil {
		data = make([]byte,0)
	}
	appendToFile(filename(uid), now(), data, aMapLoc)

}

func appendToFile(file string, time string, str []byte, loc string) {

	pth := path.Dir(file)
	if !isExist(pth) {
		err := os.MkdirAll(pth, os.ModePerm)
		if err != nil{
			fmt.Println(err)
			return
		}
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		fmt.Printf("Cannot open file %s!\n", file)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = f.WriteString(time) // time
	if err != nil {
		fmt.Println(err)
	}
	_, err = f.Write(append(str, ' '))  // time data
	if err != nil {
		fmt.Println(err)
	}
	_, err = f.WriteString(loc) // time data loc
	if err != nil {
		fmt.Println(err)
	}
	_, err = f.Write([]byte{'\n'})
	if err != nil {
		fmt.Println(err)
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil{
		if os.IsExist(err){
			return true
		}
		if os.IsNotExist(err){
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

func filename(uid string) string {
	return fmt.Sprintf("./loc/%s/%s", uid,time.Now().Format("2006-01-02"))
}

func now() string {
	return time.Now().Format("2006-01-02 15:04:05 ")
}

func filter(str string) string {
	switch str {
	case "%%gl_coordinates_accuracy":
	case "%gl_altitude":
	case "%gl_altitude_accuracy":
	case "%gl_bearing":
	case "%gl_bearing_accuracy":
	case "%gl_latitude":
	case "%gl_longitude":
	case "%gl_speed":
	case "%gl_speed_accuracy":
	case "%gl_time_seconds":
		return ""
	default:
		return str
	}
	return ""
}

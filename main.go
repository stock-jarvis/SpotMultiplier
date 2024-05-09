package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	input_path1 := "/home/ajeeb/sandbox/src/SpotMultiplier/SpotData22"
	input_path2 := "/home/ajeeb/sandbox/src/SpotMultiplier/spot multiplier - Full n final.csv"
	paths := []string{}
	pc := 0
	filepath.Walk(input_path1, func(path string, info os.FileInfo, err error) error {
		if err == nil && strings.Contains(info.Name(), ".csv") {
			paths = append(paths, path)
			pc++
		}
		return nil
	})
	log.Printf("\nPaths:= %d", pc)
	file, _ := os.Open(input_path2)
	csvReader := csv.NewReader(file)
	contents, _ := csvReader.ReadAll()
	// inst := 0
	// from := []string{}
	// to := []string{}
	// mul := []string{}
	for _, filepath := range paths {
		path_name := strings.Split(filepath, "/")
		name := strings.Trim(path_name[len(path_name)-1], "_1m.csv")
		for i := 1; i <= len(contents); i++ {
			if name == contents[i][0] {
				log.Printf("\nPathname matched: %v %v", name, contents[i][0])
				// from[inst] = contents[i][1]
				// to[inst] = contents[i][2]
				// mul[inst] = contents[i][3]
				// inst++
				SpotMultiplier(filepath, contents[i][1], contents[i][2], contents[i][3])
				break
			}
			//SpotMultiplier(filepath, inst, from, to, mul)
		}
	}
}

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func SpotMultiplier(inPath1 string, fromDate string, toDate string, multiplier string) {
	log.Println("Entering Spot Multiplier")
	file1, err := os.OpenFile(inPath1, os.O_RDWR, 2) //Original spot data (read/overwrite with multiplued data)
	if err != nil {
		log.Println("Error in opening file1")
		return
	}
	csvWriter1 := csv.NewWriter(file1)
	csvReader1 := csv.NewReader(file1)
	name := strings.Split(inPath1, "/")
	outpath1 := fmt.Sprintf("/home/ajeeb/sandbox/src/SpotMultiplier/Multiplied files/%v", name[len(name)-1])
	outpath2 := fmt.Sprintf("/home/ajeeb/sandbox/src/SpotMultiplier/Original Files/%v", name[len(name)-1])
	file2, err := os.OpenFile(outpath1, os.O_WRONLY, 1) //Stores a copy of multiplied files
	if err != nil {
		file2, err = os.Create(outpath1)
		if err != nil {
			log.Println("Error in opening file2")
			return
		} else {
			log.Println("File2 created")
		}
	} else {
		log.Println("File opened")
	}
	file3, err := os.OpenFile(outpath2, os.O_WRONLY, 1) //stores a copy of the original file
	if err != nil {
		file3, err = os.Create(outpath2)
		if err != nil {
			log.Println("Error in opening file")
			return
		} else {
			log.Println("File3 created")
		}
	} else {
		log.Println("File3 opened")
	}
	csvWriter2 := csv.NewWriter(file2)
	csvWriter3 := csv.NewWriter(file3)
	contents, err := csvReader1.ReadAll()
	if err != nil {
		log.Println("Error in reading contents")
		return
	}

	//for x := 0; x < inst; x++ {
	mul, err := strconv.ParseFloat(multiplier, 64)
	if err != nil {
		log.Println("Error in parsing multiplier")
		return
	}
	from := strings.Split(fromDate, "-") //splits from date dd-mm-yy
	to := strings.Split(toDate, "-")     //splits to date dd-mm-yy
	from1, err := strconv.ParseInt(from[0], 10, 64)
	if err != nil {
		log.Println("Error in parsing from1")
		return
	}
	from2, err := strconv.ParseInt(from[1], 10, 64)
	if err != nil {
		log.Println("Error in parsing from2")
		return
	}
	from3, err := strconv.ParseInt(from[2], 10, 64)
	if err != nil {
		log.Println("Error in parsing from3")
		return
	}
	to1, err := strconv.ParseInt(to[0], 10, 64)
	if err != nil {
		log.Println("Error in parsing to1")
		return
	}
	to2, err := strconv.ParseInt(to[1], 10, 64)
	if err != nil {
		log.Println("Error in parsing to2")
		return
	}
	to3, err := strconv.ParseInt(to[2], 10, 64)
	if err != nil {
		log.Println("Error in parsing to3")
		return
	}

	for i := 1; i < len(contents)-1; i++ {
		//log.Println("Entering for loop")
		log.Printf(" ")
		ts, _ := strconv.ParseInt(contents[i][0], 10, 64)
		ts = ts + 19800 //IST to GMT
		x := time.Unix(ts, 0)
		y, m, d := x.Date() //extracts date from unix ts in yy-mm-dd
		o, _ := strconv.ParseFloat(contents[i][1], 64)
		h, _ := strconv.ParseFloat(contents[i][2], 64)
		l, _ := strconv.ParseFloat(contents[i][3], 64)
		c, _ := strconv.ParseFloat(contents[i][4], 64)
		v, _ := strconv.ParseFloat(contents[i][5], 64)
		ogdata := []string{
			fmt.Sprintf("%d", ts),
			fmt.Sprintf("%0.2f", o),
			fmt.Sprintf("%0.2f", h),
			fmt.Sprintf("%0.2f", l),
			fmt.Sprintf("%0.2f", c),
			fmt.Sprintf("%0.2f", v),
		}
		//log.Println("Printing og file")
		csvWriter3.Write(ogdata) //making copy of og file
		csvWriter3.Flush()
		d1 := int64(d)
		m1 := int64(m)
		y1 := int64(y)
		//log.Printf("%d %d %d", d1, m1, y1)
		// log.Printf("%d %d %d", from1, from2, from3)
		// log.Printf("%d %d %d", to1, to2, to3)
		if d1 >= from1 && m1 >= from2 && y1 >= from3 {
			if to1 >= d1 && to2 >= m1 && to3 >= y1 {
				o = o * mul
				h = h * mul
				l = l * mul
				c = c * mul
			}
			rep_data := []string{
				fmt.Sprintf("%d", ts),
				fmt.Sprintf("%0.2f", Round(o, 0.05)),
				fmt.Sprintf("%0.2f", Round(h, 0.05)),
				fmt.Sprintf("%0.2f", Round(l, 0.05)),
				fmt.Sprintf("%0.2f", Round(c, 0.05)),
				fmt.Sprintf("%0.2f", v),
			}
			log.Println("Printing replaced file")
			csvWriter1.Write(rep_data) //overwrites og file with new data
			csvWriter1.Flush()
			csvWriter2.Write(rep_data) //makes a copy of file with replaced data
			csvWriter2.Flush()
		}
	}
}

//}
